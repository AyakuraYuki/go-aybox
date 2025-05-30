package email

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net/smtp"
	"regexp"

	em "github.com/jordan-wright/email"
)

func IsValidAddress(emailAddr string) bool {
	re := regexp.MustCompile(`^([a-zA-Z0-9._%-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,})$`)
	return re.MatchString(emailAddr)
}

type SMTPConfig struct {
	Host     string // SMTP endpoint host
	Port     int    // SMTP endpoint port
	SSL      bool   // enable SSL or use TLS
	Username string // sender's email address
	Password string // sender's email password (auth code, or independent password)
	From     string // sender's screen name, e.g. Sender <sender@exampe.com>
}

func (c *SMTPConfig) Addr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

func SMTPSendEmail(ctx context.Context, conf *SMTPConfig, subject, content string, to, cc, bcc []string) (err error) {
	return sendEmail(ctx, conf, subject, content, "", to, cc, bcc)
}

func SMTPSendHTML(ctx context.Context, conf *SMTPConfig, subject, html string, to, cc, bcc []string) (err error) {
	return sendEmail(ctx, conf, subject, "", html, to, cc, bcc)
}

func sendEmail(ctx context.Context, conf *SMTPConfig, subject, content, html string, to, cc, bcc []string) (err error) {
	if conf == nil {
		return errors.New("smtp config not found")
	}
	if len(to) == 0 {
		return errors.New("missing receivers: to")
	}

	auth := smtp.PlainAuth("", conf.Username, conf.Password, conf.Host)

	e := em.NewEmail()
	e.From = conf.From
	e.Subject = subject

	if content != "" {
		e.Text = []byte(content)
	} else if html != "" {
		e.HTML = []byte(html)
	} else {
		return errors.New("missing content")
	}

	// to
	for _, v := range to {
		if IsValidAddress(v) {
			e.To = append(e.To, v)
		}
	}

	// cc
	if len(cc) > 0 {
		for _, v := range cc {
			if IsValidAddress(v) {
				e.Cc = append(e.Cc, v)
			}
		}
	}

	// bcc
	if len(bcc) > 0 {
		for _, v := range bcc {
			if IsValidAddress(v) {
				e.Bcc = append(e.Bcc, v)
			}
		}
	}

	if conf.SSL {
		err = e.SendWithTLS(conf.Addr(), auth, &tls.Config{ServerName: conf.Host})
	} else {
		err = e.Send(conf.Addr(), auth)
	}

	return err
}
