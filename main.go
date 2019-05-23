package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/smtp"
	"strings"
	"sync"

	"github.com/paulstuart/ping"
)

var cfg *Config

// Config ...
type Config struct {
	MailSubject   string
	MailServer    string
	MailSender    string
	MailSenderPwd string
	MailReciver   string
	IPList        []string
}

// ReadConfig returns config
func ReadConfig() *Config {
	jsb, err := ioutil.ReadFile("./config.json")
	if err != nil {
		log.Fatal("read config.json error.")
	}
	var config Config
	err = json.Unmarshal(jsb, &config)
	if err != nil {
		log.Fatal("Unmarshal config.json error.")
	}
	return &config
}

// SendMail ...
func SendMail(user, password, host, to, subject, body string) error {
	hp := strings.Split(host, ":")
	auth := smtp.PlainAuth("", user, password, hp[0])
	contentType := "Content-Type: text/html; charset=UTF-8"
	msg := []byte("To: " + to + "\r\nFrom: " + user +
		"\r\nSubject: " + subject + "\r\n" + contentType + "\r\n\r\n" + body)
	sendTo := strings.Split(to, ";")
	err := smtp.SendMail(host, auth, user, sendTo, msg)
	return err
}

// SendReport ...
func SendReport(msg string) error {
	return SendMail(
		cfg.MailSender,
		cfg.MailSenderPwd,
		cfg.MailServer,
		cfg.MailReciver,
		cfg.MailSubject, msg)
}

// PingHost ...
func PingHost(ip string) {
	if ping.Ping(ip, 2) {
		fmt.Printf("ping %s success!\n", ip)
	} else {
		fmt.Printf("ping %s timeout!\n", ip)
		//log to file

		// send mail
		err := SendReport(ip)
		if err != nil {
			log.Println(err)
		}
	}
}

func main() {

	cfg = ReadConfig()
	fmt.Printf("%#v\n", cfg)

	iplist := cfg.IPList
	var wg sync.WaitGroup
	for _, ip := range iplist {
		wg.Add(1)
		go func(ip string) {
			defer wg.Done()
			PingHost(ip)
		}(ip)

	}
	wg.Wait()
}
