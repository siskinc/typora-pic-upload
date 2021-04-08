package pathx

import (
	"fmt"
	"os"
	"strings"
)

func ExistsPath(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

func IsFile(path string) bool {
	return !IsDir(path)
}

func IsHttpFile(path string) bool {
	if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
		return true
	}
	return false
}

func GenUploadFilePathFormURL(url string) (string, error) {
	if !IsHttpFile(url) {
		return "", fmt.Errorf("URL is not valid")
	}
	index := strings.LastIndex(url, "/")
	uploadFilePath := url[index+1:]
	return uploadFilePath, nil
}
