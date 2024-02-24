package logger

import (
	"auth/pkg/errify"
	"fmt"
	"log/slog"
)

type dev struct {
	logger *slog.Logger
}

func (d *dev) Info(err errify.IError) {
	d.logger.Info(format(err))
	saveErrify(err, LevelInfo)
}

func (d *dev) Infof(format string, args ...interface{}) {
	d.logger.Info(format, args...)
	save(fmt.Sprintf(format, args...), LevelInfo)
}

func (d *dev) Error(err errify.IError) {
	d.logger.Error(format(err))
	saveErrify(err, LevelError)
}

func (d *dev) Errorf(format string, args ...interface{}) {
	d.logger.Error(format, args...)
	save(fmt.Sprintf(format, args...), LevelError)
}

func (d *dev) Warn(err errify.IError) {
	d.logger.Warn(format(err))
	saveErrify(err, LevelWarning)
}

func (d *dev) Warnf(format string, args ...interface{}) {
	d.logger.Warn(format, args...)
	save(fmt.Sprintf(format, args...), LevelWarning)
}

func (d *dev) Debug(err errify.IError) {
	d.logger.Debug(format(err))
	saveErrify(err, LevelDebug)
}

func (d *dev) Debugf(format string, args ...interface{}) {
	d.logger.Debug(format, args...)
	save(fmt.Sprintf(format, args...), LevelDebug)
}
