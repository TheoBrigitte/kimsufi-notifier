package main

import (
	"log"
	"os"

	"github.com/TheoBrigitte/kimsufi-notify/pkg/sms"
)

func main() {
	logger := log.New(os.Stdout, "", log.Lshortfile)

	c := sms.Config{}

	s, err := sms.NewService(c)
	if err != nil {
		logger.Fatalf("error: %v\n", err)
	}

	err = s.SendMessage("test")
	if err != nil {
		logger.Fatalf("error: %v\n", err)
	}
	logger.Printf("message sent\n")
}
