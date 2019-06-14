# mail

## Install
```shell
go get github.com/zs5460/mail
```

## Usage
```go

cfg := mail.Config{
	MailSubject:   "test",
	MailServer:    "smtp.xxx.com:25",
	MailSender:    "xxx@xxx.com",
	MailSenderPwd: "******",
	MailReciver:   "abc@abc.com",
}

err := mail.SendMail(
	cfg.MailSender,
	cfg.MailSenderPwd,
	cfg.MailServer,
	cfg.MailReciver,
	cfg.MailSubject,
	"this is a test mail",
)
if err != nil {
	log.Println(err)
}

```

### Licence

Released under MIT license, see [LICENSE](LICENSE) for details.