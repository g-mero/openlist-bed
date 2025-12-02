package storages

import (
	"net/url"
	"openlist-bed/pkg/vimage"
	path2 "path"
	"time"

	"github.com/duke-git/lancet/v2/netutil"
	"github.com/spf13/viper"
)

type Storage interface {
	SaveImg(img *vimage.Image) (ImageUrl, error)
	GetImg(path string) ([]byte, error)
}

type ImageUrl struct {
	Url  string `json:"url"`
	Path string `json:"path"`
}

// MakeImgUrl basePath： eg 20230627/test.jpg
func MakeImgUrl(basePath string) (imgUrl ImageUrl) {
	joinPath, _ := url.JoinPath(viper.GetString("host"), "pic", basePath)
	imgUrl.Url, _ = netutil.EncodeUrl(joinPath)
	imgUrl.Path = path2.Clean(basePath)

	return imgUrl
}

// GetFileNameFromPath 获取路径中图片的文件名
func GetFileNameFromPath(path string) string {
	base := path2.Base(path)

	return base[:len(base)-len(path2.Ext(base))]
}

// MakeDateDir 生成由 YYYY/MM/DD 组成的文件夹路径名
func MakeDateDir() string {
	now := time.Now()

	return now.Format("2006/01/02") + "/"
}
