package sms

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

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

	s.logger.Printf("SendMessage:values %s\n", b)
	req, err := http.NewRequest("POST", s.url.String(), bytes.NewBuffer(b))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if len(content) > 0 {
		fmt.Printf("SendMessage:body %s\n", content)
	}
	s.logger.Printf("SendMessage:response\n\tstatus: %v\n\theader: %v\n\tbody: %s\n", resp.Status, resp.Header, content)

	return nil
}
