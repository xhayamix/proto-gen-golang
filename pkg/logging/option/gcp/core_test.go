package gcp

import (
	"bytes"
	"encoding/json"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func newCore(writer io.Writer) zapcore.Core {
	enc := zapcore.NewJSONEncoder(encoderConfig)
	c := zapcore.NewCore(enc, zapcore.AddSync(writer), zapcore.DebugLevel)
	return NewCore(c)
}

func TestNewCore(t *testing.T) {
	writer := bytes.NewBuffer(nil)
	c := newCore(writer)
	logger := zap.New(c, zap.AddCaller())

	logger.Info("test")

	entry := make(map[string]interface{})
	err := json.Unmarshal(writer.Bytes(), &entry)
	assert.NoError(t, err)
	assert.Equal(t, "INFO", entry["severity"])
	assert.Equal(t, "test", entry["message"])
	assert.NotNil(t, entry["time"])
	locKey := (&SourceLocation{}).Key()
	if assert.NotNil(t, entry[locKey]) {
		loc := entry[locKey].(map[string]interface{})
		assert.Equal(t, float64(26), loc["line"])
		assert.True(t, strings.HasSuffix(loc["file"].(string), "core_test.go"))
		assert.True(t, strings.HasSuffix(loc["function"].(string), "TestNewCore"))
	}
}
