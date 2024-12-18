package controller

import (
	"github.com/gin-gonic/gin"
	"html/template"
	"net/http"
	"os/exec"
	"os_manage/config"
	"os_manage/log"
	"path/filepath"
	"time"
)

var vncPath = filepath.Join(config.ProcessWorkDir, "vnc/UltraVNC_1436/x64/winvnc.exe")
var vncServerCmd = exec.Command(vncPath)

func StartVNC(c *gin.Context) {
	// 如果vnc server未启动 或者挂了 则尝试协程启动
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

	//c.HTML(http.StatusOK, "vnc/web/vnc.html", gin.H{})
	t, err := template.ParseFiles("vnc/web/vnc.html")
	if err != nil {
		c.Writer.WriteString(err.Error())
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	for {
		// 等待vnc server启动 约6ms
		time.Sleep(time.Millisecond * 5)
		if vncServerCmd.Process != nil {
			t.Execute(c.Writer, "vnc")
			return
		}
	}
}

func ShutdownVNC(c *gin.Context) {
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

func QuitVNCServer() {
	if vncServerCmd.Process != nil && vncServerCmd.ProcessState == nil {
		err := vncServerCmd.Process.Kill()
		if err != nil {
			log.Error("cmd.Process.Kill() err:", err)
		}
	}
}
