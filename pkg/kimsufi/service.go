package kimsufi

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/ovh/go-ovh/ovh"
	"github.com/sirupsen/logrus"
)

type Config struct {
	URL    string
	Logger logrus.StdLogger
}

type Service struct {
	url    *url.URL
	client *ovh.Client
}

func NewService(config Config) (*Service, error) {
	c, err := ovh.NewClient(config.URL, "none", "none", "none")
	//c, err := ovh.NewEndpointClient(ovh.KimsufiEU)
	if err != nil {
		fmt.Println("nope")
		return nil, err
	}
	c.Logger = NewRequestLogger(config.Logger)

	s := &Service{
		client: c,
	}

	return s, nil
}

func (s *Service) GetAvailabilities(countryCode, hardware string) (*Availabilities, error) {
	u, err := url.Parse("/dedicated/server/availabilities")
	q := u.Query()
	q.Set("country", countryCode)
	if hardware != "" {
		q.Set("hardware", hardware)
	}
	u.RawQuery = q.Encode()

	var availabilities Availabilities
	err = s.client.GetUnAuth(u.String(), &availabilities)
	if err != nil {
		return nil, err
	}

	return &availabilities, nil
}

type Logger struct {
	logger logrus.StdLogger
}

func NewRequestLogger(l logrus.StdLogger) *Logger {
	logger := &Logger{
		logger: l,
	}

	return logger
}

func (l *Logger) LogRequest(r *http.Request) {
	l.logger.Printf("kimsufi: %s %s %s %v\n", r.Method, r.URL.String(), r.Proto, r.Header)
}

func (l *Logger) LogResponse(r *http.Response) {
	l.logger.Printf("kimsufi: %s %s %v\n", r.Status, r.Proto, r.Header)
}
