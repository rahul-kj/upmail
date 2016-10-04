package email

import (
	"fmt"
	"net/smtp"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/rahul-kj/checkup"
)

// Notifier sends an email notification when something is wrong.
type Notifier struct {
	// Recipient is the email address to send the notification to.
	Recipient string
	// Server is the email server.
	Server string
	// Sender is the email address to send the notification from.
	Sender string
	// Auth holds the authentication details for the email server.
	Auth smtp.Auth
}

// Notify checks the health status of the result and sends an email if
// something is not healthy.
func (n Notifier) Notify(results []checkup.Result) error {
	for _, r := range results {
		logrus.Debugf("%s is %s: sending email", r.Title, r.Status())
		if err := n.sendEmail(r); err != nil {
			return err
		}
	}

	return nil
}

func (n Notifier) sendEmail(result checkup.Result) error {
	// create the template
	body := fmt.Sprintf(`From: %s
To: %s
Subject: %s %s

Checkup run at %s for the link %s resulted in %s
`, n.Sender, n.Recipient, result.Title, result.Status(), time.Now().Format(time.UnixDate), result.Title, result.Status())

	// send the email
	if err := smtp.SendMail(n.Server, n.Auth, n.Sender, []string{n.Recipient}, []byte(body)); err != nil {
		return fmt.Errorf("Send mail failed: %v", err)
	}

	return nil
}
