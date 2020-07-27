package v1

import (
	"context"
	"fmt"
	"net/http"
	"set-flags/global"
	"set-flags/models"
	"set-flags/pkg/e"
	"set-flags/pkg/logging"
	"set-flags/pkg/setting"
	"set-flags/schemas"
	"strconv"

	"github.com/fox-one/mixin-sdk"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
)

// UploadEvidence upload evidence
// only payer can upload evidence
func UploadEvidence(c *gin.Context) {
	code := e.INVALID_PARAMS

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
		code = e.ERROR_NOT_FOUND_USER
		c.JSON(http.StatusNotFound, gin.H{
			"code": code,
			"msg":  e.GetMsg(code),
			"data": make(map[string]interface{}),
		})
		return
	}

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
		code = e.ERROR_NOT_FOUND_FLAG
		c.JSON(http.StatusNotFound, gin.H{
			"code": code,
			"msg":  e.GetMsg(code),
			"data": make(map[string]interface{}),
		})
		return
	}

	flag := models.FindFlagByID(flagID)

	if flag.PayerID != userID {
		code = e.ERROR_FLAGER_NOT_CURRENT_USER
		c.JSON(http.StatusBadRequest, gin.H{
			"code": code,
			"msg":  e.GetMsg(code),
			"data": make(map[string]interface{}),
		})
		return
	}

	// flag's creator not yet paid.
	if flag.Period == 0 {
		code = e.ERROR_NO_PAID
		c.JSON(http.StatusBadRequest, gin.H{
			"code": code,
			"msg":  e.GetMsg(code),
			"data": make(map[string]interface{}),
		})
		return
	}

	ctx := context.Background()
	attachment, err := global.Bot.CreateAttachment(ctx)
	if err != nil {
		code = e.ERROR
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": code,
			"msg":  err.Error(),
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

	logging.Info("attachmentId: %s, flagId: %s", attachment.AttachmentID, flagID)

	models.CreateEvidence(flagID, attachment.AttachmentID, mediaType, attachment.ViewURL, flag.Period)

	// update flag period status to `done`
	models.UpdateFlagPeriodStatus(flagID, "done")

	code = e.SUCCESS
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  fmt.Sprintf("'%s' uploaded!", fileHeader.Filename),
		"data": make(map[string]interface{}),
	})
}

// ListEvidencesWithPeriod list the evidences with period
func ListEvidencesWithPeriod(c *gin.Context) {

	code := e.INVALID_PARAMS

	var pagination schemas.Pagination

	c.ShouldBindQuery(&pagination)

	if pagination.CurrentPage < 1 {
		pagination.CurrentPage = 1
	}

	if pagination.PageSize < 1 {
		pagination.PageSize = setting.GetConfig().App.PageSize
	}

	flagID, err := uuid.FromString(c.Param("id"))

	logging.Info(fmt.Sprintf("flag_id %v", flagID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": code,
			"msg":  err.Error(),
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

	// 0 means all,
	// -1 missing means current,
	// greater than 0 means specific period
	period, err := strconv.Atoi(c.DefaultQuery("period", "0"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": code,
			"msg":  err.Error(),
			"data": make(map[string]interface{}),
		})
		return
	}

	var evidences []models.Evidence
	var total int

	if period == -1 {
		flag := models.FindFlagByID(flagID)
		evidences, total = models.FindEvidencesByFlagAndPeriod(flagID, pagination.CurrentPage, pagination.PageSize, flag.Period)
	} else if period == 0 {
		evidences, total = models.GetAllEvidenceByFlagID(flagID, pagination.CurrentPage, pagination.PageSize)
	} else if period > 0 {
		evidences, total = models.FindEvidencesByFlagAndPeriod(flagID, pagination.CurrentPage, pagination.PageSize, period)
	}

	code = e.SUCCESS
	c.JSON(http.StatusOK, gin.H{
		"code":  code,
		"msg":   e.GetMsg(code),
		"data":  evidences,
		"total": total,
	})
}
