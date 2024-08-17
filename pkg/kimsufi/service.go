package kimsufi

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

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

func (s *Service) GetAvailabilities(datacenters []string, planCode string) (*Availabilities, error) {
	u, err := url.Parse("/dedicated/server/datacenter/availabilities")
	q := u.Query()
	if len(datacenters) > 0 {
		q.Set("datacenters", strings.Join(datacenters, ","))
	}
	if planCode != "" {
		q.Set("planCode", planCode)
	}
	u.RawQuery = q.Encode()

	var availabilities Availabilities
	err = s.client.GetUnAuth(u.String(), &availabilities)
	if err != nil {
		return nil, err
	}

	return &availabilities, nil
}

func (s *Service) ListServers(ovhSubsidiary string) (*Catalog, error) {
	u, err := url.Parse("/order/catalog/public/eco")
	q := u.Query()
	q.Set("ovhSubsidiary", ovhSubsidiary)
	u.RawQuery = q.Encode()

	var catalog Catalog
	err = s.client.GetUnAuth(u.String(), &catalog)
	if err != nil {
		return nil, err
	}

	return &catalog, nil
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
