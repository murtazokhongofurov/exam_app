package api

import (
	_ "exam/api-gateway/api/docs"
	v1 "exam/api-gateway/api/handlers/v1"
	"exam/api-gateway/api/middleware"
	token "exam/api-gateway/api/tokens"
	"exam/api-gateway/config"
	"exam/api-gateway/pkg/logger"
	"exam/api-gateway/services"
	"exam/api-gateway/storage/repo"

	"github.com/gin-contrib/cors"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	swaggerfile "github.com/swaggo/files"

	ginSwagger "github.com/swaggo/gin-swagger"
)

type Option struct {
	Conf           config.Config
	Logger         logger.Logger
	ServiceManager services.IServiceManager
	Redis          repo.RedisRepo
	CasbinEnforcer *casbin.Enforcer
}

// New ...
// @title           exam api
// @version         2.0
// @description     This is exam server api server
// @termsOfService  2 term exam

// @contact.name   Murtazoxon
// @contact.url    https://t.me/murtazokhon_gofurov
// @contact.email  gofurovmurtazoxon@gmail.com

// @host    	   localhost:9090

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func New(option Option) *gin.Engine {
	router := gin.New()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	jwtHandler := token.JWTHandler{
		SigninKey: option.Conf.SignInKey,
		Log:       option.Logger,
	}

	handlerV1 := v1.New(&v1.HandlerV1Config{
		Logger:         option.Logger,
		ServiceManager: option.ServiceManager,
		Cfg:            option.Conf,
		Redis:          option.Redis,
		JWTHandler:     jwtHandler,
	})

	corConfig := cors.DefaultConfig()
	corConfig.AllowAllOrigins = true
	corConfig.AllowCredentials = true
	corConfig.AllowHeaders = []string{"*"}
	corConfig.AllowBrowserExtensions = true
	corConfig.AllowMethods = []string{"*"}
	router.Use(cors.New(corConfig))
	router.Use(middleware.NewAuth(option.CasbinEnforcer, jwtHandler, config.Load()))

	api := router.Group("/v1")
	// customer service post, get...
	api.GET("/customer/post/review/:id", handlerV1.GetCustomerById)
	api.PUT("/customer/update", handlerV1.UpdateCustomer)
	api.DELETE("/customer/:id", handlerV1.DeleteCustomer)
	api.GET("/customers/:page/:limit/:order/:search", handlerV1.GetCustomerBySearchOrder)

	//post service post, get...
	api.POST("/post", handlerV1.CreatePost)
	api.GET("/post/get/:id", handlerV1.GetPostReview)
	api.PUT("/post/update", handlerV1.UpdatePost)
	api.DELETE("/post/delete/:id", handlerV1.DeletePost)
	api.GET("/post/list/:page/:limit/:search", handlerV1.ListPost)

	//review service post, get...
	// api.POST("/review", handlerV1.CreateReview)
	api.GET("/review/get/:id", handlerV1.GetReviewById)
	api.PUT("/review/update", handlerV1.UpdateReview)
	api.DELETE("/review/delete/:id", handlerV1.DeleteReview)

	// register customer
	api.POST("/customer/register", handlerV1.RegisterCustomer)
	api.GET("/verify/:email/:code", handlerV1.Verify)
	api.GET("/login/:email/:password", handlerV1.Login)

	api.GET("/admin/login/:admin_name/:password", handlerV1.LoginAdmin)
	api.GET("/moderator/login/:name/:password", handlerV1.LoginModerator)

	url := ginSwagger.URL("swagger/doc.json")
	api.GET("swagger/*any", ginSwagger.WrapHandler(swaggerfile.Handler, url))

	return router
}
