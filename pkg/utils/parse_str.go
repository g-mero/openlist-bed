package utils

import (
	"errors"
	"path"
	"strings"

	"github.com/bytedance/sonic"
	"github.com/duke-git/lancet/v2/convertor"
	"github.com/gofiber/utils/v2"
)

func ParseJsonStr[T any](s string) (T, error) {
	var zero T

	err := sonic.Unmarshal([]byte(s), &zero)
	if err != nil {
		return zero, err
	}
	return zero, nil
}

func ParseStr[T any](s string) (T, error) {
	var zero T

	switch any(zero).(type) {
	case string:
		return any(s).(T), nil
	case int:
		res, err := utils.ParseInt(s)
		if err != nil {
			return zero, err
		}
		return any(int(res)).(T), nil
	case int64:
		res, err := utils.ParseInt(s)
		if err != nil {
			return zero, err
		}
		return any(res).(T), nil
	case float64:
		res, err := convertor.ToFloat(s)
		if err != nil {
			return zero, err

		}
		return any(res).(T), nil
	case bool:
		if s == "true" || s == "\"true\"" {
			return any(true).(T), nil
		}
		if s == "false" || s == "\"false\"" {
			return any(false).(T), nil
		}
		return zero, errors.New("it is not a bool or bool-like")
	}

	return zero, errors.New("un-support type")
}

func FilenameWithoutExt(filepath string) string {
	return strings.Split(path.Base(filepath), ".")[0]
}
