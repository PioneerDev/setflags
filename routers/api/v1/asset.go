package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"set-flags/models"
)

func AssetInfos(c *gin.Context) {
	assetId := c.Param("id")

	asset := models.FindAssetByID(assetId)

	c.JSON(http.StatusOK, gin.H{
		"code" : 200,
		"msg" : "ok",
		"data" : asset,
	})
}
