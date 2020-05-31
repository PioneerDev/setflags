package v1

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"set-flags/models"
	"set-flags/pkg/e"
	"set-flags/pkg/setting"
	"strconv"
)

// list all the flags
func ListFlags(c *gin.Context) {
	code := e.INVALID_PARAMS
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

	data := models.GetAllFlags(pageSize, currentPage)

	code = e.SUCCESS
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg": e.GetMsg(code),
		"data": data,
	})
}

// create a flag
func CreateFlag(c *gin.Context) {
	code := e.INVALID_PARAMS
	var flag map[string]interface{}

	if c.ShouldBind(&flag) == nil {
		fmt.Println(flag)
		payerId, err := uuid.FromString(flag["payer_id"].(string))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": code,
				"msg":  err.Error(),
				"data": make(map[string]interface{}),
			})
			return
		}

		assetId, err := uuid.FromString(flag["asset_id"].(string))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": code,
				"msg":  err.Error(),
				"data": make(map[string]interface{}),
			})
			return
		}

		flag["payer_id"] = payerId
		flag["asset_id"] = assetId
		models.CreateFlag(flag)
	}
	code = e.SUCCESS
	c.JSON(http.StatusCreated, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": make(map[string]interface{}),
	})
}

// Update an existing flag
func UpdateFlag(c *gin.Context) {
	code := e.INVALID_PARAMS
	flagId := c.Param("id")

	_, err := uuid.FromString(flagId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": code,
			"msg":  err.Error(),
			"data": make(map[string]interface{}),
		})
		return
	}

	op := c.Param("op")

	if op != "yes" && op != "no" && op != "done" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": code,
			"msg":  fmt.Sprintf("op: %s is invalid.", op),
			"data": make(map[string]interface{}),
		})
		return
	}

	if !models.FLagExists(flagId) {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "Flag not found.",
			"data": make(map[string]interface{}),
		})
		return
	}

	flag := models.FindFlagByID(flagId)

	if flag.Status != "done" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": code,
			"msg":  "not yet upload the evidence.",
			"data": make(map[string]interface{}),
		})
		return
	}

	code = e.SUCCESS
	c.JSON(http.StatusBadRequest, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": make(map[string]interface{}),
	})
}

// list all flags of the user
func FindFlagsByUserID(c *gin.Context) {
	code := e.INVALID_PARAMS
	userId := c.GetHeader("x-user-id")

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

	_, err = uuid.FromString(userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": code,
			"msg":  err.Error(),
			"data": make(map[string]interface{}),
		})
		return
	}

	flags := models.FindFlagsByUserID(userId, currentPage, pageSize)

	code = e.SUCCESS
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": flags,
	})
}

// upload evidence
func UploadEvidence(c *gin.Context) {
	code := e.INVALID_PARAMS
	userId := c.GetHeader("user_id")
	flagId := c.Query("flag_id")
	attachmentId := c.Param("attachment_id")

	fmt.Sprintf("attachmentId: %s, flagId: %s", attachmentId, flagId)

	type_ := c.Query("type")

	if type_ != "image" && type_ != "audio" && type_ != "video" && type_ != "document" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": code,
			"msg":  fmt.Sprintf("type: %s is invalid.", type_),
			"data": make(map[string]interface{}),
		})
		return
	}

	if !models.FLagExists(flagId) {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "not found specific flag.",
			"data": make(map[string]interface{}),
		})
		return
	}

	fileHeader, err := c.FormFile("file")
	if err != nil {
		code = e.ERROR
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": code,
			"msg":  err.Error(),
			"data": make(map[string]interface{}),
		})
		return
	}

	log.Println(fileHeader.Filename)

	client := &http.Client{}
	// read access token from db
	accessToken, _ := models.FindUserToken(userId)

	viewUrl, err := UploadAttachment(client, fileHeader, accessToken)
	if err != nil {
		code = e.ERROR_UPLOAD_ATTACHMENT
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"msg": e.GetMsg(code),
			"data": make(map[string]interface{}),
		})
		return
	}

	attachmentID, _ := uuid.FromString(attachmentId)
	flagID, _ := uuid.FromString(flagId)
	models.CreateEvidence(attachmentID, flagID, type_, viewUrl)

	// update flag status to `done`
	models.UpdateFlagStatus(flagId, "done")

	code = e.SUCCESS
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg": fmt.Sprintf("'%s' uploaded!", fileHeader.Filename),
		"data": make(map[string]interface{}),
	})
}

// list all the evidences since yesterday
func ListEvidences(c *gin.Context) {

	code := e.INVALID_PARAMS

	flagId := c.Param("flag_id")
	flagID, _ := uuid.FromString(flagId)

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

	data := models.FindEvidencesByFlag(flagID, currentPage, pageSize)

	code = e.SUCCESS
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg": e.GetMsg(code),
		"data": data,
	})
}

// https://developers.mixin.one/api/l-messages/create-attachment/
func UploadAttachment(client *http.Client, fileHeader *multipart.FileHeader, accessToken string) (string, error) {

	f, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	size := fileHeader.Size

	buffer := make([]byte, size)
	_, err = f.Read(buffer)

	if err != nil {
		return "", err
	}


	url := fmt.Sprintf("%s/attachments", setting.MixinAPIDomain)

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(buffer))

	req.Header.Add("Content-Type", http.DetectContentType(buffer))
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	resp, err := client.Do(req)

	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var authResp map[string]map[string]interface{}

	data, _ := ioutil.ReadAll(resp.Body)

	_ = json.Unmarshal(data, &authResp)

	viewUrl, _ := authResp["data"]["view_url"].(string)

	return viewUrl, nil
}