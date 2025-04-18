package server

import (
	"Kashyap-site/config"
	"embed"
	"html/template"
	"log/slog"
	"net/http"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

// common server type
type _server map[string]http.HandlerFunc

// error page type
type errorType struct {
	Common     commonType
	ErrorTitle string
	ErrorMsg   string
}

// send type is used to send data
type sendType struct {
	Data   string
	Common string
}

// layout page type
type commonType struct {
	Top   string
	Links struct {
		Email    string
		Telegram string
		Linkedin string
	}
}

// other page type
type otherType struct {
	Msg []struct {
		Short string
		Long  string
	}
}

// home page type
type homeType struct {
	Hero struct {
		Title     string
		Sub_title string
	}
	About struct {
		Image   string
		Details string
		People  []struct {
			Name    string
			Image   string
			Details string
			Links   struct {
				Linkedin string
				Whatsapp string
				Email    string
			}
		}
	}
	Services []struct {
		Image   string
		Name    string
		Details string
	}
	Others struct {
		Details string
		Tasks   []struct {
			Name    string
			Title   string
			Details string
		}
	}
	Contacts struct {
		Addr     string
		Phone    string
		Email    string
		Responce string
	}
}

// about page type
type aboutType struct {
	IsPerson  bool
	Story     string
	Intro     string
	Choose_us string
	Contact   struct {
		Intro string
		Email string
	}
	Services struct {
		Intro   string
		Options []string
	}
	Patners struct {
		Intro   string
		Members []struct {
			Name string
			Img  string
		}
	}
	Team struct {
		Intro   string
		Members []struct {
			Img      string
			Name     string
			Position string
			Email    string
		}
	}
	Person struct {
		Name       string
		Intro      string
		Startup    string
		Background string
		Feedback   struct {
			Intro string
			Quote struct {
				Name string
				Said string
			}
		}
		Links struct {
			Img      string
			Email    string
			Whatsapp string
			Linkedin string
		}
	}
}

// services page type
type servicesType struct {
	Intro string
	Data  []struct {
		Image        string
		Name         string
		Time         string
		Deliverable  []string
		Regulation   string
		Catagory     string
		Availability string
		Audience     string
		Details      string
		Scope        string
	}
}

// admin server type
type adminType struct {
	Login    bool
	Email    string
	Password string
	Main     interface{}
}

var (
	//go:embed pages/*
	files embed.FS

	//go:embed assets/*
	assets embed.FS

	// will couple middlewares
	middleware Middleware

	// template variables
	indexTmpl    *template.Template
	adminTmpl    *template.Template
	otherTmpl    *template.Template
	errorTmpl    *template.Template
	aboutTmpl    *template.Template
	servicesTmpl *template.Template
)

const (
	// pages dir
	pages = "pages/"
	// common layout page
	layout = pages + "layout.html"
)

// will fetch the templates
func get_template(file string) *template.Template {
	tmpl, err := template.ParseFS(files, layout, pages+file)
	if err != nil {
		slog.Error(err.Error())
		return nil
	}
	return tmpl
}

// function will run on start
func init() {
	indexTmpl = get_template("home.html")
	otherTmpl = get_template("other.html")
	adminTmpl = get_template("admin.html")
	errorTmpl = get_template("error.html")
	aboutTmpl = get_template("about.html")
	servicesTmpl = get_template("services.html")

	middleware = createStack(
		// loggingMiddleware,
		gzipMiddleware,
	)
}

// _server type method will serve the handler
func (routes _server) serve() http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/assets/", http.StripPrefix("/",
		http.FileServerFS(assets),
	))
	mux.HandleFunc("/robots.txt", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/assets/robots.txt", http.StatusFound)
	})
	for route, handler := range routes {
		mux.HandleFunc(route, handler)
	}
	return middleware(mux)
}

// main server
func New() http.Handler {
	return _server{
		"/":            home,
		"/thanks":      thanks,
		"/about/":      about,
		"/services":    services,
		"/annotations": annotation,
		"/calculator":  calculator,
		"/about/{id}":  about_with_id,
		"/404":         page_not_found,
		"/error":       error_page,
		"POST /chat":   chat,
		"POST /email":  email,
	}.serve()
}

// admin server
func NewAdmin() http.Handler {
	return _server{
		"/":            admin_login,
		"/dashboard":   admin,
		"/404":         page_not_found,
		"/error":       error_page,
		"/denied":      access_denied,
		"POST /login":  auth_login,
		"POST /update": update_data,
		"POST /logout": auth_logout,
	}.serve()
}

func Telegram(token string) {
	var err error
	if token == "" {
		slog.Error("Token not found!")
		return
	}
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	var bot *tgbotapi.BotAPI
	if bot, err = tgbotapi.NewBotAPI(token); err != nil {
		slog.Error(err.Error())
		return
	}
	var updates tgbotapi.UpdatesChannel
	if updates, err = bot.GetUpdatesChan(u); err != nil {
		slog.Error(err.Error())
		return
	}
	slog.Info("Telegram Server started!")
	for update := range updates {
		if update.Message == nil {
			continue
		}
		var data string
		if update.Message.Text != "/start" {
			arrMsg := []config.Msg{
				{
					Role:    "user",
					Content: update.Message.Text,
				},
			}
			prompt := config.CreatePrompt(arrMsg)
			data = prompt.TalkToAI()
		} else {
			data = "Hi There, I'm AuditIQ\nHow can I help you today ?"
		}
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, data)
		_, err := bot.Send(msg)
		if err != nil {
			slog.Error(err.Error())
			return
		}
	}
}
