package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/smtp"
	"os"
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
	var patner_msg []Msg
	patner_msg = append(patner_msg, Msg{
		Role:    "system",
		Content: fmt.Sprint("Your Partners are mentioned below"),
	})
	for _, v := range config_data.About.Patners.Members {
		patner_msg = append(patner_msg, Msg{
			Role: "system",
			Content: fmt.Sprintf("Partner Name: %s\nDetails: %s\nStartup Story: %s\nBackground: %s\nFor more information, please visit the About page. which is accessible through the navigation bar.",
				v.Name, v.Details, v.Startup_story, v.Background),
		})
	}
	var services_msg []Msg
	services_msg = append(services_msg, Msg{
		Role:    "system",
		Content: fmt.Sprint("The Services provided by the firm are mentioned below, they can get started with a service by sending an message from contactus section, which can be accessed from navigation bar, there are no fix prices or price plans for these services, they are decided by the firm."),
	})
	for _, v := range config_data.Services.Options {
		services_msg = append(services_msg, Msg{
			Role:    "system",
			Content: fmt.Sprintf("Service Details: %s\nFor more information, please visit the Services page. You can find the link in the Services section, which is accessible through the navigation bar.", v),
		})
	}
	return Prompt{
		Model: "smollm2",
		Messages: append(
			append(
				append(
					[]Msg{
						{
							Role:    "system",
							Content: "You are an AI chatbot, AuditIQ, designed to play a vital role in chartered accounting by providing accurate and reliable support for various accounting inquiries. Your primary objective is to empower individuals with essential information regarding accounting principles, practices, and services relevant to chartered accountants, ensuring they have the necessary resources to address their financial questions effectively. Users are encouraged to ask pertinent questions related to chartered accounting. Consulting a qualified professional from the firm is always the best approach for specialized matters or uncertainties. Remember, the user directs the questions, so your responses must and should be formatted in plain text, and you can not use markdown. Give short but accurate responses, if incase you don't have an answer, tell them to contact the firm for more details, prices are decided by the firm and should contact them.",
						},
					}, patner_msg...,
				), services_msg...,
			), msg...),
		Stream: false,
		Options: Option{
			Temperature:    0.7,
			Repeat_penalty: 0.9,
			Num_ctx:        2050,
			Seed:           42,
			Num_predict:    75,
			Top_k:          40,
			Top_p:          0.85,
		},
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
		return fmt.Sprint(config_data.Contacts.Responce, "\n\nHere's what we have recived:\n\n", data)
	}

	// will send response to us
	return fmt.Sprint("You have a new email from:\n", data)
}

// sends a response to the specified user's email
func (email Email) send(to bool) {
	// requirements
	auth_email := os.Getenv("ADMIN_EMAIL_ID")
	password := os.Getenv("ADMIN_EMAIL_PASSWORD")
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	// default user values
	to_email := []string{"support@ca-gk.org"}
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
