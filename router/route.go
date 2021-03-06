package router

import (
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/go1234.cn/gin_scaffold/controller"
	"github.com/go1234.cn/gin_scaffold/docs"
	"github.com/go1234.cn/gin_scaffold/golang_common/lib"
	"github.com/go1234.cn/gin_scaffold/middleware"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"log"
)

// @title Swagger Example API
// @version 1.0
// @description This is a sample server celler server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1
// @query.collection.format multi

// @securityDefinitions.basic BasicAuth

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

// @securitydefinitions.oauth2.application OAuth2Application
// @tokenUrl https://example.com/oauth/token
// @scope.write Grants write access
// @scope.admin Grants read and write access to administrative information

// @securitydefinitions.oauth2.implicit OAuth2Implicit
// @authorizationurl https://example.com/oauth/authorize
// @scope.write Grants write access
// @scope.admin Grants read and write access to administrative information

// @securitydefinitions.oauth2.password OAuth2Password
// @tokenUrl https://example.com/oauth/token
// @scope.read Grants read access
// @scope.write Grants write access
// @scope.admin Grants read and write access to administrative information

// @securitydefinitions.oauth2.accessCode OAuth2AccessCode
// @tokenUrl https://example.com/oauth/token
// @authorizationurl https://example.com/oauth/authorize
// @scope.admin Grants read and write access to administrative information

// @x-extension-openapi {"example": "value on a json format"}

func InitRouter(middlewares ...gin.HandlerFunc) *gin.Engine {
	// programatically set swagger info
	docs.SwaggerInfo.Title = lib.GetStringConf("base.swagger.title")
	docs.SwaggerInfo.Description = lib.GetStringConf("base.swagger.desc")
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = lib.GetStringConf("base.swagger.host")
	docs.SwaggerInfo.BasePath = lib.GetStringConf("base.swagger.base_path")
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	router := gin.Default()
	router.Use(middlewares...)
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	//设置admin login的路由及中间件
	adminLoginRouter := router.Group("/admin_login")
	//使用redis保存session
	store, err := sessions.NewRedisStore(10, "tcp", "localhost:6379", "", []byte("secret"))
	if err != nil {
		log.Fatalf("session.NewRedisStore err:%v", err)
	}
	//设置router的中间件
	adminLoginRouter.Use(
		sessions.Sessions("mysession", store), //设置session的中间件
		middleware.RecoveryMiddleware(),
		middleware.RequestLog(),
		middleware.TranslationMiddleware(),
	)
	{
		//子路由的注册
		controller.AdminRegister(adminLoginRouter)
	}

	adminRouter := router.Group("/admin")
	//设置router的中间件
	adminRouter.Use(
		sessions.Sessions("mysession", store), //设置session的中间件
		middleware.RecoveryMiddleware(),
		middleware.RequestLog(),
		middleware.SessionAuthMiddleware(), //权限验证
		middleware.TranslationMiddleware(),
	)
	{
		//子路由的注册
		controller.AdminsRegister(adminRouter)
	}

	//服务
	serviceRouter := router.Group("/service")
	//设置router的中间件
	serviceRouter.Use(
		sessions.Sessions("mysession", store), //设置session的中间件
		middleware.RecoveryMiddleware(),
		middleware.RequestLog(),
		middleware.SessionAuthMiddleware(), //权限验证
		middleware.TranslationMiddleware(),
	)
	{
		//子路由的注册
		controller.ServiceRegister(serviceRouter)
	}

	//dashboard大盘
	dashboardRouter := router.Group("/dashboard")
	//设置router的中间件
	serviceRouter.Use(
		sessions.Sessions("mysession", store), //设置session的中间件
		middleware.RecoveryMiddleware(),
		middleware.RequestLog(),
		middleware.SessionAuthMiddleware(), //权限验证
		middleware.TranslationMiddleware(),
	)
	{
		//子路由的注册
		controller.ServiceRegister(dashboardRouter)
	}
	return router
}
