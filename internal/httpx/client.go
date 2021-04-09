package httpx

import (
	"net/http"
	"net/url"
	"sync"

	"github.com/spf13/viper"
)

var (
	client *http.Client
	once   = &sync.Once{}
)

func GetClinet() *http.Client {
	once.Do(func() {
		proxyUrl := viper.GetString("proxy")
		//设置代理
		var proxyURI *url.URL
		if proxyUrl != "" {
			proxyURI, _ = url.Parse(proxyUrl)
		}

		//构建代理client对象
		client = &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyURL(proxyURI),
			},
		}
	})
	return client
}
