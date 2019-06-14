package mail

import (
	"crypto/tls"
	"net/smtp"
	"strings"
)

type Config struct {
	MailSubject   string
	MailServer    string
	MailSender    string
	MailSenderPwd string
	MailReciver   string
}

// SendMail ...
func SendMail(user, password, addr, to, subject, body string) error {
	hp := strings.Split(addr, ":")
	host := hp[0]
	port := hp[1]
	auth := smtp.PlainAuth("", user, password, host)
	contentType := "Content-Type: text/html; charset=UTF-8"
	msg := []byte("To: " + to + "\r\nFrom: " + user +
		"\r\nSubject: " + subject + "\r\n" + contentType + "\r\n\r\n" + body)
	sendTo := strings.Split(to, ";")
	if port == "465" {
		return sendMailWithTLS(addr, auth, user, sendTo, msg)
	}
	return smtp.SendMail(addr, auth, user, sendTo, msg)
}

func sendMailWithTLS(addr string, auth smtp.Auth, user string, sendTo []string, msg []byte) error {

	conn, err := tls.Dial("tcp", addr, nil)
	if err != nil {
		return err
	}

	host := strings.Split(addr, ":")[0]
	c, err := smtp.NewClient(conn, host)
	if err != nil {
		return err
	}
	defer c.Close()

	if auth != nil {
		if ok, _ := c.Extension("AUTH"); ok {
			if err = c.Auth(auth); err != nil {
				return err
			}
		}
	}

	if err = c.Mail(user); err != nil {
		return err
	}

	for _, addr := range sendTo {
		if err = c.Rcpt(addr); err != nil {
			return err
		}
	}

	w, err := c.Data()
	if err != nil {
		return err
	}

	_, err = w.Write(msg)
	if err != nil {
		return err
	}

	err = w.Close()
	if err != nil {
		return err
	}
	return c.Quit()
}
