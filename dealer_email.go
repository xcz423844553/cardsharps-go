package main

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"html/template"
	"net/smtp"
	"strconv"
	"strings"
	"time"
)

type Email struct {
	senderId string
	toIds    []string
	subject  string
	body     string
	password string
}

type SmtpServer struct {
	host string
	port string
}

func (s *SmtpServer) ServerName() string {
	return s.host + ":" + s.port
}

func (email *Email) GetTemplateFileName() string {
	return "template_daily_monitor.html"
}

func (email *Email) BuildMessage() string {
	message := ""
	message += fmt.Sprintf("From: %s\r\n", email.senderId)
	if len(email.toIds) > 0 {
		message += fmt.Sprintf("To: %s\r\n", strings.Join(email.toIds, ";"))
	}
	message += fmt.Sprintf("Subject: %s\r\n", email.subject)
	message += "\r\n"
	message += email.body
	return message
}

func (email *Email) sendEmail() error {
	if email.senderId == "" || len(email.toIds) == 0 {
		return errors.New("Failed to create email. SenderId and ToIds are required.")
	}
	messageBody := email.BuildMessage()
	smtpServer := SmtpServer{host: "smtp.gmail.com", port: "465"}
	auth := smtp.PlainAuth("", email.senderId, email.password, smtpServer.host)
	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         smtpServer.host,
	}

	conn, connErr := tls.Dial("tcp", smtpServer.ServerName(), tlsconfig)
	if connErr != nil {
		return connErr
	}
	defer conn.Close()
	client, smtpErr := smtp.NewClient(conn, smtpServer.host)
	if smtpErr != nil {
		return smtpErr
	}
	defer client.Close()
	if authErr := client.Auth(auth); authErr != nil {
		return authErr
	}
	if fromToErr := client.Mail(email.senderId); fromToErr != nil {
		return fromToErr
	}
	for _, k := range email.toIds {
		if rcptErr := client.Rcpt(k); rcptErr != nil {
			return rcptErr
		}
	}
	writeCloser, emailErr := client.Data()
	if emailErr != nil {
		return emailErr
	}
	_, writeErr := writeCloser.Write([]byte(messageBody))
	if writeErr != nil {
		return writeErr
	}
	defer writeCloser.Close()
	client.Quit()
	return nil
}

func (email *Email) sendEmailTemplate() error {
	if email.senderId == "" || len(email.toIds) == 0 {
		return errors.New("Failed to create email. SenderId and ToIds are required.")
	}
	data := struct {
		Name string
		URL  string
	}{
		Name: "Az",
		URL:  "www.google.com",
	}
	if parseErr := email.ParseTemplate(email.GetTemplateFileName(), data); parseErr != nil {
		return parseErr
	}
	subject := "Subject: " + email.subject + "\n"
	mime := "MIME-version: 1.0;\nContent-Type: text/plain; charset=\"UTF-8\";\n\n"
	messageBody := []byte(subject + mime + "\n" + email.body)
	smtpServer := SmtpServer{host: "smtp.gmail.com", port: "465"}
	auth := smtp.PlainAuth("", email.senderId, email.password, smtpServer.host)
	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         smtpServer.host,
	}

	conn, connErr := tls.Dial("tcp", smtpServer.ServerName(), tlsconfig)
	if connErr != nil {
		return connErr
	}
	defer conn.Close()
	client, smtpErr := smtp.NewClient(conn, smtpServer.host)
	if smtpErr != nil {
		return smtpErr
	}
	defer client.Close()
	if authErr := client.Auth(auth); authErr != nil {
		return authErr
	}
	if fromToErr := client.Mail(email.senderId); fromToErr != nil {
		return fromToErr
	}
	for _, k := range email.toIds {
		if rcptErr := client.Rcpt(k); rcptErr != nil {
			return rcptErr
		}
	}
	writeCloser, emailErr := client.Data()
	if emailErr != nil {
		return emailErr
	}
	_, writeErr := writeCloser.Write(messageBody)
	if writeErr != nil {
		return writeErr
	}
	defer writeCloser.Close()
	client.Quit()
	return nil
}

func (email *Email) SendEmailInTemplate() error {
	if email.senderId == "" || len(email.toIds) == 0 {
		return errors.New("Failed to create email. SenderId and ToIds are required.")
	}
	data := struct {
		Name string
		URL  string
	}{
		Name: "Az",
		URL:  "www.google.com",
	}
	if parseErr := email.ParseTemplate(email.GetTemplateFileName(), data); parseErr != nil {
		return parseErr
	}
	to := "To: " + email.toIds[0] + "\r\n"
	subject := "Subject: " + email.subject + "\r\n"
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	messageBody := []byte(to + subject + mime + "\r\n" + email.body)
	smtpServer := SmtpServer{host: "smtp.gmail.com", port: "587"}
	if err := smtp.SendMail(smtpServer.ServerName(), smtp.PlainAuth("", email.senderId, email.password, smtpServer.host), email.senderId, email.toIds, messageBody); err != nil {
		fmt.Println(err.Error(), "\r\nError Sent email")
		return err
	}
	return nil
}

//SendEmailInTemplateModel2 sends the result of monitor model2 via email in the template
func (email *Email) SendEmailInTemplateModel2(reports []OptionReport) error {
	templateFileName := "template_daily_monitor_model_2.html"
	datum := struct {
		Symbol             string
		OptionType         string
		ExpirationDate     string
		Strike             string
		CurrentPrice       string
		Dates              []string
		DoublePrice        []string
		DoublePriceChange  []string
		CurrentOptionPrice string
		DoubleOptionPrice  string
		Atr                string
		AtrPercent         string
		Ma                 string
		Std                string
		StdPercent         string
	}{}
	data := []interface{}{}
	for _, rpt := range reports {
		datum.Symbol = rpt.Symbol
		datum.OptionType = rpt.OptionType
		datum.ExpirationDate = fmt.Sprintf("%d", rpt.ExpirationDate)
		datum.Strike = fmt.Sprintf("%.2f", rpt.Strike)
		datum.CurrentPrice = fmt.Sprintf("%.2f", rpt.CurrentPrice)
		dt := []string{}
		dp := []string{}
		dpc := []string{}
		for i := range rpt.DoublePrice {
			if i >= 30 {
				break
			}
			t, _ := time.Parse("20060102", strconv.FormatInt(GetTimeInYYYYMMDD64(), 10))
			tStr := t.AddDate(0, 0, i+1).Format("01/02")
			dt = append(dt, fmt.Sprintf("%s", tStr))
			dp = append(dp, fmt.Sprintf("%.2f", rpt.DoublePrice[i]))
			dpc = append(dpc, fmt.Sprintf("%.2f%%", rpt.DoublePriceChange[i]*100))
		}
		datum.Dates = dt
		datum.DoublePrice = dp
		datum.DoublePriceChange = dpc
		datum.CurrentOptionPrice = fmt.Sprintf("%.2f", rpt.CurrentOptionPrice)
		datum.DoubleOptionPrice = fmt.Sprintf("%.2f", rpt.DoubleOptionPrice)
		datum.Atr = fmt.Sprintf("%.2f", rpt.Atr)
		datum.AtrPercent = fmt.Sprintf("%.2f%%", rpt.AtrPercent)
		datum.Ma = fmt.Sprintf("%.2f", rpt.Ma)
		datum.Std = fmt.Sprintf("%.2f", rpt.Std)
		datum.StdPercent = fmt.Sprintf("%.2f%%", rpt.StdPercent)
		data = append(data, datum)
	}
	if email.senderId == "" || len(email.toIds) == 0 {
		return errors.New("Failed to create email. SenderId and ToIds are required.")
	}
	if parseErr := email.ParseTemplate(templateFileName, data); parseErr != nil {
		return parseErr
	}
	to := "To: " + email.toIds[0] + "\r\n"
	subject := "Subject: " + email.subject + "\r\n"
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	messageBody := []byte(to + subject + mime + "\r\n" + email.body)
	smtpServer := SmtpServer{host: "smtp.gmail.com", port: "587"}
	if err := smtp.SendMail(smtpServer.ServerName(), smtp.PlainAuth("", email.senderId, email.password, smtpServer.host), email.senderId, email.toIds, messageBody); err != nil {
		fmt.Println(err.Error(), "\r\nError Sent email")
		return err
	}
	return nil
}

//ParseTemplate parses the data struct into the template
func (email *Email) ParseTemplate(fileName string, data interface{}) error {
	temp, tempErr := template.ParseFiles(fileName)
	if tempErr != nil {
		return tempErr
	}
	buf := new(bytes.Buffer)
	if execErr := temp.Execute(buf, data); execErr != nil {
		return execErr
	}
	email.body = buf.String()
	return nil
}
