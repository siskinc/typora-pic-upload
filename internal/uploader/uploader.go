package uploader

import (
	"fmt"
	"strings"
)

const (
	UploaderTypeQiniu = "qiniu"
)

type Uploader interface {
	Upload(path string) (string, error)
}

type UploaderFactory struct {
}

func (factory *UploaderFactory) CreateUploader(uploaderType string) Uploader {
	switch strings.ToLower(uploaderType) {
	case UploaderTypeQiniu:
		return CreateQiniuUploaderByConfig()
	default:
		panic(fmt.Errorf("can not find uploader: %s", uploaderType))
	}
}
