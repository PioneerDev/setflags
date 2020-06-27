package v1

import (
	"context"
	"net/http"
	"set-flags/models"
	"set-flags/pkg/e"
	"set-flags/pkg/setting"
	"set-flags/schemas"

	"github.com/fox-one/mixin-sdk"
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

// ReadUserAssets ReadUserAssets
func ReadUserAssets(c *gin.Context) {
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

	accessToken := models.FindUserAccessToken(userID)
	ctx := context.Background()
	assets, err := mixin.ReadAssets(ctx, accessToken)

	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"code": 400,
			"msg":  err.Error(),
			"data": make(map[string]interface{}),
		})
		return
	}

	assetInfos := make([]schemas.AssetSchema, 0, 0)

	for _, a := range assets {

		priceUSD, _ := a.PriceUsd.Float64()
		balance, _ := a.Balance.Float64()
		assetInfos = append(assetInfos, schemas.AssetSchema{
			AssetID:  a.AssetID,
			Symbol:   a.Symbol,
			PriceUSD: priceUSD,
			Balance:  balance,
		})
	}

	code = e.SUCCESS
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": assetInfos,
	})
}
