package router

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"os_manage/controller"
)

func NewRouter() *gin.Engine {
	router := gin.Default()
	RouteGroups(router)

	return router
}

func RouteGroups(router *gin.Engine) {
	router.GET("/help", RouteHelp)

	router.GET("/shutdown", controller.Shutdown)
	router.GET("/shutdown/cancel", controller.ShutdownCancel)

	router.GET("/monitor/turn_off", controller.MonitorTurnOff)
	router.GET("/monitor/turn_on", controller.MonitorTurnOn)

	// 非常吃系统资源 且只能预览
	//router.GET("/capture", controller.Capture)
	//router.GET("/capture/ws", controller.CaptureWSConnection)

	router.GET("/vnc", controller.StartVNC)
	router.GET("/vnc/quit", controller.ShutdownVNC)

	controller.RegisterProxy(router)

	router.GET("/bili/flush", controller.FlushBiliCookie)
	router.GET("/bili/:bv/*ex", controller.BiliDownload)

	// todo 浏览本地文件
}

func RouteHelp(c *gin.Context) {
	c.String(http.StatusOK, `
/shutdown: arg[after] second, default[0]
/shutdown/cancel: none arg.
/monitor/turn_off: turn off the power of monitor.
/monitor/turn_on: turn on the power of monitor.
`)
}
