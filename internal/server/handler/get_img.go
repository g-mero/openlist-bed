package handler

import (
	"openlist-bed/internal/storages"
	"openlist-bed/pkg/cache"
	"openlist-bed/pkg/logger"
	"openlist-bed/pkg/utils"
	"openlist-bed/pkg/vimage"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/spf13/viper"
)

func GetImage(c fiber.Ctx) error {
	webPath := c.Params("+")
	webp := c.Query("webp") == "true"
	if viper.GetBool("auto_webp") {
		webp = strings.Contains(c.Get("Accept"), "webp")
	}
	cacheKey := "img_cache:" + webPath
	if webp {
		cacheKey += ":webp"
	} else {
		cacheKey += ":orig"
	}
	typeCacheKey := cacheKey + ":type"
	cachedBuf := cache.Get(cacheKey)
	cachedType := cache.Get(typeCacheKey)
	if cachedBuf != nil && cachedType != nil {
		logger.Debug("cache hit")
		c.Set("Content-Disposition", "inline")
		c.Set("Content-Type", string(cachedType))
		c.Set("Cache-Control", "public, max-age=31536000, immutable")
		return c.Send(cachedBuf)
	}

	storage := storages.NewOpenlistStorage()

	buf, err := storage.GetImg(webPath)

	if err != nil {
		return err
	}

	img, err := vimage.LoadFromBuffer(buf, utils.FilenameWithoutExt(webPath))
	if err != nil {
		return err
	}

	buf, format, err := img.SmartCompress(webp)
	if err != nil {
		return err
	}

	contentType := vimage.FormatToContentType(format)
	c.Set("Content-Disposition", "inline")
	c.Set("Content-Type", contentType)
	c.Set("Cache-Control", "public, max-age=31536000, immutable")
	cache.Set(cacheKey, buf)
	cache.Set(typeCacheKey, []byte(contentType))
	return c.Send(buf)
}
