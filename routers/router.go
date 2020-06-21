package routers

import (
	"net/http"
	"set-flags/pkg/setting"
	v1 "set-flags/routers/api/v1"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// InitRouter gin router
func InitRouter() *gin.Engine {
	r := gin.New()

	r.Use(gin.Logger())

	r.Use(gin.Recovery())

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"PUT", "POST", "GET"},
		AllowHeaders:     []string{"Origin", "x-user-id", "content-type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return origin == "https://github.com"
		},
		MaxAge: 12 * time.Hour,
	}))

	gin.SetMode(setting.RunMode)

	apiv1 := r.Group("")

	apiv1.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	apiv1.GET("/flags", v1.ListFlags)
	apiv1.POST("/flag", v1.CreateFlag)
	apiv1.PUT("/flags/:id/:op", v1.UpdateFlag)
	apiv1.GET("/flags/:id", v1.GetWitnesses)
	apiv1.GET("/myflags", v1.FindFlagsByUserID)
	apiv1.POST("/attachments/:attachment_id", v1.UploadEvidence)
	apiv1.GET("/flags/:flag_id/evidences", v1.ListEvidences)
	apiv1.GET("/me", v1.Me)
	apiv1.POST("/auth", v1.Auth)
	apiv1.GET("/users/:user_id/rewards/:flag_id", v1.CheckRewards)
	apiv1.GET("/assets/:id", v1.AssetInfos)

	return r
}
