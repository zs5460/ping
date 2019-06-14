package main

import (
	"fmt"
	"log"
	"time"

	"github.com/zs5460/mail"
	"github.com/zs5460/my"

	"github.com/paulstuart/ping"
)

var (
	cfg = config{}
)

type config struct {
	mail.Config
	IPList []string
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

// SendReport ...
func SendReport(msg string) error {
	return mail.SendMail(
		cfg.MailSender,
		cfg.MailSenderPwd,
		cfg.MailServer,
		cfg.MailReciver,
		cfg.MailSubject, msg)
}

func main() {
	my.MustLoadConfig("./config.json", &cfg)

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
