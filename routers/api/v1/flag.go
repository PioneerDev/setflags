package v1

import (
	"fmt"
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

	var pagination schemas.Pagination

	c.ShouldBindQuery(&pagination)

	if pagination.CurrentPage == 0 {
		pagination.CurrentPage = 1
	}

	if pagination.PageSize == 0 {
		pagination.PageSize = setting.PageSize
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

	witnesses := models.GetWitnesses(flagID, pagination.CurrentPage, pagination.PageSize)

	code = e.SUCCESS
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": witnesses,
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
