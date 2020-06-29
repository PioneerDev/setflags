package routers

import (
	"net/http"
	_ "set-flags/docs"
	"set-flags/middleware/jwt"
	"set-flags/pkg/setting"
	v1 "set-flags/routers/api/v1"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

// InitRouter gin router
func InitRouter() *gin.Engine {
	r := gin.New()

	r.Use(gin.Logger())

	r.Use(gin.Recovery())

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"PUT", "POST", "GET"},
		AllowHeaders:     []string{"Origin", "x-user-id", "content-type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return origin == "https://github.com"
		},
		MaxAge: 12 * time.Hour,
	}))

	gin.SetMode(setting.GetConfig().RUNMODE)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	apiv1 := r.Group("")
	apiv1.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})
	apiv1.GET("/auth", v1.Auth)

	apiv1.Use(jwt.JWT())
	{
		apiv1.GET("/flags", v1.ListFlags)
		apiv1.POST("/flag", v1.CreateFlag)
		apiv1.PUT("/flags/:id/:op", v1.UpdateFlag)
		apiv1.GET("/flags/:id/witnesses", v1.GetWitnesses)
		apiv1.GET("/flags/:id/evidences", v1.ListEvidences)
		apiv1.GET("/myflags", v1.FindFlagsByUserID)
		apiv1.POST("/attachments/", v1.UploadEvidence)
		apiv1.GET("/me", v1.Me)
		apiv1.GET("/users/:user_id/rewards/:flag_id", v1.CheckRewards)
		apiv1.GET("/assets", v1.ReadAssets)
		apiv1.GET("/assets/:id", v1.AssetInfos)
	}

	return r
}
