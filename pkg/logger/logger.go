package logger

import (
	"auth/internal/config"
	"auth/pkg/errify"
	"fmt"
	"log/slog"
	"os"
)

type Logger interface {
	Info(err errify.IError)
	Infof(format string, args ...interface{})

	Error(err errify.IError)
	Errorf(format string, args ...interface{})

	Warn(err errify.IError)
	Warnf(format string, args ...interface{})

	Debug(err errify.IError)
	Debugf(format string, args ...interface{})
}

func GetLogger(env string) Logger {
	switch env {
	case config.EnvLocal:
		return &local{}
	case config.EnvDev:
		return &dev{
			logger: slog.New(
				slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
			),
		}
	default:
		return &prod{
			logger: slog.New(
				slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
			),
		}
	}
}

func format(err errify.IError) string {
	return fmt.Sprintf("{\n  Error: %s\n  Location: %s\n  Message: %s\n  Details: %s\n  Read: %s\n}\n\n",
		err.Error(), err.Location(), err.Message(), err.Details(), err.Read())
}

const (
	LevelInfo    = "info"
	LevelError   = "error"
	LevelWarning = "warning"
	LevelDebug   = "debug"
)
