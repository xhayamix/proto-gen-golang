package action

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewConfig(topicID string) zap.Config {
	path := fmt.Sprintf("pubsub://%s", topicID)
	return zap.Config{
		Level:            zap.NewAtomicLevelAt(zap.InfoLevel),
		Development:      false,
		Encoding:         "json",
		EncoderConfig:    encoderConfig,
		OutputPaths:      []string{path},
		ErrorOutputPaths: []string{path},
	}
}

var encoderConfig = zapcore.EncoderConfig{
	EncodeLevel:    zapcore.CapitalLevelEncoder,
	EncodeTime:     zapcore.RFC3339NanoTimeEncoder,
	LineEnding:     zapcore.DefaultLineEnding,
	EncodeDuration: zapcore.SecondsDurationEncoder,
	EncodeCaller:   zapcore.ShortCallerEncoder,
}
