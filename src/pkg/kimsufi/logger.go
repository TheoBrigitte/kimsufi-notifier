package kimsufi

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

type Logger struct {
	*logrus.Logger
}

func NewRequestLogger(l *logrus.Logger) *Logger {
	if l == nil {
		// No-op logger
		l = &logrus.Logger{}
	}

	logger := &Logger{
		l,
	}

	return logger
}

func (l *Logger) LogRequest(r *http.Request) {
	l.Tracef("request: %s %s %s %v\n", r.Method, r.URL.String(), r.Proto, r.Header)
}

func (l *Logger) LogResponse(r *http.Response) {
	l.Tracef("response: %s %s %v\n", r.Status, r.Proto, r.Header)
}
