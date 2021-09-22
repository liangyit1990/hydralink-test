package monitor

import (
	"go.uber.org/zap"
)

// Logger wraps underlying logger
// The underlying logger can be either logrus or uber zap
type Logger struct {
	zapLogger *zap.SugaredLogger
}

// NewLogger creates new logger
func NewLogger() Logger {
	z, _ := zap.NewProduction()
	return Logger{
		zapLogger: z.Sugar(),
	}
}

// Infof uses fmt.Sprintf to log a templated message.
func (l Logger) Infof(template string, args ...interface{}) {
	if len(args) == 0 {
		l.zapLogger.Info(template)
		return
	}
	l.zapLogger.Infof(template, args)
}

// Errorf uses fmt.Sprintf to log a templated message.
func (l Logger) Errorf(template string, args ...interface{}) {
	if len(args) == 0 {
		l.zapLogger.Error(template)
		return
	}
	l.zapLogger.Errorf(template, args)
}

// Flush flushes any logs in buffer out
func (l Logger) Flush() {
	_ = l.zapLogger.Sync()
}
