package controller

import (
	"github.com/cloverzrg/go-portforward/logger"
	"github.com/cloverzrg/go-portforward/service/forward"
	"github.com/cloverzrg/go-portforward/utils"
	"github.com/cloverzrg/go-portforward/web/dto"
	"github.com/cloverzrg/go-portforward/web/resp"
	"github.com/gin-gonic/gin"
	"net"
)

// @Summary get network interface list
// @Description ""
// @Tags network
// @Produce  json
// @Success 200 {object} resp.DataResp{data=[]dto.NetworkInterface}
// @Router /v1/network/interfaces [get]
func GetNetworkInterfaces(c *gin.Context) {
	defaultGateway := utils.GetLocalHostAddress()
	list := []dto.NetworkInterface{
		{"", "all", false},
		{"0.0.0.0", "all ipv4", false},
		{"127.0.0.1", "ipv4 localhost", false},
	}
	ifaces, err := net.Interfaces()
	if err != nil {
		c.JSON(resp.UnexpectedError(err))
		return
	}
	for _, v := range ifaces {
		item, _ := v.Addrs()
		for _, addr := range item {
			switch vv := addr.(type) {
			case *net.IPNet:
				if !vv.IP.IsLoopback() {
					if vv.IP.To4() != nil {
						isDefaultGateway := false
						if defaultGateway == vv.IP.String() {
							isDefaultGateway = true
						}
						list = append(list, dto.NetworkInterface{
							Address:        vv.IP.String(),
							Desc:           v.Name,
							DefaultGateway: isDefaultGateway,
						})
					}
				}
			}
		}
	}
	c.JSON(resp.Data(list))
}

// @Summary add a forward and start
// @Description ""
// @Tags network
// @Accept  json
// @Produce  json
// @Param json body dto.AddForward true "请求json"
// @Success 200 {object} resp.DataResp{}
// @Router /v1/forward/ [post]
func AddForward(c *gin.Context) {
	var req dto.AddForward
	var err error
	err = c.BindJSON(&req)
	if err != nil {
		logger.Error(err)
		c.JSON(resp.UnexpectedError(err))
		return
	}
	err = forward.Add(c, req)
	if err != nil {
		c.JSON(resp.UnexpectedError(err))
		return
	}
	c.JSON(resp.Data("ok"))
}
