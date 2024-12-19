package controller

import (
	"github.com/gin-gonic/gin"
	"html/template"
	"net/http"
	"os"
	"os/exec"
	"os_manage/config"
	"os_manage/log"
	"os_manage/tools"
	"path/filepath"
	"time"
)

var vncPath = filepath.Join(config.ProcessWorkDir, "vnc/UltraVNC_1436/x64/winvnc.exe")
var vncServerCmd = exec.Command(vncPath)

const (
	downloadLinkUltraVNC = "https://uvnc.eu/download/1436/UltraVNC_1436.zip"
	savePathUltraVNC     = "vnc/UltraVNC_1436.zip"
	destPathUltraVNC     = "vnc/UltraVNC_1436"
	downloadLinkNoVNC    = "https://github.com/novnc/noVNC/archive/refs/tags/v1.5.0.zip"
	savePathNoVNC        = "vnc/noVNC-1.5.0.zip"
	destPathNoVNC        = "vnc/noVNC-1.5.0"
)

func prepareVNCEnv() error {
	_, err := os.Stat("vnc")
	if err != nil {
		_ = os.MkdirAll("vnc", 0644)
	}

	// 本项目的VNC server端服务基于UltraVNC-1.4.3.6实现 遵守UltraVNC的相关协议
	// 可执行文件下载来源：https://uvnc.com/downloads/ultravnc.html
	// 另外注意第一次使用UltraVNC时可能会要求设置连接密码
	// UltraVNC仓库详情可参见:
	// https://github.com/ultravnc/UltraVNC
	_, err = os.Stat(vncPath)
	if err != nil {
		log.Info("try to download UltraVNC")
		err = tools.DownloadFile(downloadLinkUltraVNC, savePathUltraVNC)
		if err != nil {
			return err
		}
		err = tools.Unzip(savePathUltraVNC, destPathUltraVNC)
		if err != nil {
			return err
		}
	}
	// 本项目的noVNC 网页端服务基于noVNC-1.5.0实现 遵守noVNC相关协议
	// 相关文件来源：https://github.com/novnc/noVNC/archive/refs/tags/v1.5.0.zip
	// 在下载过程中可能因为网络问题而下载失败，因此下面的tryDownloadNoVNC可能会使用代理下载相关文件
	// noVNC仓库详情可参见:
	// https://github.com/novnc/noVNC
	_, err = os.Stat(filepath.Join(config.ProcessWorkDir, "vnc/web/vnc.html"))
	if err != nil {
		log.Info("try to download noVNC")
		err = tryDownloadNoVNC()
		if err != nil {
			return err
		}
		err = tools.Unzip(savePathNoVNC, "vnc")
		if err != nil {
			return err
		}
		_ = os.Rename(destPathNoVNC, "vnc/web")
	}

	_ = os.RemoveAll(savePathUltraVNC)
	_ = os.RemoveAll(savePathNoVNC)

	return nil
}

func tryDownloadNoVNC() error {
	// noVNC的资源下载优先直连github 失败则尝试通过代理下载
	proxySites := []string{"", "https://gh.llkk.cc/", "https://ghproxy.cn/", "https://github.moeyy.xyz"}

	for _, proxySite := range proxySites {
		err := tools.DownloadFile(proxySite+downloadLinkNoVNC, savePathNoVNC)
		if err != nil {
			continue
		}
		if fileBytes, err := os.ReadFile(savePathNoVNC); err != nil || len(fileBytes) == 0 {
			continue
		}
		return nil
	}

	return nil
}

func StartVNC(c *gin.Context) {
	err := prepareVNCEnv()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	// 如果vnc server未启动 或者挂了 则尝试协程启动
	if vncServerCmd.Process == nil ||
		(vncServerCmd.Process != nil && vncServerCmd.ProcessState != nil) {
		vncServerCmd = exec.Command(vncPath)
		go func() {
			log.Info("try to start vnc server")
			// 启动vnc服务
			if err := vncServerCmd.Run(); err != nil {
				log.Error("start vnc server error: ", err)
				return
			}
		}()
	}

	//c.HTML(http.StatusOK, "vnc/web/vnc.html", gin.H{})
	t, err := template.ParseFiles("vnc/web/vnc.html")
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
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
			log.Error("start vnc server error: ", err)
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
	}

	c.String(http.StatusOK, "ok")
}

func QuitVNCServer() {
	if vncServerCmd.Process != nil && vncServerCmd.ProcessState == nil {
		err := vncServerCmd.Process.Kill()
		if err != nil {
			log.Error("quit UltraVNC err:", err)
		}
	}
}
