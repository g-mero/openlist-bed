package handler

import (
	"openlist-bed/internal/storages"
	"openlist-bed/pkg/utils"
	"openlist-bed/pkg/vimage"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/spf13/viper"
)

func GetImage(c fiber.Ctx) error {
	webPath := c.Params("+")
	webp := c.Query("webp") == "true"

	storage := storages.NewOpenlistStorage()

	buf, err := storage.GetImg(webPath)

	if err != nil {
		return err
	}

	img, err := vimage.LoadFromBuffer(buf, utils.FilenameWithoutExt(webPath))
	if err != nil {
		return err
	}

	if viper.GetBool("auto_webp") {
		webp = strings.Contains(c.Get("Accept"), "webp")
	}

	buf, format, err := img.SmartCompress(webp)
	if err != nil {
		return err
	}

	c.Set("Content-Disposition", "inline")
	c.Set("Content-Type", vimage.FormatToContentType(format))
	return c.Send(buf)
}
