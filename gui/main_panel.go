package gui

import (
	"fmt"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"github.com/lxn/win"
	"os_manage/config"
	"os_manage/log"
	"strconv"
	"time"
)

func newMainPanel() *MyWindow {
	mw := &MyWindow{}

	if err := (MainWindow{
		AssignTo: &mw.MainWindow,
		Title:    "notify icon",
		Size:     Size{550, 380},
		Layout:   VBox{},
		Children: []Widget{
			//CheckBox{
			//	AssignTo:            &mw.minimizeBox,
			//	Text:                "显示最小化按钮",
			//	Checked:             true,
			//	OnCheckStateChanged: mw.SetMinimizeBox,
			//},
			//CheckBox{
			//	AssignTo:            &mw.maximizeBox,
			//	Text:                "显示最大化按钮",
			//	Checked:             true,
			//	OnCheckStateChanged: mw.SetMaximizeBox,
			//},
			//CheckBox{
			//	AssignTo:            &mw.closeBox,
			//	Text:                "显示关闭按钮",
			//	Checked:             true,
			//	OnCheckStateChanged: mw.SetCloseBox,
			//},
			CheckBox{
				AssignTo:            &mw.autoBootBox,
				Text:                "开机启动",
				OnCheckStateChanged: mw.SetAutoBootBox,
			},
			LineEdit{
				ToolTipText: "tips",
			},
			PushButton{
				MinSize: Size{Width: 60, Height: 37},

				Text: "Show progress",
				OnClicked: func() {
					doProgress(mw)
				},
			},
			TextEdit{
				ReadOnly: true,
			}, // todo log print
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
			log.Info("main panel hided")
		}
	})

	// 设置窗口句柄
	mw.hWnd = mw.Handle()
	config.MainPanelHWND = mw.hWnd

	mw.AddIcon(config.AppIconPath)

	go func() {
		// mw.logArea
	}()

	return mw
}

var hWnd win.HWND
var dialog *walk.Dialog
var jindu *walk.LineEdit
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
				AssignTo: &jindu,
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
			jindu.SetText(strconv.FormatInt(n, 10))
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
