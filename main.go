package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/smtp"
	"strings"
	"time"

	"github.com/paulstuart/ping"
)

var (
	cfg = Config{}
)

// Config ...
type Config struct {
	MailSubject   string
	MailServer    string
	MailSender    string
	MailSenderPwd string
	MailReciver   string
	IPList        []string
}

// Host ...
type Host struct {
	IP              string
	Online          bool
	OnlineReported  bool
	OfflineReported bool
	IsRestore       bool
	Count           int
	Interval        time.Duration
	LastOnlineTime  time.Time
	RestoreTime     time.Time
	OfflineTime     time.Time
	OfflineTimes    int
	OfflineDuration time.Duration
}

// Watch ...
func (h *Host) Watch() {
	for {
		if ping.Ping(h.IP, 2) {
			// offline => online
			if !h.Online {
				h.Count = 0
				h.IsRestore = true
				h.RestoreTime = time.Now()
				h.Online = true
				h.OfflineReported = false
				h.OfflineDuration += time.Since(h.OfflineTime)
			}
			h.Count++
			h.Interval = 1 * time.Second

			if h.IsRestore {
				if time.Since(h.OfflineTime) > 60*time.Second {
					if h.Count > 10 {
						if !h.OnlineReported {
							h.OnlineReport()
						}
					}
				}
			}
			h.LastOnlineTime = time.Now()
		} else {
			// online => offline
			if h.Online {
				h.Count = 0
				h.Online = false
				h.IsRestore = false
				h.OnlineReported = false
				h.OfflineTime = time.Now()
				h.OfflineTimes++
			}
			h.Count++
			h.Interval = 10 * time.Second
			if h.Count > 3 {
				if !h.OfflineReported {
					h.OfflineReport()
				}
			}
		}
		time.Sleep(h.Interval)
	}
}

// OfflineReport ...
func (h *Host) OfflineReport() {
	msg := fmt.Sprintf("Host %s offline at %s",
		h.IP,
		h.OfflineTime.Format("2006-01-02 15:04:05"))
	err := SendReport(msg)
	if err == nil {
		h.OfflineReported = true
	} else {
		log.Println(err)
	}
}

// OnlineReport ...
func (h *Host) OnlineReport() {
	msg := fmt.Sprintf("Host %s offline at %s,restore online at %s,offline %d times,total %s.",
		h.IP,
		h.OfflineTime.Format("2006-01-02 15:04:05"),
		h.RestoreTime.Format("2006-01-02 15:04:05"),
		h.OfflineTimes,
		h.OfflineDuration)
	err := SendReport(msg)
	if err == nil {
		h.OnlineReported = true
	} else {
		log.Println(err)
	}
}

// LoadConfig ...
func LoadConfig(fn string, v interface{}) {
	jsonData, err := ioutil.ReadFile(fn)
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(jsonData, v)
	if err != nil {
		log.Fatal(err)
	}
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

func main() {
	LoadConfig("./config.json", &cfg)
	//log.Printf("%#v\n",cfg)

	iplist := cfg.IPList
	log.Println("service start...")
	for _, ip := range iplist {
		h := &Host{
			IP:             ip,
			Online:         true,
			Interval:       1 * time.Second,
			LastOnlineTime: time.Now(),
			OfflineTime:    time.Now(),
		}
		go h.Watch()
	}

	done := make(chan struct{})
	<-done

}
