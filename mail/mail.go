package mail

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/mail"
	"net/smtp"
	"os"
)

type Message struct {
	To      string
	From    string
	Subject string
	Body    string
}

type Email struct {
	To    mail.Address
	From  mail.Address
	Email string
}

// Sends a message using SMTP over TLS
func SendMessage(message *Message) {
	to := mail.Address{"", message.To}
	from := mail.Address{"", message.From}

	// Set up headers
	headers := make(map[string]string)
	headers["From"] = from.String()
	headers["To"] = to.String()
	headers["Subject"] = message.Subject

	// Set Up message
	body := ""
	for k, v := range headers {
		body += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	body += "\r\n" + message.Body

	email := Email{
		To:    to,
		From:  from,
		Email: body,
	}

	go sendMessage(&email)
}

func sendMessage(email *Email) {
	host := os.Getenv("SMTP_HOST")
	port := os.Getenv("SMTP_PORT")
	username := os.Getenv("SMTP_USERNAME")
	password := os.Getenv("SMTP_PASSWORD")

	// Set up auth
	auth := smtp.PlainAuth(
		"",
		username,
		password,
		host)

	// TLS config
	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         host,
	}

	// Dial mail server
	conn, err := tls.Dial("tcp", host+":"+port, tlsconfig)
	if err != nil {
		log.Println(err)
	}

	// Create a new mail client using the connection
	c, err := smtp.NewClient(conn, host)
	if err != nil {
		log.Println(err)
	}

	// Auth
	if err = c.Auth(auth); err != nil {
		log.Println(err)
	}

	// Set "To" address
	if err = c.Mail(email.From.Address); err != nil {
		log.Println(err)
	}

	// Set "From" address
	if err = c.Rcpt(email.To.Address); err != nil {
		log.Println(err)
	}

	// Open data stream for message contents
	w, err := c.Data()
	if err != nil {
		log.Println(err)
	}

	// Create the message
	_, err = w.Write([]byte(email.Email))
	if err != nil {
		log.Println(err)
	}

	// Close the data for the message
	err = w.Close()
	if err != nil {
		log.Println(err)
	}

	// Finalize the send
	c.Quit()

	log.Println("Successfully sent email.")
}
