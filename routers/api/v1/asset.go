package v1

import (
	"net/http"
	"set-flags/models"
	"set-flags/pkg/e"

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

	asset := models.FindAssetByID(assetID)

	code = e.SUCCESS
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": asset,
	})
}
