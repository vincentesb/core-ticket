package main

import (
	"core-ticket/base/helpers/gin_helper"
	"core-ticket/base/helpers/http_helper"
	"core-ticket/base/helpers/struct_validator"
	"core-ticket/helpers/shutdown"
	"core-ticket/middlewares"
	"core-ticket/modules/common/healthcheck"
	"core-ticket/modules/ip_whitelist"
	"core-ticket/modules/return_ticket"
	"time"

	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/jmoiron/sqlx"
	"gopkg.in/guregu/null.v4"
)

func SetupRouter(db map[string]*sqlx.DB, shutdownMgr *shutdown.Manager) *gin.Engine {

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(sentrygin.New(sentrygin.Options{Repanic: true}))

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"*"}
	config.AllowCredentials = true
	router.Use(cors.New(config))

	validate := GetValidator()
	validate.RegisterCustomTypeFunc(
		struct_validator.NullFlake,
		null.String{},
		null.Int{},
		null.Float{},
		null.Bool{},
	)
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("alphanum_space", struct_validator.ValidateAlphaNumericSpaceOnly)
		v.RegisterValidation("not_emoji", struct_validator.NotEmoji())
	}

	// Simple health check endpoint without authentication or database dependency
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":     "ok",
			"message":    "Application is running",
			"version":    Version,
			"commit":     CommitSHA,
			"build_time": BuildTime,
			"timestamp":  time.Now().UTC().Format(time.RFC3339),
		})
	})

	// Panic recovery in middleware
	router.Use(middlewares.PanicMiddleware())
	ipWhitelistService, _ := ip_whitelist.InitializeIpWhitelistService(db)
	middleware := middlewares.NewMiddleware(ipWhitelistService)

	// Custom gin engine
	routerEngine := gin_helper.Engine(router, db, middleware)
	{
		routerEngine.RegisterRouter(healthcheck.Router)
		routerEngine.RegisterRouter(return_ticket.Router)
	}

	router.NoRoute(func(c *gin.Context) {
		http_helper.NotFoundResponse(c, "Page Not Found", nil)
	})

	router.HandleMethodNotAllowed = true

	return router
}
