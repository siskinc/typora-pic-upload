package download

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/siskinc/typora-pic-upload/internal/httpx"
	"github.com/siskinc/typora-pic-upload/internal/pathx"
	"github.com/siskinc/typora-pic-upload/internal/strx"
)

type Downloader struct {
	downloadedFilePathList []string
}

func (d *Downloader) Clear() {
	for _, filePath := range d.downloadedFilePathList {
		if !pathx.ExistsPath(filePath) {
			continue
		}
		err := os.Remove(filePath)
		if err != nil {
			logrus.Errorf("delete file: %s have an err: %v", filePath, err)
		}
	}
}

func (d *Downloader) isValidFileName(fileName string) bool {
	dotCount := 0
	for i := 0; i <= len(fileName); i++ {
		if fileName[i] == '.' {
			dotCount++
			continue
		}
		if !strx.IsNumberOrAlpha(fileName[i]) {
			return false
		}
	}
	return dotCount <= 1
}

func (d *Downloader) newFileName() string {
	return fmt.Sprintf("%d", time.Now().Unix())
}

func (d *Downloader) DownloadFile(url string) (string, error) {
	fileName, err := pathx.GenUploadFilePathFormURL(url)
	if err != nil {
		logrus.Errorf("generator file path from url have an err: %v, url: %s", err, url)
		return "", nil
	}
	if !d.isValidFileName(fileName) {
		fileName = d.newFileName()
	}
	client := httpx.GetClinet()
	resp, err := client.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	out, err := os.Create(fileName)
	defer out.Close()
	if err != nil {
		return "", err
	}
	io.Copy(out, bytes.NewReader(body))
	d.downloadedFilePathList = append(d.downloadedFilePathList, fileName)
	return fileName, nil
}
