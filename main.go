package main

import (
	"errors"
	"net/http"
	"os_manage/config"
	"os_manage/controller"
	"os_manage/gui"
	"os_manage/log"
	"os_manage/router"
	"syscall"
)

func main() {
	logger := log.NewLogger(
		log.WithLogLevel(config.GlobalConfig.Log.LogLevel),
		log.WithStorePath(config.GlobalConfig.Log.LogPath),
	)

	server := &http.Server{
		Addr:    config.GlobalConfig.App.ServeAddr,
		Handler: router.NewRouter(),
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal("listen:", err)
		}
	}()

	// create a channel to receive signal
	quit := config.GlobalQuit

	go func() {
		gui.NewTray()
		// 托盘退出时 主程序同步退出
		quit <- syscall.SIGQUIT
	}()

	<-quit
	beforeQuit(controller.QuitVNCServer)
}

func beforeQuit(fns ...func()) {
	for _, fn := range fns {
		fn()
	}
}
