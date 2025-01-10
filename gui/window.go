package gui

import (
	"fmt"
	"github.com/lxn/walk"
	"github.com/lxn/win"
	"golang.org/x/sys/windows/registry"
	"os"
	"os/exec"
	"os_manage/config"
	"os_manage/controller"
	"os_manage/log"
	"path/filepath"
	"regexp"
)

// ----------------------------- gui -----------------------------------

type MyWindow struct {
	*walk.MainWindow
	hWnd           win.HWND
	minimizeBox    *walk.CheckBox
	maximizeBox    *walk.CheckBox
	closeBox       *walk.CheckBox
	autoBootBox    *walk.CheckBox
	biliLineEdit   *walk.LineEdit
	biliCookieEdit *walk.LineEdit
	progressBar    *walk.ProgressBar
	logArea        *walk.TextEdit
}

var myWindow *MyWindow

func (mw *MyWindow) SetMinimizeBox() {
	if mw.minimizeBox.Checked() {
		mw.addStyle(win.WS_MINIMIZEBOX)
		return
	}
	mw.removeStyle(^win.WS_MINIMIZEBOX)
}

func (mw *MyWindow) SetMaximizeBox() {
	if mw.maximizeBox.Checked() {
		mw.addStyle(win.WS_MAXIMIZEBOX)
		return
	}
	mw.removeStyle(^win.WS_MAXIMIZEBOX)
}

func (mw *MyWindow) addStyle(style int32) {
	currStyle := win.GetWindowLong(mw.hWnd, win.GWL_STYLE)
	win.SetWindowLong(mw.hWnd, win.GWL_STYLE, currStyle|style)
}

func (mw *MyWindow) removeStyle(style int32) {
	currStyle := win.GetWindowLong(mw.hWnd, win.GWL_STYLE)
	win.SetWindowLong(mw.hWnd, win.GWL_STYLE, currStyle&style)
}

func (mw *MyWindow) SetCloseBox() {
	if mw.closeBox.Checked() {
		win.GetSystemMenu(mw.hWnd, true)
		return
	}
	hMenu := win.GetSystemMenu(mw.hWnd, false)
	win.RemoveMenu(hMenu, win.SC_CLOSE, win.MF_BYCOMMAND)
}

// 检测注册表key 决定初始是否勾选
func isAutoBoot() bool {
	key, err := registry.OpenKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Run`, registry.ALL_ACCESS)
	if err != nil {
		return false
	}
	defer key.Close()

	value, _, err := key.GetStringValue(config.RegeditKey)
	if err != nil {
		return false
	}

	return value == filepath.Join(config.HomeDir, "auto_boot.bat")
}

func autoBoot(state bool) {
	scriptPath := filepath.Join(config.HomeDir, "auto_boot.bat")
	// the regedit path
	keyPath := `Software\Microsoft\Windows\CurrentVersion\Run`
	key, _, err := registry.CreateKey(registry.CURRENT_USER, keyPath, registry.ALL_ACCESS)
	if err != nil {
		log.Error("Error creating registry key:", err)
		return
	}
	defer key.Close()

	if state {
		// 创建启动脚本
		if _, err := os.Stat(scriptPath); err != nil {
			_ = os.MkdirAll(config.HomeDir, 0644)
			nowPath, err := os.Executable()
			if err != nil {
				log.Error("set auto boot failed cause getting executable path:", err)
				return
			}
			nowDir, nowName := filepath.Split(nowPath)
			err = os.WriteFile(scriptPath, []byte(fmt.Sprintf(`cd /d %s
.\%s
`, nowDir, nowName)), 0644)
			if err != nil {
				log.Error("set auto boot failed cause writing to file:", err)
				return
			}
		}

		// 设置开机自启
		err = key.SetStringValue(config.RegeditKey, scriptPath)
		if err != nil {
			log.Error("setting auto_boot.bat failed:", err)
			return
		}

		log.Debug("set auto boot success", scriptPath)
	} else {
		err = key.DeleteValue(config.RegeditKey)
		if err != nil {
			log.Errorf("set program autoBoot state[%v] failed: %v\n", state, err)
		}
		_ = os.RemoveAll(scriptPath)

		log.Debug("set not auto boot success", scriptPath)
	}
}

func (mw *MyWindow) SetAutoBootBox() {
	if mw.autoBootBox.Checked() {
		autoBoot(true)
		return
	}

	autoBoot(false)
}

// AddIcon 新建图标
func (mw *MyWindow) AddIcon(iconPath string) {
	icon, err := walk.Resources.Image(iconPath)
	if err != nil {
		log.Errorf("walk get image icon failed: %v\n", err)
		return
	}
	_ = mw.SetIcon(icon)
}

func GetPanelUI() *MyWindow {
	if myWindow == nil {
		log.Debug("try to new main panel")
		myWindow = newMainPanel()
	}

	return myWindow
}

func (mw *MyWindow) DownloadBiliVideo() {
	bvReg := regexp.MustCompile(`BV[a-zA-Z0-9]+`)
	bvId := bvReg.FindString(mw.biliLineEdit.Text())
	if bvId == "" {
		log.Error("there is not found bv in:", mw.biliLineEdit.Text())
		return
	}
	log.Debug("try to download bv", bvId)

	err := controller.DownloadByBvID(bvId, config.GlobalConfig.Bili.SavePath, config.GlobalConfig.Bili.SaveCover)
	if err != nil {
		log.Errorf("download bv[%s] error: %v", bvId, err)
		return
	}
	log.Info("download video to", config.GlobalConfig.Bili.SavePath, "success")
	mw.biliLineEdit.SetText("")
}

func (mw *MyWindow) OpenBiliSavePath() {
	abs, err := filepath.Abs(config.GlobalConfig.Bili.SavePath)
	if err != nil {
		log.Error("get abs path error:", err)
		return
	}
	err = exec.Command("explorer", abs).Start()
	if err != nil {
		log.Error("open bili save path error:", err)
	}

	log.Debug("open bili save path:", abs)
}

func (mw *MyWindow) SaveBiliCookie() {
	cookie := mw.biliCookieEdit.Text()
	config.GlobalConfig.Bili.Cookie = cookie
	err := config.SaveConfigToTomlFile()
	if err != nil {
		log.Error("save cookie error:", err)
	} else {
		mw.biliCookieEdit.SetText("")
		log.Debug("save new cookie success:", config.GlobalConfig.Bili.Cookie)
	}
}
