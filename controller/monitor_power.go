package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/go-vgo/robotgo"
	"os_manage/config"
	"os_manage/log"

	"net/http"
	"syscall"
)

// ----------------------------- about monitor -----------------------------------

const (
	SCMonitorPower = 0xF170
	MonitorOff     = 2
	MonitorOn      = -1
)

func setMonitorPower(state int) {
	user32 := syscall.NewLazyDLL("user32.dll")
	sendMessage := user32.NewProc("SendMessageW") // 异步可以使用PostMessageW
	//desktopHwnd, _, _ := user32.NewProc("GetDesktopWindow").Call()

	_, _, err := sendMessage.Call(
		//desktopHwnd, // 桌面窗口句柄 // HWND_BROADCAST-0xFFFF
		uintptr(config.MainPanelHWND), // 直接使用当前程序的句柄
		0x0112,                        // WM_SYSCOMMAND
		uintptr(SCMonitorPower),
		uintptr(state),
	)
	if err != nil {
		log.Error("set monitor power failed:", err)
	}
}

func MonitorTurnOff(c *gin.Context) {
	// cmd
	// %systemroot%\system32\scrnsave.scr /s
	setMonitorPower(MonitorOff)
	c.String(http.StatusOK, "turn off")
}

func MonitorTurnOn(c *gin.Context) {
	// setMonitorPower(MonitorOn)

	robotgo.ScrollDir(10, "down") // 模拟鼠标滚动唤醒屏幕电源
	c.String(http.StatusOK, "turn on")
}
