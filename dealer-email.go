package main

// import (
// 	"crypto/tls"
// 	"errors"
// 	"fmt"
// 	"net/smtp"
// 	"strings"
// )

// type Email struct {
// 	senderId string
// 	toIds    []string
// 	subject  string
// 	body     string
// 	password string
// }

// type SmtpServer struct {
// 	host string
// 	port string
// }

// func (s *SmtpServer) ServerName() string {
// 	return s.host + ":" + s.port
// }

// func (email *Email) BuildMessage() string {
// 	message := ""
// 	message += fmt.Sprintf("From: %s\r\n", email.senderId)
// 	if len(email.toIds) > 0 {
// 		message += fmt.Sprintf("To: %s\r\n", strings.Join(email.toIds, ";"))
// 	}
// 	message += fmt.Sprintf("Subject: %s\r\n", email.subject)
// 	message += "\r\n"
// 	message += email.body
// 	return message
// }

// func (email *Email) sendEmail() error {
// 	if email.senderId == "" || len(email.toIds) == 0 {
// 		return errors.New("Failed to create email. SenderId and ToIds are required.")
// 	}
// 	messageBody := email.BuildMessage()
// 	smtpServer := SmtpServer{host: "smtp.gmail.com", port: "465"}
// 	auth := smtp.PlainAuth("", email.senderId, email.password, smtpServer.host)
// 	tlsconfig := &tls.Config{
// 		InsecureSkipVerify: true,
// 		ServerName:         smtpServer.host,
// 	}

// 	conn, connErr := tls.Dial("tcp", smtpServer.ServerName(), tlsconfig)
// 	if connErr != nil {
// 		return connErr
// 	}
// 	defer conn.Close()
// 	client, smtpErr := smtp.NewClient(conn, smtpServer.host)
// 	if smtpErr != nil {
// 		return smtpErr
// 	}
// 	defer client.Close()
// 	if authErr := client.Auth(auth); authErr != nil {
// 		return authErr
// 	}
// 	if fromToErr := client.Mail(email.senderId); fromToErr != nil {
// 		return fromToErr
// 	}
// 	for _, k := range email.toIds {
// 		if rcptErr := client.Rcpt(k); rcptErr != nil {
// 			return rcptErr
// 		}
// 	}
// 	writeCloser, emailErr := client.Data()
// 	if emailErr != nil {
// 		return emailErr
// 	}
// 	_, writeErr := writeCloser.Write([]byte(messageBody))
// 	if writeErr != nil {
// 		return writeErr
// 	}
// 	defer writeCloser.Close()
// 	client.Quit()
// 	return nil
// }
