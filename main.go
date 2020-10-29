package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/TheoBrigitte/kimsufi-notify/pkg/kimsufi"
	"github.com/TheoBrigitte/kimsufi-notify/pkg/sms"
)

func main() {
	var (
		country  string
		hardware string

		smsUser string
		smsPass string
	)

	flag.StringVar(&country, "country", "", "kimsufi country code (e.g. fr)")
	flag.StringVar(&hardware, "hardware", "", "kimsufi hardware code (e.g. 1801sk143)")
	flag.StringVar(&smsUser, "smsUser", "", "sms api username")
	flag.StringVar(&smsPass, "smsPass", "", "sms api password")
	flag.Parse()

	logger := log.New(os.Stderr, "", log.Lshortfile)

	logger.Println("start")
	d := kimsufi.Config{
		Logger: logger,
	}
	k, err := kimsufi.NewService(d)
	if err != nil {
		logger.Fatalf("error: %v\n", err)
	}

	a, err := k.GetAvailabilities(country, hardware)
	if err != nil {
		logger.Fatalf("error: %v\n", err)
	}
	logger.Printf("kimsufi:availabilities\n")
	formatter := kimsufi.DatacenterFormatter(kimsufi.IsDatacenterAvailable, kimsufi.DatacenterKey)
	result := a.Format(kimsufi.HardwareKey, kimsufi.RegionKey, formatter)
	data, err := json.Marshal(result)
	fmt.Printf("%s", data)
	logger.Println("end")

	os.Exit(0)

	c := sms.Config{
		URL:    "https://smsapi.free-mobile.fr/sendmsg",
		Logger: logger,
		User:   smsUser,
		Pass:   smsPass,
	}

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
