package internal

import (
	"bic-cd/internal/admin"
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

var ginInstance *gin.Engine

func Setup() {
	ginInstance = gin.New()
	ginInstance.ContextWithFallback = true
	ginInstance.Use(gin_recover.Recovery())
	ginInstance.Use(log.MidLogger())
	ginInstance.Use(Cors())
	ginInstance.GET("/health/", func(ctx *gin.Context) {
		ctx.JSON(200, "success")
	})
	ginInstance.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(200, "success")
	})
}

// SetupRouter 启动路由，添加中间件
func SetupRouter() *gin.Engine {
	g := ginInstance
	apiNoJwt := g.Group("")
	api.RegisterAuthServiceHttpHandler(apiNoJwt, admin.Auth{})
	if gin.IsDebugging() {
		g.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}
	//注册、设置路由
	return g
}
