package uploader

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/qiniu/go-sdk/v7/auth"
	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/siskinc/typora-pic-upload/internal/download"
	"github.com/siskinc/typora-pic-upload/internal/filex"
	"github.com/siskinc/typora-pic-upload/internal/pathx"
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

	//设置代理
	var proxyURI *url.URL
	if uld.Proxy != "" {
		proxyURI, _ = url.Parse(uld.Proxy)
	}

	//构建代理client对象
	client := http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyURI),
		},
	}
	// 构建表单上传的对象
	formUploader := storage.NewFormUploaderEx(&cfg, &storage.Client{Client: &client})
	ret := storage.PutRet{}
	err = formUploader.PutFile(context.Background(), &ret, upToken, key, filePath, nil)
	if err != nil {
		err = fmt.Errorf("put file to qiniu have an err: %v", err)
		return "", err
	}
	return urlx.Join(uld.Url, ret.Key), err
}
