package kimsufi

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

// Logger is a wrapper around logrus.Logger.
// It is used to log requests and responses.
type Logger struct {
	*logrus.Logger
}

// NewRequestLogger creates a new Logger for requests.
// If l is nil, a no-op logger is used.
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

// LogRequest logs the HTTP request.
// Format: request: method url proto header
func (l *Logger) LogRequest(r *http.Request) {
	l.Tracef("request: %s %s %s %v\n", r.Method, r.URL.String(), r.Proto, r.Header)
}

// LogResponse logs the HTTP response.
// Format: response: status proto header
func (l *Logger) LogResponse(r *http.Response) {
	l.Tracef("response: %s %s %v\n", r.Status, r.Proto, r.Header)
}
