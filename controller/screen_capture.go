package controller

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/kbinani/screenshot"
	"image/png"
	"net/http"
	"os_manage/log"
	"time"
)

func Capture(ctx *gin.Context) {
	// 设置响应头，指定内容类型为 HTML
	ctx.Writer.Header().Set("Content-Type", "text/html; charset=utf-8")

	// 写入 HTML 内容
	htmlContent := `
    <!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title>Screen Capture</title>
	</head>
	<body>
	<h1>实时屏幕捕获</h1>
	<img id="screen" src="" style="max-width: 100%;"/>
	
	<script>
		const ws = new WebSocket('ws://192.168.1.102:7799/capture/ws');
	
		ws.binaryType = 'arraybuffer';
	
		ws.onmessage = function(event) {
			const blob = new Blob([event.data], { type: 'image/png' });
			const url = URL.createObjectURL(blob);
			document.getElementById('screen').src = url;
		};
	</script>
	</body>
	</html>
    `
	fmt.Fprint(ctx.Writer, htmlContent)
}

var upgrader = websocket.Upgrader{
	//ReadBufferSize:  1024,
	//WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // 为了简化，不检查来源，注意在生产环境中需要谨慎
	},
}

func CaptureWSConnection(ctx *gin.Context) {
	ws, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		fmt.Println("could not upgrade connection", err)
		return
	}
	captureScreen(ws)
}

func captureScreen(ws *websocket.Conn) {
	defer ws.Close()
	for {
		img, err := screenshot.CaptureDisplay(0)
		if err != nil {
			// shot display failed: GetDIBits failed
			log.Error("shot display failed:", err)

			time.Sleep(time.Millisecond * 16)
			continue
		}

		buf := new(bytes.Buffer)
		err = png.Encode(buf, img)
		if err != nil {
			log.Error("encode img failed", err)
			return
		}

		// 发送图像数据
		err = ws.WriteMessage(websocket.BinaryMessage, buf.Bytes())
		if err != nil {
			log.Error("send img failed", err)
			return
		}

		// 每隔16毫秒捕获一次屏幕
		time.Sleep(16 * time.Millisecond)
	}
}
