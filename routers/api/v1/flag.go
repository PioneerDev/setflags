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

	if !models.UserExist(userID) {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": code,
			"msg":  fmt.Sprintf("not found current user."),
			"data": make(map[string]interface{}),
		})
		return
	}

	var pagination schemas.Pagination

	c.ShouldBindQuery(&pagination)

	if pagination.CurrentPage < 1 {
		pagination.CurrentPage = 1
	}

	if pagination.PageSize < 1 {
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

	if flag.MaxWitness <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": code,
			"msg":  fmt.Sprintf("max witness must greater than zero."),
			"data": make(map[string]interface{}),
		})
		return
	}

	if flag.Amount <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": code,
			"msg":  fmt.Sprintf("amount must greater than zero."),
			"data": make(map[string]interface{}),
		})
		return
	}

	if flag.DaysPerPeriod <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": code,
			"msg":  fmt.Sprintf("days per period must greater than zero."),
			"data": make(map[string]interface{}),
		})
		return
	}

	if flag.TotalPeriod <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": code,
			"msg":  fmt.Sprintf("total period must greater than zero."),
			"data": make(map[string]interface{}),
		})
		return
	}

	// check asset id
	if !models.ExistAsset(flag.AssetID) {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": code,
			"msg":  fmt.Sprintf("asset not exist, asset_id: %s, symbol: %s", flag.AssetID, flag.Symbol),
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
		"flag_id":   flagID,
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
	if !models.UserExist(userID) {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": code,
			"msg":  fmt.Sprintf("not found current user."),
			"data": make(map[string]interface{}),
		})
		return
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
		models.UpdateFlagPeriodStatus(flagID, op)
	} else if flag.PayerID != userID && (op == "yes" || op == "no") {
		err := models.UpsertWitness(flagID, userID, flag.AssetID, op, flag.Symbol, flag.Period, flag.MaxWitness)
		if err != nil {
			code = e.ERROR
			c.JSON(http.StatusOK, gin.H{
				"code": code,
				"msg":  err.Error(),
				"data": make(map[string]interface{}),
			})
			return
		}
		code = e.SUCCESS
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

	var pagination schemas.Pagination

	c.ShouldBindQuery(&pagination)

	if pagination.CurrentPage < 1 {
		pagination.CurrentPage = 1
	}

	if pagination.PageSize < 1 {
		pagination.PageSize = setting.GetConfig().App.PageSize
	}

	flags, total := models.FindFlagsByUserID(userID, pagination.CurrentPage, pagination.PageSize)

	code = e.SUCCESS
	c.JSON(http.StatusOK, gin.H{
		"code":  code,
		"msg":   e.GetMsg(code),
		"data":  flags,
		"total": total,
	})
}

// FlagDetail FlagDetail
func FlagDetail(c *gin.Context) {
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

	flagSchema := schemas.FlagSchema{
		ID:              flag.ID,
		PayerID:         flag.PayerID,
		PayerName:       flag.PayerName,
		PayerAvatarURL:  flag.PayerAvatarURL,
		Task:            flag.Task,
		Days:            flag.Days,
		MaxWitness:      flag.MaxWitness,
		AssetID:         flag.AssetID,
		Symbol:          flag.Symbol,
		Amount:          flag.Amount,
		TimesAchieved:   flag.TimesAchieved,
		Status:          flag.Status,
		DaysPerPeriod:   flag.DaysPerPeriod,
		PeriodStatus:    flag.PeriodStatus,
		Verified:        "UNSET",
		RemainingAmount: flag.RemainingAmount,
		RemainingDays:   flag.RemainingDays,
		Period:          flag.Period,
		TotalPeriod:     flag.TotalPeriod,
	}
	// current user is not flag creator
	// fetch witness
	if userID != flag.PayerID {
		witness := models.GetWitnessByFlagIDAndPayeeID(flagID, userID, flag.Period)
		if witness.Verified != "" {
			flagSchema.Verified = witness.Verified
		}
	}

	code = e.SUCCESS
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": flagSchema,
	})
}
