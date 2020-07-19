package v1

import (
	"net/http"
	"set-flags/models"
	"set-flags/pkg/e"
	"set-flags/pkg/setting"
	"set-flags/schemas"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
)

// GetWitnessesByPeriod GetWitnessesByPeriod
func GetWitnessesByPeriod(c *gin.Context) {
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
	period, err := strconv.Atoi(c.DefaultQuery("period", "-1"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": code,
			"msg":  err.Error(),
			"data": make(map[string]interface{}),
		})
		return
	}

	var witnesses []schemas.WitnessSchema
	var total int

	if period == -1 {
		flag := models.FindFlagByID(flagID)
		witnesses, total = models.GetWitnessWithPeriod(flagID, pagination.PageSize, pagination.CurrentPage, flag.Period)
	} else if period == 0 {
		witnesses, total = models.GetAllWitnessByFlagID(flagID, pagination.PageSize, pagination.CurrentPage)
	} else if period > 0 {
		witnesses, total = models.GetWitnessWithPeriod(flagID, pagination.PageSize, pagination.CurrentPage, period)
	}

	code = e.SUCCESS
	c.JSON(http.StatusOK, gin.H{
		"code":  code,
		"msg":   e.GetMsg(code),
		"data":  witnesses,
		"total": total,
	})
}
