package log

import (
	"code_game/config"
	"os"

	"github.com/sirupsen/logrus"
)

const (
	XB3TraceID = "X-B3-Traceid"
)

type TraceIDHook struct {
	levels []logrus.Level
}

func (*TraceIDHook) Fire(e *logrus.Entry) error {
	if e.Context == nil {
		return nil
	}
	e.WithField("trace_id", e.Context.Value(XB3TraceID))
	return nil
}

func (h *TraceIDHook) WithLevel(level ...logrus.Level) *TraceIDHook {
	h.levels = append(h.levels, level...)
	return h
}

func (h *TraceIDHook) Levels() []logrus.Level {
	return h.levels
}

func InitLogger(cfg config.LogCfg) error {
	logrus.SetReportCaller(true)
	logrus.AddHook(new(TraceIDHook).WithLevel(
		logrus.DebugLevel,
		logrus.InfoLevel,
		logrus.WarnLevel,
		logrus.ErrorLevel,
		logrus.FatalLevel,
		logrus.PanicLevel,
	))
	if cfg.Path != "" {
		logFile, err := os.OpenFile(cfg.Path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			return err
		}
		logrus.SetOutput(logFile)
	}

	switch cfg.Format {
	case "text":
		logrus.SetFormatter(new(logrus.TextFormatter))
	default:
		logrus.SetFormatter(new(logrus.JSONFormatter))
	}

	if cfg.Level != "" {
		lvl, _ := logrus.ParseLevel(cfg.Level)
		logrus.SetLevel(lvl)
	}
	return nil
}
