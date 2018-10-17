package logger

import (
	"github.com/go-xorm/core"
	"github.com/golang/glog"
)

type Logger struct {
	showSQL bool
}

func (s *Logger) Error(v ...interface{}) {
	glog.Error(v...)
}

// Errorf implement core.ILogger
func (s *Logger) Errorf(format string, v ...interface{}) {
	glog.Errorf(format, v...)
}

// Debug implement core.ILogger
func (s *Logger) Debug(v ...interface{}) {
	glog.Info(v...)
}

// Debugf implement core.ILogger
func (s *Logger) Debugf(format string, v ...interface{}) {
	glog.Infof(format, v...)
}

// Info implement core.ILogger
func (s *Logger) Info(v ...interface{}) {
	glog.Info(v...)
}

// Infof implement core.ILogger
func (s *Logger) Infof(format string, v ...interface{}) {
	glog.Infof(format, v...)
}

// Warn implement core.ILogger
func (s *Logger) Warn(v ...interface{}) {
	glog.Warning(v...)
}

// Warnf implement core.ILogger
func (s *Logger) Warnf(format string, v ...interface{}) {
	glog.Warningf(format, v...)
}

// Level implement core.ILogger
func (s *Logger) Level() core.LogLevel {
	return 1
}

// SetLevel implement core.ILogger
func (s *Logger) SetLevel(l core.LogLevel) {
	glog.V(glog.Level(l))
}

// ShowSQL implement core.ILogger
func (s *Logger) ShowSQL(show ...bool) {
	if len(show) == 0 {
		s.showSQL = true
		return
	}
	s.showSQL = show[0]
}

// IsShowSQL implement core.ILogger
func (s *Logger) IsShowSQL() bool {
	return s.showSQL
}
