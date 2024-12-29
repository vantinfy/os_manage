package controller

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"os"
	"os/exec"
	"os_manage/config"
	"os_manage/log"
	"os_manage/tools"
	"path/filepath"
	"regexp"
)

func FlushBiliCookie(c *gin.Context) {
	config.LoadConfig()
	c.String(http.StatusOK, "flush cookie done. new cookie len: %d", len(config.GlobalConfig.Bili.Cookie))
}

func BiliDownload(c *gin.Context) {
	bvId := c.Param("bv")
	err := DownloadByBvID(bvId, config.GlobalConfig.Bili.SavePath, config.GlobalConfig.Bili.SaveCover)
	if err != nil {
		c.String(http.StatusServiceUnavailable, err.Error())
	}
	c.String(http.StatusOK, "download success")
}

// DownloadByBvID 通过bv号下载b站视频
//
// 原作者b站@゚゚未闻花名 https://space.bilibili.com/630468506
//
// 源代码python https://github.com/T-Tedebug/python-bilibili-downloads
func DownloadByBvID(bvId, savePath string, saveCover bool) error {
	title, cover, aid, cid, err := GetVideoInfo(bvId)
	if err != nil {
		return err
	}

	if saveCover { // 保存封面
		req, _ := http.NewRequest(http.MethodGet, cover, nil)
		setBiliHeader(req)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Error("save cover failed", bvId, title, err)
		}
		defer resp.Body.Close()

		coverBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Error("save cover failed", bvId, title, err)
		}
		_ = os.WriteFile(filepath.Join(savePath, title+".jpg"), coverBytes, 0644)
	}

	audioStreamUrl, videoStreamUrl, err := getVideoUrl(aid, cid)
	if err != nil {
		return err
	}

	err = downloadAVStream(audioStreamUrl, videoStreamUrl, title, savePath)
	if err != nil {
		return err
	}

	return mergeVideo(title, savePath)
}

// GetVideoInfo 根据bv解析返回视频标题、视频封面、aid、cid和错误
//
// aid和cid在获取音视频时需要用到
func GetVideoInfo(bvId string) (title string, cover string, aid int, cid int, err error) {
	infoUrl := fmt.Sprintf("https://api.bilibili.com/x/web-interface/view?bvid=%s", bvId)
	req, _ := http.NewRequest(http.MethodGet, infoUrl, nil)
	setBiliHeader(req)
	resp, err := http.Get(infoUrl)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	respMap := make(map[string]interface{})
	err = json.Unmarshal(respBytes, &respMap)
	if err != nil {
		return
	}
	if respMap["code"].(float64) != 0 {
		err = fmt.Errorf("%v", respMap["message"])
		return
	}

	title, _ = respMap["data"].(map[string]interface{})["title"].(string) // 视频标题
	illegalChars := regexp.MustCompile(`[<>:"/\\|?*\x00-\x1F]`)           // 确保文件名合法
	title = illegalChars.ReplaceAllString(title, "-")

	aidF, _ := respMap["data"].(map[string]interface{})["aid"].(float64)
	cidF, _ := respMap["data"].(map[string]interface{})["cid"].(float64)
	cover, _ = respMap["data"].(map[string]interface{})["pic"].(string) // 视频封面
	aid = int(aidF)
	cid = int(cidF)

	return
}

func getVideoUrl(aid, cid int) (audioStreamUrl, videoStreamUrl string, err error) {
	targetQuality := 127 // 直接默认最高画质 127对应8K
	srcUrl := fmt.Sprintf("https://api.bilibili.com/x/player/playurl?avid=%d&cid=%d&qn=%d&fnval=4048&fourk=1&fnver=0&session=", aid, cid, targetQuality)
	req, _ := http.NewRequest(http.MethodGet, srcUrl, nil)
	setBiliHeader(req)
	req.Header.Set("Referer", fmt.Sprintf("https://www.bilibili.com/video/av%d", aid))
	req.Header.Set("Origin", "https://www.bilibili.com")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	respMap := make(map[string]interface{})
	respBytes, _ := io.ReadAll(resp.Body)
	err = json.Unmarshal(respBytes, &respMap)
	if err != nil {
		return
	}
	if respMap["code"].(float64) != 0 {
		err = fmt.Errorf("%v", respMap["message"])
		return
	}

	data, _ := respMap["data"].(map[string]interface{})
	acceptList, _ := data["accept_quality"].([]interface{})
	if len(acceptList) > 0 {
		targetQuality = int(acceptList[0].(float64)) // 直接使用最高质量
	}
	dash, _ := data["dash"].(map[string]interface{})

	audios, _ := dash["audio"].([]interface{})
	if len(audios) > 0 {
		// 找到音频流
		audioStreamUrl = audios[0].(map[string]interface{})["baseUrl"].(string)
	}

	videos, _ := dash["video"].([]interface{})
	for _, videoI := range videos {
		video, _ := videoI.(map[string]interface{})
		if int(video["id"].(float64)) == targetQuality {
			// 找到视频流
			videoStreamUrl = video["baseUrl"].(string)
			break
		}
	}
	// 没有找到期望的画质 直接使用第一个
	if len(videoStreamUrl) == 0 {
		videoStreamUrl = videos[0].(map[string]interface{})["baseUrl"].(string)
	}

	return
}

func downloadAVStream(audioStreamUrl, videoStreamUrl, title, savePath string) (err error) {
	if _, err = os.Stat(savePath); os.IsNotExist(err) {
		_ = os.MkdirAll(savePath, 0644)
	}

	// 下载音频流
	req, _ := http.NewRequest(http.MethodGet, audioStreamUrl, nil)
	setBiliHeader(req)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	audioBytes, _ := io.ReadAll(resp.Body)
	err = os.WriteFile(filepath.Join(savePath, title+"_audio.m4s"), audioBytes, 0644)
	if err != nil {
		return err
	}

	// 下载视频流
	req, _ = http.NewRequest(http.MethodGet, videoStreamUrl, nil)
	setBiliHeader(req)
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	videoBytes, _ := io.ReadAll(resp.Body)
	return os.WriteFile(filepath.Join(savePath, title+"_video.m4s"), videoBytes, 0644)
}

func setBiliHeader(req *http.Request) {
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
	req.Header.Set("Cookie", config.GlobalConfig.Bili.Cookie)
	req.Header.Set("Referer", "https://www.bilibili.com/")
}

func mergeVideo(title, savePath string) error {
	videoFile := filepath.Join(savePath, title+"_video.m4s") // 视频流文件
	audioFile := filepath.Join(savePath, title+"_audio.m4s") // 音频流文件
	outputFile := filepath.Join(savePath, title+".mp4")      // 合并后的文件

	ffmpegPath, err := prepareFFMpegEnv()
	if err != nil {
		return err
	}

	// 使用 ffmpeg 合并音频和视频
	cmd := exec.Command(ffmpegPath, "-i", videoFile, "-i", audioFile, "-c", "copy", outputFile)
	err = cmd.Run()
	if err != nil {
		return err
	}

	_ = os.RemoveAll(videoFile)
	_ = os.RemoveAll(audioFile)
	return nil
}

func prepareFFMpegEnv() (string, error) {
	curlPath, err := exec.LookPath("ffmpeg")
	if err == nil {
		return curlPath, nil
	}

	_, err = os.Stat("ffmpeg/bin/ffmpeg.exe")
	if err == nil {
		return "ffmpeg/bin/ffmpeg.exe", nil
	}

	// 下载并解压ffmpeg
	if _, err = os.Stat("ffmpeg-essentials.zip"); os.IsNotExist(err) {
		err = tools.DownloadFile("https://www.gyan.dev/ffmpeg/builds/packages/ffmpeg-7.1-essentials_build.zip", "ffmpeg-essentials.zip")
		if err != nil {
			return "", err
		}
	}

	err = tools.Unzip("ffmpeg-essentials.zip", "./")
	if err != nil {
		return "", err
	}
	_ = os.Rename("ffmpeg-7.1-essentials_build", "ffmpeg")
	_ = os.RemoveAll("ffmpeg-essentials.zip")

	return "ffmpeg/bin/ffmpeg.exe", nil
}
