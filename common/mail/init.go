package mail

import (
	"bytes"
	"html/template"
	"log"
	"net/smtp"
	"os"
)

type client struct {
	auth         smtp.Auth
	host         string
	port         string
	username     string
	templatePath string
}

type Client interface {
	SendMail(templateName string, emailDest string, data AgreementMailReq) error
}

func Init() Client {
	return &client{
		auth:         smtp.PlainAuth("", os.Getenv("SMTP_USERNAME"), os.Getenv("SMTP_PASSWORD"), os.Getenv("SMTP_HOST")),
		host:         os.Getenv("SMTP_HOST"),
		port:         os.Getenv("SMTP_PORT"),
		username:     os.Getenv("SMTP_USERNAME"),
		templatePath: "./etc/template/",
	}
}

func (c *client) SendMail(templateName string, emailDest string, data AgreementMailReq) error {
	tmpl, err := template.ParseFiles(c.templatePath + templateName)
	if err != nil {
		log.Println("[SendMail] error loading template:", err.Error())
		return err
	}

	var body bytes.Buffer
	err = tmpl.Execute(&body, data)
	if err != nil {
		log.Println("[SendMail] error executing template:", err.Error())
		return err
	}

	// Compose email headers and body
	subject := "Subject: Thank You for Your Investment\r\n" +
		"MIME-version: 1.0;\r\n" +
		"Content-Type: text/html; charset=\"UTF-8\";\r\n\r\n"

	msg := []byte(subject + body.String())

	// Auth and send
	err = smtp.SendMail(c.host+":"+c.port, c.auth, c.username, []string{emailDest}, msg)
	if err != nil {
		log.Println("[SendMail] error sending email:", err.Error())
		return err
	}

	return nil
}
