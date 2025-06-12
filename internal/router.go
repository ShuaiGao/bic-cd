package internal

import (
	"bic-cd/internal/admin"
	"bic-cd/internal/manager"
	"bic-cd/pkg/config"
	"bic-cd/pkg/gen/api"
	"bic-cd/pkg/gin_recover"
	"bic-cd/pkg/log"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
)

func Cors() gin.HandlerFunc {
	return func(context *gin.Context) {
		method := context.Request.Method
		context.Header("Access-Control-Allow-Origin", "*")
		context.Header("Access-Control-Allow-Headers", "Content-Type,X-CSRF-Token,Authorization,x-request-id,X-API-KEY")
		context.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, DELETE, PATCH, PUT")
		context.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		context.Header("Access-Control-Allow-Credentials", "true")

		if method == "OPTIONS" {
			context.AbortWithStatus(http.StatusNoContent)
		}
	}
}

func Setup() *gin.Engine {
	gin.SetMode(config.GlobalConf.App.RunMode)
	g := gin.New()
	g.Use(gin_recover.Recovery())
	g.Use(log.MidLogger())
	g.Use(Cors())
	g.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(200, "success")
	})
	if gin.IsDebugging() {
		g.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}
	apiNoJwt := g.Group("")
	api.RegisterAuthServiceHttpHandler(apiNoJwt, admin.Auth{})
	api.RegisterManagerServiceHttpHandler(apiNoJwt, manager.Manager{})
	//注册、设置路由
	return g
}
