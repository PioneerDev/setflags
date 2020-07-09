package v1

import (
	"fmt"
	"net/http"
	"set-flags/models"
	"set-flags/pkg/e"
	"set-flags/pkg/logging"
	"set-flags/pkg/setting"
	"set-flags/schemas"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
)

// ListFlags list all the flags
// @Summary 获取全部Flag
// @Produce  json
// @Param current_page query int false "CurrentPage"
// @Param page_size query int false "PageSize"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"ok"}"
// @Router /flags [get]
func ListFlags(c *gin.Context) {
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

	var pagination schemas.Pagination

	c.ShouldBindQuery(&pagination)

	if pagination.CurrentPage == 0 {
		pagination.CurrentPage = 1
	}

	if pagination.PageSize == 0 {
		pagination.PageSize = setting.GetConfig().App.PageSize
	}

	// data, total := models.GetAllFlags(pagination.PageSize, pagination.CurrentPage)
	data, total := models.GetFlagsWithVerified(pagination.PageSize, pagination.CurrentPage, userID)

	code = e.SUCCESS
	c.JSON(http.StatusOK, gin.H{
		"code":  code,
		"msg":   e.GetMsg(code),
		"data":  data,
		"total": total,
	})
}

// CreateFlag create a flag
// @Summary 创建Flag
// @Produce json
// @Param payer_name body string false "创建者的name"
// @Param task body string true "任务名称"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"ok"}"
// @Router /flag [post]
func CreateFlag(c *gin.Context) {
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
		c.JSON(http.StatusBadRequest, gin.H{
			"code": code,
			"msg":  fmt.Sprintf("not found current user."),
			"data": make(map[string]interface{}),
		})
		return
	}

	var flag schemas.FlagSchema

	if err := c.ShouldBindJSON(&flag); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": code,
			"msg":  err.Error(),
			"data": make(map[string]interface{}),
		})
		return
	}

	// find user
	user := models.FindUserByID(userID)
	flag.PayerID = userID
	traceID, _ := uuid.NewV1()

	flagID := models.CreateFlag(&flag, user)

	assetID := flag.AssetID.String()
	assetID = "965e5c6e-434c-3fa9-b780-c50f43cd955c"
	memo := "转账给励志机器人."
	appID := setting.GetConfig().Bot.ClientID.String()

	payment := models.Payment{
		TraceID:    traceID,
		FlagID:     flagID,
		AssetID:    assetID,
		OpponentID: appID,
		Amount:     fmt.Sprintf("%f", flag.Amount),
		Memo:       memo,
	}

	models.CreatePayment(payment)

	data := map[string]interface{}{
		"recipient": appID,
		"asset":     assetID,
		"amount":    flag.Amount,
		"trace":     traceID,
		"memo":      memo,
	}

	code = e.SUCCESS
	c.JSON(http.StatusCreated, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": data,
	})
}

// UpdateFlag Update an existing flag
// only witness can update flag
func UpdateFlag(c *gin.Context) {
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

	flagID, err := uuid.FromString(c.Param("id"))
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
	flag := models.FindFlagByID(flagID)

	op := c.Param("op")

	if flag.PayerID == userID && op == "done" {
		code = e.SUCCESS
		models.UpdateFlagStatus(flagID, op)
	} else if flag.PayerID != userID && (op == "yes" || op == "no") {
		code = e.SUCCESS
		models.UpsertWitness(flagID, userID, op)
	}

	if code != e.SUCCESS {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": code,
			"msg":  "invalid opreration",
			"data": make(map[string]interface{}),
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"msg":  e.GetMsg(code),
			"data": make(map[string]interface{}),
		})
	}
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
		pagination.PageSize = setting.GetConfig().App.PageSize
	}

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

	flags, total := models.FindFlagsByUserID(userID, pagination.CurrentPage, pagination.PageSize)

	code = e.SUCCESS
	c.JSON(http.StatusOK, gin.H{
		"code":  code,
		"msg":   e.GetMsg(code),
		"data":  flags,
		"total": total,
	})
}

// GetWitnesses list all witnesses of the flag
func GetWitnesses(c *gin.Context) {
	code := e.INVALID_PARAMS

	var pagination schemas.Pagination

	c.ShouldBindQuery(&pagination)

	if pagination.CurrentPage == 0 {
		pagination.CurrentPage = 1
	}

	if pagination.PageSize == 0 {
		pagination.PageSize = setting.GetConfig().App.PageSize
	}

	flagID, err := uuid.FromString(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": code,
			"msg":  err.Error(),
			"data": make(map[string]interface{}),
		})
		return
	}

	witnesses, total := models.GetWitnessSchema(flagID, pagination.PageSize, pagination.CurrentPage)

	code = e.SUCCESS
	c.JSON(http.StatusOK, gin.H{
		"code":  code,
		"msg":   e.GetMsg(code),
		"data":  witnesses,
		"total": total,
	})
}

// ListEvidences list all the evidences since yesterday
func ListEvidences(c *gin.Context) {

	code := e.INVALID_PARAMS

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

	var pagination schemas.Pagination

	c.ShouldBindQuery(&pagination)

	if pagination.CurrentPage == 0 {
		pagination.CurrentPage = 1
	}

	if pagination.PageSize == 0 {
		pagination.PageSize = setting.GetConfig().App.PageSize
	}

	data, total := models.FindEvidencesByFlag(flagID, pagination.CurrentPage, pagination.PageSize)

	code = e.SUCCESS
	c.JSON(http.StatusOK, gin.H{
		"code":  code,
		"msg":   e.GetMsg(code),
		"data":  data,
		"total": total,
	})
}

// FlagDetail FlagDetail
func FlagDetail(c *gin.Context) {
	code := e.INVALID_PARAMS

	flagID, err := uuid.FromString(c.Query("flag_id"))

	logging.Info(fmt.Sprintf("flag_id %v", flagID))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": code,
			"msg":  err.Error(),
			"data": make(map[string]interface{}),
		})
		return
	}

	flag := models.FindFlagByID(flagID)

	code = e.SUCCESS
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": flag,
	})
}
