package tools

import (
	"io"
	"net/http"
	"os"
	"os/exec"
)

func DownloadFile(downloadURL, savePath string) error {
	// 如果系统支持curl
	curlPath, err := exec.LookPath("curl")
	if err == nil {
		return exec.Command(curlPath, "-o", savePath, downloadURL).Run()
	}

	resp, err := http.DefaultClient.Get(downloadURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	fileBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	err = os.WriteFile(savePath, fileBytes, 0644)
	if err != nil {
		return err
	}

	return nil
}
