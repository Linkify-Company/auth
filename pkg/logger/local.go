package logger

import (
	"auth/pkg/errify"
	"github.com/fatih/color"
)

type local struct{}

func (l *local) Info(err errify.IError) {
	color.Green(format(err))
}

func (l *local) Infof(format string, args ...interface{}) {
	color.Green(format, args...)
}

func (l *local) Error(err errify.IError) {
	color.Red(format(err))
}

func (l *local) Errorf(format string, args ...interface{}) {
	color.Red(format, args...)
}

func (l *local) Warn(err errify.IError) {
	color.Blue(format(err))
}

func (l *local) Warnf(format string, args ...interface{}) {
	color.Blue(format, args...)
}

func (l *local) Debug(err errify.IError) {
	color.Cyan(format(err))
}

func (l *local) Debugf(format string, args ...interface{}) {
	color.Cyan(format, args...)
}
