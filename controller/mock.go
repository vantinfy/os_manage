package controller

import "github.com/go-vgo/robotgo"

func mock() {
	// 处理交互，例如移动鼠标或模拟按键

	//"github.com/go-vgo/robotgo" 需要安装gcc https://blog.csdn.net/qq_35977117/article/details/139070603
	//robotgo.Move(100, 100, 0) // 移动鼠标到 (100, 100)

	for {
		robotgo.ScrollDir(10, "down")
		robotgo.MilliSleep(5000)
	}
}
