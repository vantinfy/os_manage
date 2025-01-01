package gui

import (
	"github.com/getlantern/systray"
	"os"
	"os_manage/config"
	"os_manage/log"
)

// ----------------------------- hide program -----------------------------------

var iconData []byte

func getIconData(iconPath string) []byte {
	iconData, _ = os.ReadFile(iconPath)
	// todo 读取不到文件
	return iconData
}

func NewTray() {
	systray.Run(onReady, onExit)
}

func onReady() {
	systray.SetIcon(getIconData(config.AppIconPath)) // 设置托盘图标
	systray.SetTitle("my app")
	systray.SetTooltip("OSManage")

	mOpen := systray.AddMenuItem("Open", "Open My App")
	//systray.AddSeparator() // 分隔线
	mExit := systray.AddMenuItem("Exit", "Exit My App")

	go func() {
		GetPanelUI().Run()
	}()

	go func() {
		for {
			select {
			case <-mOpen.ClickedCh: // 用户点击了Open菜单项
				openMainPanel()
			case <-mExit.ClickedCh: // 用户点击了Exit菜单项
				systray.Quit()
				return
			}
		}
	}()
}

func openMainPanel() {
	// 打开应用程序的主窗口
	GetPanelUI().Show()
}

func onExit() {
	// 托盘程序退出时的清理逻辑
	log.Info("托盘退出")
}
