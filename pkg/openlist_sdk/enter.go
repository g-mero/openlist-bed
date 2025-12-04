package openlist_sdk

import (
	"errors"
	"openlist-bed/pkg/logger"
	"openlist-bed/pkg/utils"
	"openlist-bed/pkg/vimage"
	path2 "path"

	"github.com/imroc/req/v3"
	"go.uber.org/zap"
)

type OpenlistApi struct {
	Host   string
	client *req.Client
}

type OpenlistInfo struct {
	RawUrl string
}

func NewOpenlistApi(token, host string) *OpenlistApi {
	client := req.C().SetCommonHeader("Authorization", token)

	return &OpenlistApi{
		host,
		client,
	}
}

func (that OpenlistApi) GetImgInfo(remotePath string) (OpenlistInfo, error) {
	var (
		err     error
		apiUrl  = that.Host + "/api/fs/get"
		imgInfo OpenlistInfo
	)
	var body utils.H
	_, err = that.client.R().
		SetBody(utils.H{"path": remotePath}).
		SetSuccessResult(&body).
		Post(apiUrl)

	if err != nil {
		return imgInfo, err
	}

	if body["code"].(float64) != 200 {
		msg := body["message"].(string)
		return imgInfo, errors.New(msg)
	}

	imgInfo.RawUrl = body["data"].(map[string]interface{})["raw_url"].(string)

	return imgInfo, nil
}

func (that OpenlistApi) UploadImg(remoteDir string, img *vimage.Image) error {
	header := map[string]string{
		"File-Path": path2.Clean(remoteDir + "/" + img.FullName()),
	}

	apiUrl := that.Host + "/api/fs/form"
	var body utils.H
	_, err := that.client.R().SetHeaders(header).SetFileBytes("file", img.FullName(), img.OriginalData()).
		SetSuccessResult(&body).Put(apiUrl)

	if err != nil {
		return err
	}

	if body["code"].(float64) != 200 {
		return errors.New("上传失败: " + body["message"].(string))
	}

	// refresh dir
	_ = that.refreshDir(remoteDir)

	return nil
}

func (that OpenlistApi) refreshDir(remoteDir string) error {
	apiUrl := that.Host + "/api/fs/list"

	resp, err := that.client.R().SetBodyJsonMarshal(utils.H{
		"path":     remoteDir,
		"password": "",
		"refresh":  true,
		"page":     1,
		"per_page": 0,
	}).Post(apiUrl)

	if err != nil {
		logger.Debug("[openlist] refresh dir error", zap.Error(err))
		return err
	}

	logger.Debug("[openlist] refresh dir response", zap.String("resp", resp.String()))

	return nil
}
