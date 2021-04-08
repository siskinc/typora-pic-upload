package filex

import (
	"crypto/md5"
	"encoding/hex"
	"io"
)

func GetMd5(file io.Reader) string {
	md5 := md5.New()
	io.Copy(md5, file)
	return hex.EncodeToString(md5.Sum(nil))
}
