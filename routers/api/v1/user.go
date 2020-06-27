package v1

import (
	"context"
	"net/http"
	"set-flags/models"
	"set-flags/pkg/e"
	"set-flags/pkg/logging"
	"set-flags/pkg/setting"
	"set-flags/pkg/utils"
	"set-flags/schemas"

	"github.com/fox-one/mixin-sdk"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
)

// CheckRewards check the total rewards received by the user for the flag
func CheckRewards(c *gin.Context) {
	code := e.INVALID_PARAMS

	var pagination schemas.Pagination

	c.ShouldBindQuery(&pagination)

	if pagination.CurrentPage == 0 {
		pagination.CurrentPage = 1
	}

	if pagination.PageSize == 0 {
		pagination.PageSize = setting.GetConfig().App.PageSize
	}

	var checkReward schemas.CheckReward

	if err := c.BindUri(&checkReward); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  err.Error(),
			"data": make(map[string]interface{}),
		})
		return
	}

	userID, _ := uuid.FromString(checkReward.UserID)

	flagID, _ := uuid.FromString(checkReward.FlagID)

	data, total := models.FindEvidenceByFlagIDAndAttachmentID(flagID, userID, pagination.CurrentPage, pagination.PageSize)

	code = e.SUCCESS
	c.JSON(http.StatusOK, gin.H{
		"code":  code,
		"msg":   e.GetMsg(code),
		"data":  data,
		"total": total,
	})
}

// Me current user profile
func Me(c *gin.Context) {
	code := e.INVALID_PARAMS

	var header schemas.Header

	if err := c.BindHeader(&header); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  err.Error(),
			"data": make(map[string]interface{}),
		})
		return
	}
	userID, _ := uuid.FromString(header.XUSERID)

	if !models.UserExist(userID) {
		code = http.StatusForbidden
		c.JSON(code, gin.H{
			"code": code,
			"msg":  "current user not exist.",
			"data": make(map[string]interface{}),
		})
		return
	}

	accessToken := models.FindUserAccessToken(userID)

	ctx := context.Background()
	profile, err := mixin.FetchProfile(ctx, accessToken)

	if err != nil {

		code = e.ERROR_AUTH_TOKEN

		logging.Info("fetch user profile failed", err.Error())

		c.JSON(http.StatusForbidden, gin.H{
			"code": http.StatusForbidden,
			"msg":  err.Error(),
			"data": make(map[string]interface{}),
		})
		return
	}

	userSchema := models.UserSchema{
		UserID:         profile.UserID,
		IdentityNumber: profile.IdentityNumber,
		FullName:       profile.FullName,
		AvatarURL:      profile.AvatarURL,
	}

	code = e.SUCCESS
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": userSchema,
	})
}

// Auth auth
func Auth(c *gin.Context) {
	code := e.INVALID_PARAMS

	authorizationCode := c.Query("code")
	logging.Info("authorizationCode", authorizationCode)

	ctx := context.Background()

	accessToken, _, err := mixin.AuthorizeToken(ctx,
		setting.GetConfig().Bot.ClientID.String(),
		setting.GetConfig().Bot.ClientSecret,
		authorizationCode,
		setting.GetConfig().Bot.CodeVerifier)

	if err != nil {
		code = e.ERROR_AUTH_TOKEN

		logging.Info("fetch access token failed", err.Error())

		c.JSON(http.StatusInternalServerError, gin.H{
			"code": code,
			"msg":  err.Error(),
			"data": make(map[string]interface{}),
		})
		return
	}

	profile, err := mixin.FetchProfile(ctx, accessToken)

	if err != nil {

		code = e.ERROR_AUTH_TOKEN

		logging.Info("fetch user profile failed", err.Error())

		c.JSON(http.StatusInternalServerError, gin.H{
			"code": code,
			"msg":  err.Error(),
			"data": make(map[string]interface{}),
		})
		return
	}

	userID, _ := uuid.FromString(profile.UserID)

	// update user info and access token
	if models.UserExist(userID) {
		logging.Info("update user")
		models.UpdateUser(profile, accessToken)
		models.UpdateFlagUserInfo(profile)
	} else {
		// create user
		logging.Info("create user")
		models.CreateUser(profile, accessToken)
	}
	// mixin auth success
	// generate app Bearer token
	token, err := utils.GenerateToken(profile.UserID)
	if err != nil {
		code = e.ERROR_AUTH_TOKEN
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"msg":  e.GetMsg(code),
			"data": make(map[string]interface{}),
		})
		return
	}

	authToken := schemas.AuthToken{
		Token:  token,
		UserID: profile.UserID,
	}

	code = e.SUCCESS
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": authToken,
	})
}
