package handler

import (
	"fmt"
	"openlist-bed/internal/storages"
	"openlist-bed/pkg/vimage"
	"time"

	"github.com/gofiber/fiber/v3"
)

func UploadImg(c fiber.Ctx) error {
	file, err := c.FormFile("image")
	compress := c.Query("compress") == "true"
	keepName := c.Query("keep_name") == "true"

	if err != nil {
		return err
	}

	img, err := vimage.LoadFromFile(file)
	if err != nil {
		return err
	}

	if !keepName {
		img, err = vimage.LoadFromBuffer(img.OriginalData(), fmt.Sprintf("%d_%dx%d",
			time.Now().Unix(), img.Width(), img.Height()))
		if err != nil {
			return err
		}
	}

	if compress {
		buf, _, err := img.SmartCompress(false)
		if err != nil {
			return err
		}
		img, err = vimage.LoadFromBuffer(buf, img.FileName)
		if err != nil {
			return err
		}
	}

	storage := storages.NewOpenlistStorage()
	url, err := storage.SaveImg(img)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"url":  url.Url,
		"path": url.Path,
	})
}
