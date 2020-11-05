package sms

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"

	"github.com/sirupsen/logrus"
)

type Service struct {
	client *http.Client
	logger logrus.StdLogger
	url    *url.URL
	user   string
	pass   string
}

type Config struct {
	URL    string
	Logger logrus.StdLogger
	User   string
	Pass   string
}

type Request struct {
	Msg  string `json:"msg"`
	User string `json:"user"`
	Pass string `json:"pass"`
}

func NewService(config Config) (*Service, error) {
	u, err := url.Parse(config.URL)
	if err != nil {
		return nil, err
	}

	s := &Service{
		client: &http.Client{},
		logger: config.Logger,
		url:    u,
		user:   config.User,
		pass:   config.Pass,
	}

	return s, nil
}

func (s *Service) SendMessage(msg string) error {
	r := Request{
		User: s.user,
		Pass: s.pass,
		Msg:  msg,
	}
	b, err := json.Marshal(r)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", s.url.String(), bytes.NewBuffer(b))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")

	s.LogRequest(req, r)

	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	s.LogResponse(resp)

	return nil
}

func (s *Service) LogRequest(r *http.Request, d Request) {
	d.User = strings.Repeat("*", len(d.User))
	d.Pass = strings.Repeat("*", len(d.Pass))
	s.logger.Printf("sms: %s %s %s %#v\n", r.Method, r.URL.String(), r.Proto, d)
}
func (s *Service) LogResponse(r *http.Response) {
	s.logger.Printf("sms: %s %s %v\n", r.Status, r.Proto, r.Header)
}
