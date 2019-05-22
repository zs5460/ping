package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"sync"

	"github.com/paulstuart/ping"
)

// Config ...
type Config struct {
	MailServer    string
	MailSender    string
	MailSenderPwd string
	MailReciver   []string
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

// PingHost ...
func PingHost(ip string) {
	if ping.Ping(ip, 2) {
		fmt.Printf("ping %s success!\n", ip)
	} else {
		fmt.Printf("ping %s failed!\n", ip)
		//log to file

		// send mail
	}
}

func main() {

	cfg := ReadConfig()
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
