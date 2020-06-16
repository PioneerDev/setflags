package v1

import (
	"context"
	"fmt"
	"net/http"
	"set-flags/models"
	"set-flags/pkg/e"
	"set-flags/pkg/setting"
	"strconv"

	"github.com/fox-one/mixin-sdk"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
)

// CheckRewards check the total rewards received by the user for the flag
func CheckRewards(c *gin.Context) {
	code := e.INVALID_PARAMS

	userID, err := uuid.FromString(c.Param("user_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": code,
			"msg":  err.Error(),
			"data": make(map[string]interface{}),
		})
		return
	}
	flagID, err := uuid.FromString(c.Param("flag_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": code,
			"msg":  err.Error(),
			"data": make(map[string]interface{}),
		})
		return
	}

	currentPage, err := strconv.Atoi(c.DefaultQuery("current_page", "1"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": code,
			"msg":  e.GetMsg(code),
			"data": make(map[string]interface{}),
		})
		return
	}

	pageSize, err := strconv.Atoi(c.DefaultQuery("page_size", setting.PageSize))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": code,
			"msg":  e.GetMsg(code),
			"data": make(map[string]interface{}),
		})
		return
	}

	data := models.FindEvidenceByFlagIDAndAttachmentID(flagID, userID, currentPage, pageSize)

	code = e.SUCCESS
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": data,
	})
}

// Me current user profile
func Me(c *gin.Context) {
	code := e.INVALID_PARAMS

	userID, err := uuid.FromString(c.GetHeader("x-user-id"))

	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"code": code,
			"msg":  err.Error(),
			"data": make(map[string]interface{}),
		})
		return
	}

	user := models.FindUserByID(userID)

	data := map[string]string{
		"id":         userID.String(),
		"full_name":  user.FullName,
		"avatar_url": user.AvatarURL,
	}

	code = e.SUCCESS
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": data,
	})
}

// Auth auth
func Auth(c *gin.Context) {

	authorizationCode := c.Query("token")

	ctx := context.Background()

	accessToken, _, err := mixin.AuthorizeToken(ctx, setting.ClientID.String(), setting.ClientSecret, authorizationCode, setting.CodeVerifier)

	profile, err := mixin.FetchProfile(ctx, accessToken)

	code := e.INVALID_PARAMS

	if err != nil {
		code = e.ERROR_AUTH_TOKEN
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": code,
			"msg":  e.GetMsg(code),
			"data": make(map[string]interface{}),
		})
		return
	}

	// update user info and access token
	if models.UserExist(profile.UserID) {
		models.UpdateUser(profile, accessToken)
		models.UpdateFlagUserInfo(profile)
	} else {
		// create user
		models.CreateUser(profile, accessToken)
	}

	code = e.SUCCESS
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": make(map[string]interface{}),
	})
}
