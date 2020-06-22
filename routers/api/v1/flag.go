package v1

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"set-flags/models"
	"set-flags/pkg/e"
	"set-flags/pkg/setting"
	"set-flags/schemas"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
)

// ListFlags list all the flags
func ListFlags(c *gin.Context) {
	code := e.INVALID_PARAMS

	var pagination schemas.Pagination

	c.ShouldBindQuery(&pagination)

	if pagination.CurrentPage == 0 {
		pagination.CurrentPage = 1
	}

	if pagination.PageSize == 0 {
		pagination.PageSize = setting.PageSize
	}

	data := models.GetAllFlags(pagination.PageSize, pagination.CurrentPage)

	code = e.SUCCESS
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": data,
	})
}

// CreateFlag create a flag
func CreateFlag(c *gin.Context) {
	code := e.INVALID_PARAMS

	var flag schemas.Flag

	if err := c.ShouldBindJSON(&flag); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": code,
			"msg":  err.Error(),
			"data": make(map[string]interface{}),
		})
		return
	}

	// find user
	user := models.FindUserByID(flag.PayerID)

	models.CreateFlag(&flag, user)

	code = e.SUCCESS
	c.JSON(http.StatusCreated, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": make(map[string]interface{}),
	})
}

// UpdateFlag Update an existing flag
func UpdateFlag(c *gin.Context) {
	code := e.INVALID_PARAMS

	flagID, err := uuid.FromString(c.Param("id"))
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

	if !models.FlagExists(flagID) {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "Flag not found.",
			"data": make(map[string]interface{}),
		})
		return
	}

	flag := models.FindFlagByID(flagID)

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
// FindFlagsByUserID list all flags of the user
func FindFlagsByUserID(c *gin.Context) {
	code := e.INVALID_PARAMS

	var pagination schemas.Pagination

	c.ShouldBindQuery(&pagination)

	if pagination.CurrentPage == 0 {
		pagination.CurrentPage = 1
	}

	if pagination.PageSize == 0 {
		pagination.PageSize = setting.PageSize
	}

	userID, err := uuid.FromString(c.GetHeader("x-user-id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": code,
			"msg":  err.Error(),
			"data": make(map[string]interface{}),
		})
		return
	}

	flags := models.FindFlagsByUserID(userID, pagination.CurrentPage, pagination.PageSize)

	code = e.SUCCESS
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": flags,
	})
}


// GetWitnesses list all witnesses of the flag
func GetWitnesses(c *gin.Context) {
	code := e.INVALID_PARAMS

	flagID, err := uuid.FromString(c.Param("id"))

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

	witnesses := models.GetWitnesses(flagID, currentPage, pageSize)

	code = e.SUCCESS
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": witnesses,
	})
}

// UploadEvidence upload evidence
func UploadEvidence(c *gin.Context) {
	code := e.INVALID_PARAMS
	userID := c.GetHeader("user_id")

	attachmentID, err := uuid.FromString(c.Param("attachment_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": code,
			"msg":  err.Error(),
			"data": make(map[string]interface{}),
		})
		return
	}
	flagID, err := uuid.FromString(c.Query("flag_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": code,
			"msg":  err.Error(),
			"data": make(map[string]interface{}),
		})
		return
	}

	fmt.Printf("attachmentId: %s, flagId: %s", attachmentID, flagID)

	mediaType := c.Query("type")

	if mediaType != "image" && mediaType != "audio" && mediaType != "video" && mediaType != "document" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": code,
			"msg":  fmt.Sprintf("type: %s is invalid.", mediaType),
			"data": make(map[string]interface{}),
		})
		return
	}

	if !models.FlagExists(flagID) {
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
	accessToken, _ := models.FindUserToken(userID)

	viewURL, err := UploadAttachment(client, fileHeader, accessToken)
	if err != nil {
		code = e.ERROR_UPLOAD_ATTACHMENT
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"msg":  e.GetMsg(code),
			"data": make(map[string]interface{}),
		})
		return
	}

	models.CreateEvidence(attachmentID, flagID, mediaType, viewURL)

	// update flag status to `done`
	models.UpdateFlagStatus(flagID, "done")

	code = e.SUCCESS
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  fmt.Sprintf("'%s' uploaded!", fileHeader.Filename),
		"data": make(map[string]interface{}),
	})
}

// ListEvidences list all the evidences since yesterday
func ListEvidences(c *gin.Context) {

	code := e.INVALID_PARAMS

	flagID, _ := uuid.FromString(c.Param("flag_id"))

	var pagination schemas.Pagination

	c.ShouldBindQuery(&pagination)

	if pagination.CurrentPage == 0 {
		pagination.CurrentPage = 1
	}

	if pagination.PageSize == 0 {
		pagination.PageSize = setting.PageSize
	}

	data := models.FindEvidencesByFlag(flagID, pagination.CurrentPage, setting.PageSize)

	code = e.SUCCESS
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": data,
	})
}

// UploadAttachment https://developers.mixin.one/api/l-messages/create-attachment/
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

	viewURL, _ := authResp["data"]["view_url"].(string)

	return viewURL, nil
}
