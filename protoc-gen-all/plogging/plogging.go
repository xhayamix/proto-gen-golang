package plogging

import (
	"fmt"
	"os"
	"time"
)

const (
	TimeLayout = "15:04:05.000000000"
)

type Logger interface {
	Infof(format string, a ...interface{})
}

type logger struct{}

var defaultLogger = &logger{}

func GetLogger() Logger {
	return defaultLogger
}

func (l *logger) Infof(format string, a ...interface{}) {
	_, err := fmt.Fprintf(os.Stderr, time.Now().Format(TimeLayout)+" "+format, a...)
	if err != nil {
		panic(err)
	}
}
