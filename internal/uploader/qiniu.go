package uploader

import (
	"context"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/qiniu/go-sdk/v7/auth"
	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/siskinc/typora-pic-upload/internal/download"
	"github.com/siskinc/typora-pic-upload/internal/filex"
	"github.com/siskinc/typora-pic-upload/internal/httpx"
	"github.com/siskinc/typora-pic-upload/internal/pathx"
	"github.com/siskinc/typora-pic-upload/internal/strx"
	"github.com/siskinc/typora-pic-upload/internal/urlx"
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

func (uld *QiniuUploader) fixFileNameSuffix(suffix string) string {
	if suffix == "" {
		return suffix
	}
	endIndex := 1
	for i := 1; i < len(suffix); i++ {
		endIndex = i
		if !strx.IsAlpha(suffix[i]) {
			endIndex -= 1
			break
		}
	}
	return suffix[:endIndex+1]
}

func (uld *QiniuUploader) Upload(filePath string) (string, error) {
	if uld.skip(filePath) {
		return filePath, nil
	}
	downloaderObj := download.Downloader{}
	defer downloaderObj.Clear()

	if pathx.IsHttpFile(filePath) {
		var err error
		filePath, err = downloaderObj.DownloadFile(filePath)
		if err != nil {
			return "", err
		}
	}

	f, err := os.Open(filePath)
	if err != nil {
		err = fmt.Errorf("open %s have an error: %v", filePath, err)
		return "", err
	}
	defer f.Close()
	filename := filex.GetMd5(f)
	key := urlx.Join(uld.Path, filename)
	fileSuffix := path.Ext(filePath)
	fileSuffix = uld.fixFileNameSuffix(fileSuffix)
	if fileSuffix == "" {
		fileSuffix = ".png"
	}
	key = fmt.Sprintf("%s%s", key, fileSuffix)

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
	err = formUploader.PutFile(context.Background(), &ret, upToken, key, filePath, nil)
	if err != nil {
		err = fmt.Errorf("put file to qiniu have an err: %v", err)
		return "", err
	}
	return urlx.Join(uld.Url, ret.Key), err
}
