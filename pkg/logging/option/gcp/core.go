package gcp

import (
	"runtime"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type core struct {
	zapcore.Core
}

func NewCore(c zapcore.Core) zapcore.Core {
	return &core{
		Core: c,
	}
}

//nolint:gocritic // hugeParamだがzapcore.Coreの実装なので仕方ない
func (c *core) Check(ent zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if c.Enabled(ent.Level) {
		return ce.AddCore(ent, c)
	}
	return ce
}

//nolint:gocritic // hugeParamだがzapcore.Coreの実装なので仕方ない
func (c *core) Write(ent zapcore.Entry, fields []zapcore.Field) error {
	// source locationを追加
	if ent.Caller.Defined {
		loc := &SourceLocation{
			File: ent.Caller.File,
			Line: ent.Caller.Line,
		}
		if fn := runtime.FuncForPC(ent.Caller.PC); fn != nil {
			loc.Function = fn.Name()
		}
		fields = append(fields, zap.Object(loc.Key(), loc))
	}

	return c.Core.Write(ent, fields)
}
