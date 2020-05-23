package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"set-flags/models"
)

func AssetInfos(c *gin.Context) {
	assetId := c.Param("id")

	assets := models.FindAssetsByID(assetId)

	c.PureJSON(http.StatusOK, assets)
}
