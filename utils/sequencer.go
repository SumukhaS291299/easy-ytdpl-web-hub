package utils

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	litedb "ytdpldownloader/liteDB"
	"ytdpldownloader/runner"

	"github.com/charmbracelet/log"
)

var globalConfig AppData

func SaveConfig(config AppData) {
	globalConfig = config
}

func GetConfig() AppData {
	return globalConfig
}

func versionCheck() {
	version := litedb.ViewData([]byte("Version"))
	log.Info("Current yt-dpl Version :\t" + string(version))
	litedb.WriteData([]byte("version_last_checked"), []byte(time.Now().Local().String()))
}

func CheckUpdates() {
	log.Debug("Loading Database values...")
	versionCheck()
}

type FormatInfo struct {
	ID         string  `json:"id"`
	Resolution string  `json:"resolution"`
	Ext        string  `json:"ext"`
	VideoCodec string  `json:"video_codec"`
	AudioCodec string  `json:"audio_codec"`
	FPS        float64 `json:"fps,omitempty"`
	Size       any     `json:"size,omitempty"`
}

func (f FormatInfo) FormattedSize() string {
	if f.Size == nil {
		return "N/A"
	}
	var bytes float64
	switch v := f.Size.(type) {
	case float64:
		bytes = v
	case int:
		bytes = float64(v)
	case int32:
		bytes = float64(v)
	case int64:
		bytes = float64(v)
	default:
		return "N/A"
	}
	if bytes <= 0 {
		return "N/A"
	}
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%.0f B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.2f %cB", bytes/float64(div), "KMGTPE"[exp])
}

func Probeytdpl(url string) []byte {
	log.Debug("Got URL:\t", url)
	data := GetConfig()
	log.Debug("yt-dpl path", data.YTdplBinPath)
	log.Debug("URL received:\t", url)
	command := data.YTdplBinPath + " -s -O \"%(formats)j\" -- " + url
	log.Debug("Running command", command)
	ytdplProbestdout, ytdplProbestderr := runner.Run(command)
	var wg sync.WaitGroup
	wg.Add(1)
	outJSONStruct := YtdplProbe{}
	go func() {
		defer wg.Done()
		errstr := string(<-ytdplProbestderr)
		pullOut := <-ytdplProbestdout
		log.Debug(string(pullOut))
		if len(strings.TrimSpace(errstr)) > 0 {
			log.Error("[Warn]YT DPL has an error\n:", "yt_dpl error\t:", errstr)
		}
		json.Unmarshal(pullOut, &outJSONStruct)
	}()
	wg.Wait()
	// out, _ := json.Marshal(outJSONStruct)
	formats := make([]FormatInfo, 0, len(outJSONStruct))
	for _, f := range outJSONStruct {
		formats = append(formats, FormatInfo{
			ID:         f.FormatID,
			Resolution: f.Resolution,
			Ext:        f.Ext,
			VideoCodec: f.Vcodec,
			AudioCodec: f.Acodec,
			FPS:        f.Fps,
			Size:       f.Filesize,
		})
	}
	returnData, err := json.Marshal(formats)
	if err != nil {
		log.Error(err)
		return nil
	}
	return returnData
}

func Downloadytdpl(url, format string) {
	log.Debug("Got URL:\t", url)
	log.Debug("Got Format:\t", format)
	data := GetConfig()
	log.Debug("yt-dpl path", data.YTdplBinPath)
	log.Debug("URL received:\t", url)
	command := data.YTdplBinPath + " -f	 " + format + " -- " + url
	log.Debug("Running command", command)
	ytdplDownloaddout, ytdplDownloadstderr := runner.Run(command)
	var wg sync.WaitGroup
	// var ytOut, ytErr, ffOut, ffErr, fpOut, fpErr []byte
	wg.Add(1)
	go func() {
		defer wg.Done()
		errstr := string(<-ytdplDownloadstderr)
		log.Debug(<-ytdplDownloaddout)
		log.Debug(string(errstr))
		if len(strings.TrimSpace(errstr)) > 0 {
			log.Error("[Warn]YT DPL has an error\n:", "yt_dpl error\t:", errstr)
		}
	}()
	wg.Wait()
}
