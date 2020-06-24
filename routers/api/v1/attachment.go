package v1

import (
	"context"
	"fmt"
	"net/http"
	"set-flags/models"
	"set-flags/pkg/e"
	"set-flags/pkg/setting"
	"set-flags/schemas"

	"github.com/fox-one/mixin-sdk"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
)

// UploadEvidence upload evidence
func UploadEvidence(c *gin.Context) {
	code := e.INVALID_PARAMS

	// check document type
	mediaType := c.Query("type")
	if mediaType != "image" && mediaType != "audio" && mediaType != "video" && mediaType != "document" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": code,
			"msg":  fmt.Sprintf("type: %s is invalid.", mediaType),
			"data": make(map[string]interface{}),
		})
		return
	}

	// check flag id
	flagID, err := uuid.FromString(c.Query("flag_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": code,
			"msg":  err.Error(),
			"data": make(map[string]interface{}),
		})
		return
	}

	// check flag exist.
	if !models.FlagExists(flagID) {
		c.JSON(http.StatusNotFound, gin.H{
			"code": code,
			"msg":  "not found specific flag.",
			"data": make(map[string]interface{}),
		})
		return
	}

	// check user id
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

	// check flag exist.
	if !models.UserExist(userID) {
		c.JSON(http.StatusNotFound, gin.H{
			"code": code,
			"msg":  "not found specific user.",
			"data": make(map[string]interface{}),
		})
		return
	}

	// upload attachment
	user, err := mixin.NewUser(setting.ClientID.String(), setting.SessionID, setting.SessionKey, setting.PINToken)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  err.Error(),
			"data": make(map[string]interface{}),
		})
		return
	}

	ctx := context.Background()
	attachment, err := user.CreateAttachment(ctx)

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

	f, err := fileHeader.Open()
	if err != nil {
		code = e.ERROR
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": code,
			"msg":  err.Error(),
			"data": make(map[string]interface{}),
		})
		return
	}

	buffer := make([]byte, fileHeader.Size)
	_, err = f.Read(buffer)

	if err != nil {
		code = e.ERROR
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": code,
			"msg":  err.Error(),
			"data": make(map[string]interface{}),
		})
		return
	}

	err = mixin.UploadAttachment(ctx, attachment, buffer)
	if err != nil {
		code = e.ERROR
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": code,
			"msg":  err.Error(),
			"data": make(map[string]interface{}),
		})
		return
	}

	fmt.Printf("attachmentId: %s, flagId: %s", attachment.AttachmentID, flagID)

	models.CreateEvidence(flagID, attachment.AttachmentID, mediaType, attachment.ViewURL)

	// update flag status to `done`
	models.UpdateFlagStatus(flagID, "done")

	code = e.SUCCESS
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  fmt.Sprintf("'%s' uploaded!", fileHeader.Filename),
		"data": make(map[string]interface{}),
	})
}
