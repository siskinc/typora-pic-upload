package download

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/siskinc/typora-pic-upload/internal/pathx"
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

func (d *Downloader) DownloadFile(url string) (string, error) {
	fileName, err := pathx.GenUploadFilePathFormURL(url)
	if err != nil {
		logrus.Errorf("generator file path from url have an err: %v, url: %s", err, url)
		return "", nil
	}
	resp, err := http.Get(url)
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
