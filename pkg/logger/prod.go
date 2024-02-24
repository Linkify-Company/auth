package logger

import (
	"auth/pkg/errify"
	"fmt"
	"log/slog"
)

type prod struct {
	logger *slog.Logger
}

func (d *prod) Info(err errify.IError) {
	d.logger.Info(format(err))
	saveErrify(err, LevelInfo)
}

func (d *prod) Infof(format string, args ...interface{}) {
	d.logger.Info(format, args...)
	save(fmt.Sprintf(format, args...), LevelInfo)
}

func (d *prod) Error(err errify.IError) {
	d.logger.Error(format(err))
	saveErrify(err, LevelError)
}

func (d *prod) Errorf(format string, args ...interface{}) {
	d.logger.Error(format, args...)
	save(fmt.Sprintf(format, args...), LevelError)
}

func (d *prod) Warn(err errify.IError) {
	d.logger.Warn(format(err))
	saveErrify(err, LevelWarning)
}

func (d *prod) Warnf(format string, args ...interface{}) {
	d.logger.Warn(format, args...)
	save(fmt.Sprintf(format, args...), LevelWarning)
}

func (d *prod) Debug(err errify.IError) {
	d.logger.Debug(format(err))
	saveErrify(err, LevelDebug)
}

func (d *prod) Debugf(format string, args ...interface{}) {
	d.logger.Debug(format, args...)
	save(fmt.Sprintf(format, args...), LevelDebug)
}
