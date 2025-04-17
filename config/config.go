package config

import (
	"log/slog"
	"net/http"
	"os"
)

// email data
type Email struct {
	From     string
	Name     string
	Subject  string
	Phone_No string
	Message  string
}

// gpt processer
type Msg struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// gpt options
type Option struct {
	Temperature    float64 `json:"temperature"`
	Seed           int64   `json:"seed"`
	Repeat_penalty float64 `json:"repeat_penalty"`
	Num_ctx        int64   `json:"num_ctx"`
	Num_predict    int64   `json:"num_predict"`
	Top_k          int64   `json:"top_k"`
	Top_p          float64 `json:"top_p"`
}

// gpt prompts
type Prompt struct {
	Model    string `json:"model"`
	Messages []Msg  `json:"messages"`
	Stream   bool   `json:"stream"`
	Raw      bool   `json:"raw"`
	Options  Option `json:"options"`
}

// gpt result
type response struct {
	Message struct {
		Content string `json:"content"`
	} `json:"message"`
}

// config setter
func New() (string, string) {
	logger := slog.New(slog.NewJSONHandler(os.Stderr, nil))
	slog.SetDefault(logger)
	var server_port, admin_port string = os.Getenv("SERVER_PORT"), os.Getenv("ADMIN_PORT")
	if server_port != "" {
		server_port = "11000"
	}
	if admin_port != "" {
		admin_port = "8080"
	}
	return server_port, admin_port
}

func Headers(w http.ResponseWriter) {
	w.Header().Add("Content-Security-Policy", `default-src 'self'; 
		style-src 'self' https://cdn.jsdelivr.net https://fonts.googleapis.com  https://fonts.gstatic.com 'unsafe-inline';
		script-src 'nonce-xyz' 'unsafe-eval';
		font-src 'self' https://fonts.gstatic.com; 
		img-src 'self' https://api.dicebear.com data:; 
		form-action 'self';
		upgrade-insecure-requests;`)
	w.Header().Add("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
	w.Header().Add("X-Frame-Options", "SAMEORIGIN")
	w.Header().Add("X-Content-Type-Options", "nosniff")
	w.Header().Add("Referrer-Policy", "same-origin")
	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	w.Header().Add("Cache-Control", "no-store, no-cache, must-revalidate, proxy-revalidate")
	w.Header().Add("Permissions-Policy", "geolocation=(), microphone=(), camera=()")
}
