package routers

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	v1 "set-flags/routers/api/v1"
	"time"
)

func InitRouter() *gin.Engine {
	r := gin.New()

	r.Use(gin.Logger())

	r.Use(gin.Recovery())

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"PUT", "POST", "GET"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return origin == "https://github.com"
		},
		MaxAge: 12 * time.Hour,
	}))

	gin.SetMode("debug")

	apiv1 := r.Group("")

	apiv1.GET("/flags", v1.ListFlags)
	apiv1.PUT("/flags/:id/:op", v1.UpdateFlag)

	apiv1.POST("/flags", v1.CreateFlag)
	apiv1.POST("/attachments/:attachment_id", v1.UploadEvidence)

	apiv1.GET("/flags/:flag_id/evidences", v1.ListEvidences)

	apiv1.GET("/users/:user_id/rewards/:flag_id", v1.CheckRewards)

	apiv1.GET("/myflags/:id", v1.FindFlagsByUserID)

	apiv1.GET("/assets/:id", v1.AssetInfos)

	apiv1.GET("/me/:id", v1.Me)

	apiv1.POST("/auth", v1.Auth)

	return r
}





















