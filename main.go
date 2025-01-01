package main

import (
	"errors"
	"net/http"
	"os_manage/config"
	"os_manage/controller"
	"os_manage/gui"
	"os_manage/log"
	"os_manage/router"
)

func main() {
	logger := log.NewLogger(
		log.WithLogLevel(config.GlobalConfig.Log.LogLevel),
		log.WithStorePath(config.GlobalConfig.Log.LogPath),
		log.WithLogExtend(make(chan string, 10), func(eha any, msg string) {
			ch, ok := eha.(chan string)
			if ok {
				ch <- msg
			}
		}),
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

	gui.NewTray()

	beforeQuit(controller.QuitVNCServer)
}

func beforeQuit(fns ...func()) {
	for _, fn := range fns {
		fn()
	}
}
