package storages

import (
	"openlist-bed/pkg/openlist_sdk"
	"openlist-bed/pkg/vimage"

	"github.com/imroc/req/v3"
	"github.com/spf13/viper"
)

type OpenlistStorage struct {
	api  *openlist_sdk.OpenlistApi
	path string
}

func (that *OpenlistStorage) SaveImg(img *vimage.Image) (ImageUrl, error) {
	dir := MakeDateDir()
	err := that.api.UploadImg(that.path+"/"+dir, img)
	if err != nil {
		return ImageUrl{}, err
	}

	imgUrl := MakeImgUrl(dir + img.FullName())
	return imgUrl, nil
}

func (that *OpenlistStorage) GetImg(path string) ([]byte, error) {
	openlistInfo, err := that.api.GetImgInfo(that.path + "/" + path)
	if err != nil {
		return nil, err
	}

	resp, err := req.C().R().Get(openlistInfo.RawUrl)
	if err != nil {
		return nil, err
	}

	return resp.Bytes(), nil
}

func NewOpenlistStorage() Storage {
	token := viper.GetString("openlist.token")
	host := viper.GetString("openlist.host")
	path := viper.GetString("openlist.path")

	return &OpenlistStorage{
		api:  openlist_sdk.NewOpenlistApi(token, host),
		path: path,
	}
}
