package v1

import (
	"context"
	"fmt"
	"github.com/fox-one/mixin-sdk"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"net/http"
	"set-flags/models"
	"set-flags/pkg/e"
	"set-flags/pkg/setting"
	"strconv"
)

// check the total rewards received by the user for the flag
func CheckRewards(c *gin.Context) {
	code := e.INVALID_PARAMS

	userId := c.Param("user_id")
	fmt.Printf("userId: %s\n", userId)
	flagId := c.Param("flag_id")
	fmt.Printf("flagId: %s\n", flagId)
	userID, err := uuid.FromString(userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": code,
			"msg":  err.Error(),
			"data": make(map[string]interface{}),
		})
		return
	}
	flagID, err := uuid.FromString(flagId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": code,
			"msg":  err.Error(),
			"data": make(map[string]interface{}),
		})
		return
	}

	currentPage_ := c.DefaultQuery("current_page", "1")
	pageSize_ := c.DefaultQuery("page_size", setting.PageSize)

	currentPage, err := strconv.Atoi(currentPage_)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": code,
			"msg": e.GetMsg(code),
			"data": make(map[string]interface{}),
		})
		return
	}

	pageSize, err := strconv.Atoi(pageSize_)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": code,
			"msg": e.GetMsg(code),
			"data": make(map[string]interface{}),
		})
		return
	}

	data := models.FindEvidenceByFlagIdAndAttachmentId(flagID, userID, currentPage, pageSize)

	code = e.SUCCESS
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": data,
	})
}

func Me(c *gin.Context) {
	code := e.INVALID_PARAMS
	userId := c.GetHeader("x-user-id")

	_, err := uuid.FromString(userId)

	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"code": code,
			"msg":  err.Error(),
			"data": make(map[string]interface{}),
		})
		return
	}

	user := models.FindUserById(userId)

	data := map[string]string{
		"id":         userId,
		"full_name":  user.FullName,
		"avatar_url": user.AvatarUrl,
	}

	code = e.SUCCESS
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": data,
	})
}

func Auth(c *gin.Context) {

	authorizationCode := c.Query("token")

	ctx := context.Background()

	accessToken, _, err := mixin.AuthorizeToken(ctx, setting.ClientId, setting.ClientSecret, authorizationCode, setting.CodeVerifier)

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
