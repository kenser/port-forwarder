package web

import (
	"github.com/cloverzrg/go-portforwarder/config"
	"github.com/cloverzrg/go-portforwarder/constants"
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

func Start() (err error) {
	r := gin.Default()
	ENV := config.Config.ENV
	if ENV != constants.ENV_PROD {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	//r.Use(gin.LoggerWithWriter(logger.Entry.Writer()))
	//r.Use(gin.RecoveryWithWriter(logger.Entry.WriterLevel(logrus.ErrorLevel)))
	setupSwagger(r)
	SetupRoute(r)
	err = r.Run(config.Config.HTTP.Listen)
	if err != nil {
		return
	}
	return
}

func setupSwagger(r *gin.Engine) {
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
