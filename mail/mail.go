package mail

import "gopkg.in/gomail.v1"

type Mail struct {
	SmtpServer     string
	Port           int
	SenderEmail    string
	SenderName     string
	SenderPassword string
	ToEmail        string
	Subject        string
	Message        string
}

func (m Mail) SendEmail() error {
	msg := gomail.NewMessage()
	msg.SetAddressHeader("From", m.SenderEmail, m.SenderName)
	msg.SetAddressHeader("To", m.ToEmail, "")
	msg.SetHeader("Subject", m.Subject)
	msg.SetBody("text/html", m.Message)

	mailer := gomail.NewMailer(m.SmtpServer, m.SenderEmail, m.SenderPassword, m.Port)
	return mailer.Send(msg)
}
