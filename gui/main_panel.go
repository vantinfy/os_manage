package gui

import (
	"fmt"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"github.com/lxn/win"
	"os_manage/config"
	"os_manage/controller"
	"os_manage/log"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	biliGroupBox GroupBox
)

func getBiliGroupBox(mw *MyWindow) GroupBox {
	biliSyncOnce.Do(func() {
		biliGroupBox = GroupBox{
			Layout: HBox{}, Title: "b站视频下载",
			Children: []Widget{
				LineEdit{
					AssignTo:    &mw.biliLineEdit,
					ToolTipText: "支持bv或包含bv的完整链接",
				},
				PushButton{
					MinSize: Size{Width: 60, Height: 37}, Text: "下载",
					OnClicked: func() {
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
					},
				},
			},
		}
	})

	return biliGroupBox
}

func newMainPanel() *MyWindow {
	mw := &MyWindow{}

	if err := (MainWindow{
		AssignTo: &mw.MainWindow,
		Title:    "os manage",
		Size:     Size{Width: 560, Height: 480},
		Layout:   VBox{},
		Children: []Widget{
			CheckBox{
				Checked:             isAutoBoot(),
				AssignTo:            &mw.autoBootBox,
				Text:                "开机启动",
				OnCheckStateChanged: mw.SetAutoBootBox,
			},
			getBiliGroupBox(mw),
			PushButton{
				Text: "碧蓝航线科技",
				OnClicked: func() {
					azureLanePanel = &AzureLanePanel{}
					err := MainWindow{
						AssignTo: &azureLanePanel.MainWindow, Size: Size{Width: 560, Height: 370},
						Layout: VBox{}, Title: "碧蓝航线科技",
						Children: []Widget{
							getAzureLaneBox(azureLanePanel),
							Composite{
								Layout: HBox{},
								Children: []Widget{
									PushButton{Text: "重置"}, PushButton{Text: "查询"},
								},
							},
							Composite{}, // 结果
						},
					}.Create()
					if err != nil {
						log.Error("create azure lane error:", err)
					}
					azureLanePanel.Run()
				},
			},
			TextEdit{ // 日志打印区
				ReadOnly: true, VScroll: true, AssignTo: &mw.logArea,
			},
		},
	}.Create()); err != nil {
		log.Fatal(err)
	}

	// 叉掉窗口转右下角后台运行
	mw.Closing().Attach(func(canceled *bool, reason walk.CloseReason) {
		//if reason == walk.CloseReasonUser {
		if reason == walk.CloseReasonUnknown {
			*canceled = true // 阻止关闭
			mw.Hide()        // 隐藏窗口
			log.Debug("main panel hided")
		}
	})

	// 设置窗口句柄
	mw.hWnd = mw.Handle()
	config.MainPanelHWND = mw.hWnd

	mw.AddIcon(config.AppIconPath)

	go func() { // 监听log并输出到textEdit
		errChan, ok := log.GetLogger().Extend.(chan string)
		if !ok {
			return
		}

		for {
			msg := <-errChan
			text := strings.Replace(msg, "\n", "\r\n", -1)
			mw.logArea.SetText(fmt.Sprintf("%s%s", mw.logArea.Text(), text))
		}
	}()

	return mw
}

var hWnd win.HWND
var dialog *walk.Dialog
var progress *walk.LineEdit
var confirm *walk.PushButton

func doProgress(mw *MyWindow) {
	// 进度条弹窗
	err := Dialog{
		AssignTo: &dialog,
		Title:    "Progress dialog",
		MinSize:  Size{Width: 300, Height: 200},
		Layout:   VBox{},
		Name:     "ProgressBar",
		Children: []Widget{
			Label{
				Text:   "当前进度:",
				Row:    1,
				Column: 1,
			},
			LineEdit{
				AssignTo: &progress,
				ReadOnly: true,
				Row:      1,
				Column:   2,
			},
			ProgressBar{AssignTo: &mw.progressBar},
			PushButton{
				AssignTo: &confirm,
				Text:     `执行完毕，退出`,
				Enabled:  false, //默认不可用
				OnClicked: func() {
					mw.biliLineEdit.SetText("")
					// 因为弹窗关闭键被取消,按键关闭
					dialog.Close(0)
				},
			},
		},
	}.Create(mw)
	if err != nil {
		fmt.Println(err)
		return
	}
	// 禁止关闭
	hWnd = dialog.Handle()
	hMenu := win.GetSystemMenu(hWnd, false)
	win.RemoveMenu(hMenu, win.SC_CLOSE, win.MF_BYCOMMAND)
	// 开始跑进度条
	dialog.Starting().Attach(func() {
		go progressWorker(mw)
	})
	dialog.Run()
}

func progressWorker(mw *MyWindow) {
	length := 10
	dialog.Synchronize(func() {
		// progressBar.SetRange(0, int(length))
		// 进度条value的起止   value到5进度条开始做走，到10停止
		mw.progressBar.SetRange(0, 20)
	})
	workWithCallback(length+10, func(n int64) {
		fmt.Println("progress", n)
		dialog.Synchronize(func() {
			mw.progressBar.WidgetBase.SetToolTipText(strconv.FormatInt(n, 10))
			mw.progressBar.SetValue(int(n))
			progress.SetText(strconv.FormatInt(n, 10))
		})
	})
	// 跑完后按键可用
	confirm.SetEnabled(true)
}

func workWithCallback(length int, callback func(int64)) {
	// value从1到20
	for i := 1; i <= length; i++ {
		time.Sleep(time.Millisecond * 100)
		callback(int64(i))
	}
}
