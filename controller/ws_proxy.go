package controller

import (
	"net"
	"net/http"
	"os"
	"os_manage/config"
	"os_manage/log"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func RegisterProxy(r *gin.Engine) {
	serveFile := func(path string) func(ctx *gin.Context) {
		return func(ctx *gin.Context) {
			filePath := filepath.Join(config.ProcessWorkDir, "vnc/web", path)
			fileBytes, err := os.ReadFile(filePath)
			if err != nil {
				ctx.Writer.WriteHeader(http.StatusInternalServerError)
				ctx.Writer.WriteString(err.Error())
				return
			}
			ctx.Writer.WriteHeader(http.StatusOK)
			ctx.Writer.Write(fileBytes)
		}
	}
	r.GET("/package.json", serveFile("package.json"))
	r.GET("/defaults.json", serveFile("defaults.json"))
	r.GET("/mandatory.json", serveFile("mandatory.json"))

	r.Static("./app", "vnc/web/app")
	r.Static("./core", "vnc/web/core")
	r.Static("./vendor", "vnc/web/vendor")

	r.GET("/websockify", func(c *gin.Context) {
		serveWs(c.Writer, c.Request)
	})
}

func serveWs(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Errorf("%s: failed to upgrade to WS: %s", time.Now().Format(time.Stamp), err)
		return
	}

	vnc, err := net.Dial("tcp", ":5900") // todo 端口 配置化
	if err != nil {
		log.Errorf("%s: failed to bind to the VNC Server: %s", time.Now().Format(time.Stamp), err)
	}

	go forwardTcp(ws, vnc)
	go forwardWeb(ws, vnc)
}

func forwardTcp(wsConn *websocket.Conn, conn net.Conn) {
	var tcpBuffer [1024]byte
	defer func() {
		if conn != nil {
			_ = conn.Close()
		}
		if wsConn != nil {
			_ = wsConn.Close()
		}
	}()

	for {
		if (conn == nil) || (wsConn == nil) {
			return
		}
		n, err := conn.Read(tcpBuffer[0:])
		if err != nil {
			log.Errorf("%s: reading from TCP failed: %s", time.Now().Format(time.Stamp), err)
			return
		} else {
			if err := wsConn.WriteMessage(websocket.BinaryMessage, tcpBuffer[0:n]); err != nil {
				log.Errorf("%s: writing to WS failed: %s", time.Now().Format(time.Stamp), err)
			}
		}
	}
}

func forwardWeb(wsConn *websocket.Conn, conn net.Conn) {
	defer func() {
		if err := recover(); err != nil {
			log.Errorf("%s: reading from WS failed: %s", time.Now().Format(time.Stamp), err)
		}
		if conn != nil {
			_ = conn.Close()
		}
		if wsConn != nil {
			_ = wsConn.Close()
		}
	}()

	for {
		if (conn == nil) || (wsConn == nil) {
			return
		}

		_, buffer, err := wsConn.ReadMessage()
		if err == nil {
			if _, err := conn.Write(buffer); err != nil {
				log.Errorf("%s: writing to TCP failed: %s", time.Now().Format(time.Stamp), err)
			}
		}
	}
}
