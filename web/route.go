package web

import (
	_ "github.com/cloverzrg/go-portforward/docs"
	"github.com/cloverzrg/go-portforward/web/controller"
	"github.com/gin-gonic/gin"
)

func SetupRoute(r *gin.Engine) {
	g := r.Group("/v1")
	setupDefaultRoute(g)
}

func setupDefaultRoute(r *gin.RouterGroup) {
	r.GET("/network/interfaces", controller.GetNetworkInterfaces)

	// add a forward and start it
	r.POST("/forward/", controller.AddForward)
	// get forward list
	r.GET("/forward/", controller.AddForward)
	// stop forward by id
	r.POST("/forward/:id/stop", controller.StopForward)
	// start forward by id
	r.POST("/forward/:id/start", controller.StartForward)
	// get froward detail by id
	r.GET("/forward/:id")
	// delete forward by id
	r.POST("/forward/:id/delete", controller.DeleteForward)

	// restart all forward
	r.POST("/forward-manager/restart", controller.AddForward)
}
