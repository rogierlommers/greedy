package reminders

import (
	"fmt"
	"net/http"
	"net/smtp"

	"github.com/Sirupsen/logrus"
	"github.com/rogierlommers/greedy/internal/common"
	"github.com/rogierlommers/greedy/internal/render"
)

var r *Reminder

type Reminder struct {
	ToEmail      string
	FromEmail    string
	SMTPHost     string
	SMTPUser     string
	SMTPPassword string
}

func NewClient() {

	r = &Reminder{
		ToEmail:      common.ToEmail,
		FromEmail:    common.FromEmail,
		SMTPHost:     common.SMTPHost,
		SMTPUser:     common.SMTPUser,
		SMTPPassword: common.SMTPPassword,
	}

	logrus.Infof("reminders initialized for %s", r.ToEmail)
}

// AddReminderByURL adds a reminder
func AddReminderByURL(w http.ResponseWriter, r *http.Request) {
	urlTosend := r.FormValue("url")

	if len(urlTosend) == 0 || urlTosend == "about:blank" {
		renderObject := map[string]interface{}{
			"IsErrorPage":  "true",
			"errorMessage": "unable to insert empty or about:blank page",
		}
		render.DisplayPage(w, r, renderObject)
		return
	}

	err := sendReminder(urlTosend)
	if err != nil {
		renderObject := map[string]interface{}{
			"IsErrorPage":  "true",
			"errorMessage": err,
		}

		render.DisplayPage(w, r, renderObject)
		return
	}

	renderObject := map[string]interface{}{
		"IsConfirmation": "true",
		"hostname":       urlTosend,
	}

	render.DisplayPage(w, r, renderObject)
}

// AddReminderByText adds a reminder
func AddReminderByText(w http.ResponseWriter, r *http.Request) {
	text := r.FormValue("text")

	if len(text) == 0 || text == "about:blank" {
		renderObject := map[string]interface{}{
			"IsErrorPage":  "true",
			"errorMessage": "unable to remind empty text or about:blank page",
		}
		render.DisplayPage(w, r, renderObject)
		return
	}

	err := sendReminder(text)
	if err != nil {
		renderObject := map[string]interface{}{
			"IsErrorPage":  "true",
			"errorMessage": err,
		}

		render.DisplayPage(w, r, renderObject)
		return
	}

	renderObject := map[string]interface{}{
		"IsConfirmation": "true",
		"hostname":       text,
	}

	render.DisplayPage(w, r, renderObject)
}

func sendReminder(urlTosend string) error {

	auth := smtp.PlainAuth("", r.SMTPUser, r.SMTPPassword, r.SMTPHost)

	from := fmt.Sprintf("From: <%s>\r\n", r.FromEmail)
	to := fmt.Sprintf("To: <%s>\r\n", r.ToEmail)
	subject := "Subject: Reminder: " + urlTosend + "\r\n"
	body := "Reminder\r\n" + urlTosend + "\r\n"
	msg := from + to + subject + "\r\n" + body

	hostPort := fmt.Sprintf("%s:%d", r.SMTPHost, 587)

	err := smtp.SendMail(hostPort, auth, r.ToEmail, []string{r.ToEmail}, []byte(msg))
	if err != nil {
		return err
	}

	return nil
}
