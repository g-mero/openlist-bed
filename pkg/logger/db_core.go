package logger

import (
	"time"

	"go.uber.org/zap/zapcore"
)

// LogWriter database should implement this interface
type LogWriter interface {
	WriteLog(level, source, message string, detail map[string]interface{}, createdAt time.Time) error
}

type dbCore struct {
	writer LogWriter
	source string
}

func NewDBCore(writer LogWriter, source string) zapcore.Core {
	return &dbCore{
		writer: writer,
		source: source,
	}
}

func (c *dbCore) Enabled(lvl zapcore.Level) bool {
	// do not log debug
	return lvl >= zapcore.InfoLevel
}

func (c *dbCore) With(fields []zapcore.Field) zapcore.Core {
	return c
}

func (c *dbCore) Check(entry zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if c.Enabled(entry.Level) {
		return ce.AddCore(entry, c)
	}
	return ce
}

func (c *dbCore) Sync() error {
	return nil
}

func fieldsToJSON(fields []zapcore.Field) *zapcore.MapObjectEncoder {
	enc := zapcore.NewMapObjectEncoder()
	for _, f := range fields {
		f.AddTo(enc)
	}
	return enc
}

func (c *dbCore) Write(entry zapcore.Entry, fields []zapcore.Field) error {
	detail := fieldsToJSON(fields).Fields
	if entry.Caller.Defined {
		detail["caller"] = entry.Caller.TrimmedPath()
	}

	return c.writer.WriteLog(
		entry.Level.String(),
		c.source,
		entry.Message,
		detail,
		entry.Time,
	)
}
