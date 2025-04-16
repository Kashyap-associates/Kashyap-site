package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/smtp"
)

func (prompt Prompt) TalkToAI() string {
	// creating an error variable
	var err error

	// converting the prompt to json type
	var jsonBody []byte
	if jsonBody, err = json.Marshal(prompt); err != nil {
		slog.Error(err.Error())
		return ""
	}

	// creating a request
	var req *http.Request
	if req, err = http.NewRequest(
		"POST", "http://ollama-ai:11434/api/chat", bytes.NewBuffer(jsonBody),
	); err != nil {
		slog.Error(err.Error())
		return ""
	}

	// sending the request & getting response
	req.Header.Set("Content-Type", "application/json")
	var res *http.Response
	var client http.Client
	if res, err = client.Do(req); err != nil {
		slog.Error(err.Error())
		return ""
	}

	// parsing the response
	defer res.Body.Close()
	var body []byte
	if body, err = io.ReadAll(res.Body); err != nil {
		slog.Error(err.Error())
		return ""
	}

	// converting the json data to response type
	var result response
	if err = json.Unmarshal(body, &result); err != nil {
		// trying to parse the data (to see if it's an error from ai or not)
		slog.Error(err.Error())
		slog.Info("trying again...")
		var errMsg map[string]interface{}
		if err = json.Unmarshal(body, &errMsg); err != nil {
			slog.Error(err.Error())
			return ""
		}
		return "ERROR"
	}

	// returning the response message
	return result.Message.Content
}

func CreatePrompt(msg []Msg) Prompt {
	// will create & return a prompt type object
	return Prompt{
		Model: "smollm2",
		Messages: append([]Msg{
			{
				Role:    "system",
				Content: "You are an AI chatbot designed to play a vital role in chartered accounting by providing accurate and reliable support for various accounting inquiries. Your primary objective is to empower individuals with essential information regarding accounting principles, practices, and services relevant to chartered accountants, ensuring they have the necessary resources to address their financial questions effectively. Users are encouraged to ask pertinent questions related to chartered accounting. Consulting a qualified professional from the firm is always the best approach for specialized matters or uncertainties. Remember, the user directs the questions, so your responses must and should be formatted in Plain text and can not use markdown.",
			},
		}, msg...),
		Stream: false,
	}
}

// this will return the response format for the email
func (email Email) render(to bool) string {
	// main data format
	data := fmt.Sprint(
		"\nName:\n", email.Name, "\n",
		"\nEmail:\n", email.From, "\n",
		"\nPhone No:\n", email.Phone_No, "\n",
		"\nSubject:\n", email.Subject, "\n",
		"\nMessage:\n", email.Message, "\n",
	)

	// will send response to user
	if to {
		return fmt.Sprint("Thank you for choosing Gopi & Kashyap CA firm, we will get back to you shortly!\n\nHere's what we have recived:\n\n", data)
	}

	// will send response to us
	return fmt.Sprint("You have a new email from:\n", data)
}

// sends a response to the specified user's email
func (email Email) send(to bool) {
	// requirements
	auth_email := "yasasvigumma@gmail.com"
	password := "xxsz dcec vojc gzod"
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	// default user values
	to_email := []string{auth_email}
	message := []byte(fmt.Sprint("Subject: ", "New email from ", email.Name, "\n\n", email.render(false), "\n"))

	// option to change the values
	if to {
		to_email = []string{email.From}
		message = []byte(fmt.Sprint("Subject: ", "Email from Gopi & Kashyap CA firm", "\n\n", email.render(true), "\n"))
	}

	// sending the response
	auth := smtp.PlainAuth("", auth_email, password, smtpHost)
	if err := smtp.SendMail(fmt.Sprint(smtpHost, ":", smtpPort), auth, auth_email, to_email, message); err != nil {
		slog.Error(err.Error())
	}
}

func Send_email(email Email) {
	// first email
	email.send(false)

	// second email
	email.send(true)
}
