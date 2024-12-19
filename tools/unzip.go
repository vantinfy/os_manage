package tools

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
)

func Unzip(src, dest string) error {
	// 如果系统支持unzip命令
	unzipPath, err := exec.LookPath("unzip")
	if err == nil {
		return exec.Command(unzipPath, "-n", src, "-d", dest).Run()
	}

	// 打开 ZIP 文件
	r, err := zip.OpenReader(src)
	if err != nil {
		return fmt.Errorf("无法打开文件: %v", err)
	}
	defer r.Close()

	// 遍历 ZIP 文件中的每个文件
	for _, file := range r.File {
		path := filepath.Join(dest, file.Name)

		// 检查路径安全性
		if !filepath.HasPrefix(path, filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("非法文件路径: %s", path)
		}

		// 如果是目录，则创建目录
		if file.FileInfo().IsDir() {
			err = os.MkdirAll(path, os.ModePerm)
			if err != nil {
				return fmt.Errorf("创建目录失败: %v", err)
			}
			continue
		}

		// 创建文件所在的目录
		err = os.MkdirAll(filepath.Dir(path), os.ModePerm)
		if err != nil {
			return fmt.Errorf("创建父目录失败: %v", err)
		}

		// 解压文件
		err = extractFile(file, path)
		if err != nil {
			return fmt.Errorf("解压文件失败: %v", err)
		}
	}

	return nil
}

// 解压单个文件
func extractFile(file *zip.File, path string) error {
	srcFile, err := file.Open()
	if err != nil {
		return err
	}
	defer srcFile.Close()

	destFile, err := os.Create(path)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, srcFile)
	return err
}
