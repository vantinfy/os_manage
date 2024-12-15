package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os/exec"
	"os_manage/log"
)

// ----------------------------- about shutdown windows -----------------------------------

func Shutdown(c *gin.Context) {
	after := c.Query("after")

	cmd := exec.Command("shutdown", "/s", "/t", fmt.Sprintf("%v", after))
	if output, err := cmd.Output(); err != nil {
		log.Error("shutdown failed", string(output), err)
		c.String(http.StatusExpectationFailed, "%v", err.Error())
		return
	}
	c.String(http.StatusOK, "ok")
}

func ShutdownCancel(c *gin.Context) {
	cmd := exec.Command("shutdown", "/a")
	if output, err := cmd.Output(); err != nil {
		log.Error("cancel shutdown failed", string(output), err)
		c.String(http.StatusExpectationFailed, "%v", err.Error())
		return
	}
	c.String(http.StatusOK, "ok")
}
