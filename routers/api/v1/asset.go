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

// AssetInfos asset info
func AssetInfos(c *gin.Context) {
	code := e.INVALID_PARAMS
	assetID, err := uuid.FromString(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": code,
			"msg":  err.Error(),
			"data": make(map[string]interface{}),
		})
		return
	}

	if !models.ExistAsset(assetID) {
		c.JSON(http.StatusNotFound, gin.H{
			"code": code,
			"msg":  fmt.Sprintf("no asset found."),
			"data": make(map[string]interface{}),
		})
		return
	}

	asset := models.FindAssetByID(assetID)

	code = e.SUCCESS
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": asset,
	})
}

// ReadAssets ReadAssets
func ReadAssets(c *gin.Context) {
	code := e.INVALID_PARAMS

	var pagination schemas.Pagination

	c.ShouldBindQuery(&pagination)

	if pagination.CurrentPage < 1 {
		pagination.CurrentPage = 1
	}

	if pagination.PageSize < 1 {
		pagination.PageSize = setting.GetConfig().App.PageSize
	}

	assets, total := models.ReadAssets(pagination.PageSize, pagination.CurrentPage)

	code = e.SUCCESS
	c.JSON(http.StatusOK, gin.H{
		"code":  code,
		"msg":   e.GetMsg(code),
		"data":  assets,
		"total": total,
	})
}
