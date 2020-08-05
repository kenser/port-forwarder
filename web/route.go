package web

import (
	_ "github.com/cloverzrg/go-portforwarder/docs"
	"github.com/cloverzrg/go-portforwarder/web/controller"
	"github.com/gin-gonic/gin"
)

func SetupRoute(r *gin.Engine) {
	g := r.Group("/api")
	setupDefaultRoute(g)
}

func setupDefaultRoute(r *gin.RouterGroup) {
	r.GET("/network/interfaces", controller.GetNetworkInterfaces)

	// add a forward and start it
	r.POST("/forward/", controller.AddForward)
	// get forward list
	r.GET("/forward/", controller.GetForwardList)
	// stop forward by id
	r.POST("/forward/:id/stop", controller.StopForward)
	// start forward by id
	r.POST("/forward/:id/start", controller.StartForward)
	// get froward detail by id
	r.GET("/forward/:id", controller.GetForwardById)
	// delete forward by id
	r.POST("/forward/:id/delete", controller.DeleteForward)

	// restart all forward
	r.POST("/forward-manager/restart")
}
