package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"set-flags/models"
	"set-flags/pkg/e"
)

func AssetInfos(c *gin.Context) {
	code := e.INVALID_PARAMS
	assetId := c.Param("id")

	asset := models.FindAssetByID(assetId)

	code = e.SUCCESS
	c.JSON(http.StatusOK, gin.H{
		"code" : code,
		"msg" : e.GetMsg(code),
		"data" : asset,
	})
}
