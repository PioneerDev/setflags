package v1

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"io/ioutil"
	"net/http"
	"set-flags/models"
	"set-flags/pkg/e"
	"set-flags/pkg/setting"
	"set-flags/pkg/utils"
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

	authorizationCode := ""

	client := &http.Client{}

	code := e.INVALID_PARAMS

	accessToken, err := FetchAccessToken(client, authorizationCode)

	if err != nil {
		code = e.ERROR_AUTH_TOKEN
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": code,
			"msg":  e.GetMsg(code),
			"data": make(map[string]interface{}),
		})
		return
	}

	userInfo, err := FetchUserInfo(client, accessToken)

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
	if models.UserExist(userInfo.UserId) {
		models.UpdateUser(&userInfo, accessToken)
	} else {
		// create user
		models.CreateUser(&userInfo, accessToken)
	}

	code = e.SUCCESS
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": make(map[string]interface{}),
	})
}

// Fetch user access_token from Mixin
func FetchAccessToken(client *http.Client, code string) (string, error) {
	body := map[string]interface{}{}

	body["client_id"] = setting.ClientId.String()
	body["code"] = code

	bt, _ := json.Marshal(body)

	url := fmt.Sprintf("%s/oauth/token", setting.MixinAPIDomain)

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(bt))

	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)

	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var authResp map[string]interface{}

	data, _ := ioutil.ReadAll(resp.Body)

	_ = json.Unmarshal(data, &authResp)

	token, _ := authResp["access_token"].(string)

	return token, nil
}

// Fetch user info from Mixin
func FetchUserInfo(client *http.Client, accessToken string) (utils.UserInfo, error) {
	url := fmt.Sprintf("%s/me", setting.MixinAPIDomain)

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	resp, err := client.Do(req)

	if err != nil {
		return utils.UserInfo{}, err
	}
	defer resp.Body.Close()

	var authResp map[string]utils.UserInfo

	data, _ := ioutil.ReadAll(resp.Body)

	_ = json.Unmarshal(data, &authResp)

	user, _ := authResp["data"]

	return user, nil
}
