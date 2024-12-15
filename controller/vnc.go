package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"os/exec"
	"os_manage/config"
	"os_manage/log"
	"path/filepath"
)

var vncPath = filepath.Join(config.ProcessWorkDir, "vnc/UltraVNC_1436/x64/winvnc.exe")
var wsResourcePath = filepath.Join(config.ProcessWorkDir, "vnc/web")
var vncServerCmd = exec.Command(vncPath)
var wsProxyCmd = exec.Command("websockify", "5900", ":5901", "--web", wsResourcePath)

func StartVNC(c *gin.Context) {
	redirect := "http://localhost:5900/vnc.html"
	// 服务已经启动的情况下
	if wsProxyCmd.Process != nil && vncServerCmd.Process != nil {
		c.Redirect(http.StatusFound, redirect)
		return
	}

	if vncServerCmd.Process == nil ||
		(vncServerCmd.Process != nil && vncServerCmd.ProcessState != nil) {
		vncServerCmd = exec.Command(vncPath)
		go func() {
			log.Info("try to start vnc server")
			// 启动vnc服务
			if err := vncServerCmd.Run(); err != nil {
				c.Writer.WriteHeader(http.StatusInternalServerError)
				c.Writer.Write([]byte(err.Error()))
				log.Error("start vnc server error: ", err)
			}
		}()
	}
	if wsProxyCmd.Process == nil ||
		(vncServerCmd.Process != nil && vncServerCmd.ProcessState != nil) {
		wsProxyCmd = exec.Command("websockify", "5900", ":5901", "--web", wsResourcePath)
		go func() {
			log.Info("try to start ws proxy")
			// 启动ws转发服务
			if err := wsProxyCmd.Run(); err != nil {
				c.Writer.WriteHeader(http.StatusInternalServerError)
				c.Writer.Write([]byte(err.Error()))
				log.Error("start vnc server error: ", err)
			}
		}()
	}

	c.Redirect(http.StatusFound, redirect)
}

func QuitVNC(c *gin.Context) {
	if wsProxyCmd.Process != nil {
		err := wsProxyCmd.Process.Kill()
		if err != nil {
			c.Writer.WriteHeader(http.StatusInternalServerError)
			c.Writer.Write([]byte(err.Error()))
			log.Error("start vnc server error: ", err)
			return
		}
	}
	if vncServerCmd.Process != nil {
		err := vncServerCmd.Process.Kill()
		if err != nil {
			c.Writer.WriteHeader(http.StatusInternalServerError)
			c.Writer.Write([]byte(err.Error()))
			log.Error("start vnc server error: ", err)
			return
		}
	}

	c.Writer.WriteHeader(http.StatusOK)
	c.Writer.Write([]byte("ok"))
}

func QuitWsAndVNC() {
	if vncServerCmd.Process != nil && vncServerCmd.ProcessState == nil {
		err := vncServerCmd.Process.Kill()
		if err != nil {
			log.Error("cmd.Process.Kill() err:", err)
		}
	}
	if wsProxyCmd.Process != nil && wsProxyCmd.ProcessState == nil {
		err := wsProxyCmd.Process.Kill()
		if err != nil {
			log.Error("wsCmd.Process.Kill() err:", err)
		}
	}
}
