package logger

import (
	"errors"
	"testing"

	"go.uber.org/zap"
)

func TestFieldsToJSON(t *testing.T) {
	json := fieldsToJSON([]zap.Field{
		zap.Int64("a", 123),
		zap.Float64("b", 1.23),
		zap.Bool("c", true),
		zap.Int("d", 123),
		zap.Error(errors.New("err")),
		zap.Any("any", []string{"a", "b"}),
	})

	t.Log(json.Fields)
}
