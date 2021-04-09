package uploader

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/qiniu/go-sdk/v7/auth"
	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/siskinc/typora-pic-upload/internal/download"
	"github.com/siskinc/typora-pic-upload/internal/filex"
	"github.com/siskinc/typora-pic-upload/internal/pathx"
	"github.com/siskinc/typora-pic-upload/internal/urlx"
	"github.com/siskinc/typora-pic-upload/internal/httpx"
	"github.com/spf13/viper"
)

type QiniuUploader struct {
	AccessKey string
	SecretKey string
	Bucket    string
	Url       string
	Options   string
	Path      string
	Proxy     string
}

func CreateQiniuUploaderByConfig() Uploader {
	res := &QiniuUploader{}
	viper.UnmarshalKey("qiniu", res)
	return res
}

func (uld *QiniuUploader) useHTTPS() bool {
	return strings.HasPrefix(uld.Url, "https")
}

func (uld *QiniuUploader) skip(path string) bool {
	return strings.HasPrefix(path, uld.Url)
}

func (uld *QiniuUploader) Upload(path string) (string, error) {
	if uld.skip(path) {
		return path, nil
	}
	downloaderObj := download.Downloader{}
	defer downloaderObj.Clear()

	if pathx.IsHttpFile(path) {
		var err error
		path, err = downloaderObj.DownloadFile(path)
		if err != nil {
			return "", err
		}
	}

	f, err := os.Open(path)
	if err != nil {
		err = fmt.Errorf("open %s have an error: %v", path, err)
		return "", err
	}
	defer f.Close()
	filename := filex.GetMd5(f)
	key := urlx.Join(uld.Path, filename)
	if !strings.HasSuffix(key, ".png") {
		key = fmt.Sprintf("%s.png", key)
	}

	putPolicy := storage.PutPolicy{
		Scope: fmt.Sprintf("%s:%s", uld.Bucket, key),
	}

	mac := auth.New(uld.AccessKey, uld.SecretKey)
	upToken := putPolicy.UploadToken(mac)
	zone, _ := storage.GetRegion(uld.AccessKey, uld.Bucket)
	cfg := storage.Config{
		Zone:          zone,
		UseHTTPS:      uld.useHTTPS(),
		UseCdnDomains: false,
	}

	client := httpx.GetClinet()
	// 构建表单上传的对象
	formUploader := storage.NewFormUploaderEx(&cfg, &storage.Client{Client: client})
	ret := storage.PutRet{}
	err = formUploader.PutFile(context.Background(), &ret, upToken, key, path, nil)
	if err != nil {
		err = fmt.Errorf("put file to qiniu have an err: %v", err)
		return "", err
	}
	return urlx.Join(uld.Url, ret.Key), err
}
