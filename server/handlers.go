package server

import (
	"Kashyap-site/config"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"
)

var (
	content *config.DataType
	common  commonType
)

func init() {
	content = config.Read()
	common = commonType{
		Top: updateAnnotation(),
		Links: struct {
			Email    string
			Telegram string
			Linkedin string
		}{
			Email:    content.Links.Email,
			Telegram: content.Links.Telegram,
			Linkedin: content.Links.Linkedin,
		},
	}
}

func updateAnnotation() string {
	var annotation string
	if len(content.Annotations) > 0 {
		annotation = content.Annotations[0].Short
	}
	return annotation
}

func createToken() string {
	data := make([]byte, 16)
	if _, err := rand.Read(data); err != nil {
		slog.Error(err.Error())
		return ""
	}
	return base64.RawStdEncoding.EncodeToString(data)
}

// [main] / (home)
func home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Redirect(w, r, "/404", http.StatusFound)
		return
	}
	config.Headers(w)
	var people []struct {
		Name    string
		Image   string
		Details string
		Links   struct {
			Linkedin string
			Whatsapp string
			Email    string
		}
	}
	if len(content.About.Patners.Members) > 0 {
		for _, v := range content.About.Patners.Members {
			people = append(people, struct {
				Name    string
				Image   string
				Details string
				Links   struct {
					Linkedin string
					Whatsapp string
					Email    string
				}
			}{
				Name:    v.Name,
				Image:   v.Image,
				Details: v.Details,
				Links: struct {
					Linkedin string
					Whatsapp string
					Email    string
				}{
					Linkedin: v.Links.Linkedin,
					Whatsapp: v.Links.Whatsapp,
					Email:    v.Links.Email,
				},
			})
		}
	}
	var services []struct {
		Image   string
		Name    string
		Details string
	}
	if len(content.Services.Data) > 0 {
		for i, v := range content.Services.Data {
			if i >= 6 {
				break
			}
			services = append(services, struct {
				Image   string
				Name    string
				Details string
			}{
				Image:   v.Image,
				Name:    v.Name,
				Details: v.Details,
			})
		}
	}
	common.Top = updateAnnotation()
	var err error
	var homeData, commonData []byte
	if homeData, err = json.Marshal(homeType{
		Hero: struct {
			Title     string
			Sub_title string
		}{
			Title:     content.Hero.Title,
			Sub_title: content.Hero.Sub_title,
		},
		About: struct {
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
		}{
			Image:   content.About.Image,
			Details: content.About.Intro,
			People:  people,
		},
		Services: services,
		Others: struct {
			Details string
			Tasks   []struct {
				Name    string
				Title   string
				Details string
			}
		}(content.Others),
		Contacts: struct {
			Addr     string
			Phone    string
			Email    string
			Responce string
		}{
			Addr:     content.Contacts.Addr,
			Phone:    content.Contacts.Phone,
			Email:    content.Contacts.Email,
			Responce: content.Contacts.Responce,
		},
	}); err != nil {
		slog.Error(err.Error())
	}
	if commonData, err = json.Marshal(common); err != nil {
		slog.Error(err.Error())
	}
	if err := indexTmpl.Execute(w, sendType{
		Data:   string(homeData),
		Common: string(commonData),
	}); err != nil {
		handleError(w, r, err, false)
	}
}

// [main] /about
func about(w http.ResponseWriter, r *http.Request) {
	config.Headers(w)
	var patners []struct {
		Name string
		Img  string
	}
	if len(content.About.Patners.Members) > 0 {
		for _, v := range content.About.Patners.Members {
			patners = append(patners, struct {
				Name string
				Img  string
			}{
				Name: v.Name,
				Img:  v.Image,
			})
		}
	}
	var team []struct {
		Img      string
		Name     string
		Position string
		Email    string
	}
	if len(content.About.Team.Members) > 0 {
		for _, v := range content.About.Team.Members {
			team = append(team, struct {
				Img      string
				Name     string
				Position string
				Email    string
			}{
				Img:      v.Image,
				Name:     v.Name,
				Position: v.Position,
				Email:    v.Link,
			})
		}
	}
	common.Top = updateAnnotation()
	var err error
	var aboutData, commonData []byte
	if aboutData, err = json.Marshal(aboutType{
		IsPerson:  false,
		Story:     content.About.Story,
		Intro:     content.About.Intro,
		Choose_us: content.About.Choose_us,
		Contact: struct {
			Intro string
			Email string
		}{
			Intro: content.Contacts.Intro,
			Email: content.Contacts.Email,
		},
		Services: struct {
			Intro   string
			Options []string
		}{
			Intro:   content.Services.Short_intro,
			Options: content.Services.Options,
		},
		Patners: struct {
			Intro   string
			Members []struct {
				Name string
				Img  string
			}
		}{
			Intro:   content.About.Patners.Intro,
			Members: patners,
		},
		Team: struct {
			Intro   string
			Members []struct {
				Img      string
				Name     string
				Position string
				Email    string
			}
		}{
			Intro:   content.About.Team.Intro,
			Members: team,
		},
	}); err != nil {
		slog.Error(err.Error())
	}
	if commonData, err = json.Marshal(common); err != nil {
		slog.Error(err.Error())
	}
	if err := aboutTmpl.Execute(w, sendType{
		Data:   string(aboutData),
		Common: string(commonData),
	}); err != nil {
		handleError(w, r, err, false)
	}
}

// [main] /about/:id
func about_with_id(w http.ResponseWriter, r *http.Request) {
	type _patner struct {
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
	config.Headers(w)
	user := r.PathValue("id")
	common.Top = updateAnnotation()
	partner := make(map[string]_patner)
	if len(content.About.Patners.Members) > 0 {
		for _, v := range content.About.Patners.Members {
			partner[strings.ToLower(v.Name)] = _patner{
				Name:       v.Name,
				Intro:      v.Details,
				Startup:    v.Startup_story,
				Background: v.Background,
				Feedback: struct {
					Intro string
					Quote struct {
						Name string
						Said string
					}
				}(v.Feedback),
				Links: struct {
					Img      string
					Email    string
					Whatsapp string
					Linkedin string
				}{
					Img:      v.Image,
					Email:    v.Links.Email,
					Whatsapp: v.Links.Whatsapp,
					Linkedin: v.Links.Linkedin,
				},
			}
		}
	}
	var ok bool
	var info _patner
	if info, ok = partner[user]; !ok {
		http.Redirect(w, r, "/404", http.StatusFound)
		return
	}
	var err error
	var aboutData, commonData []byte
	if aboutData, err = json.Marshal(struct {
		IsPerson bool
		Person   _patner
	}{
		IsPerson: true,
		Person:   info,
	}); err != nil {
		slog.Error(err.Error())
	}
	if commonData, err = json.Marshal(common); err != nil {
		slog.Error(err.Error())
	}
	if err := aboutTmpl.Execute(w, sendType{
		Data:   string(aboutData),
		Common: string(commonData),
	}); err != nil {
		handleError(w, r, err, false)
	}
}

// [main] /services
func services(w http.ResponseWriter, r *http.Request) {
	config.Headers(w)
	common.Top = updateAnnotation()
	var data []struct {
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
	if len(content.Services.Data) > 0 {
		for _, v := range content.Services.Data {
			data = append(data, struct {
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
			}(v))
		}
	}
	var err error
	var servicesData, commonData []byte
	if servicesData, err = json.Marshal(servicesType{
		Intro: content.Services.Intro,
		Data:  data,
	}); err != nil {
		slog.Error(err.Error())
	}
	if commonData, err = json.Marshal(common); err != nil {
		slog.Error(err.Error())
	}
	if err := servicesTmpl.Execute(w, sendType{
		Data:   string(servicesData),
		Common: string(commonData),
	}); err != nil {
		handleError(w, r, err, false)
	}
}

// [main] /calculator
func calculator(w http.ResponseWriter, r *http.Request) {
	data := "/404"
	if content.Calculator != "" {
		data = content.Calculator
	}
	http.Redirect(w, r, data, http.StatusFound)
}

// [main] /annotation
func annotation(w http.ResponseWriter, r *http.Request) {
	config.Headers(w)
	var msg []struct {
		Short string
		Long  string
	}
	if len(content.Annotations) > 0 {
		for _, v := range content.Annotations {
			msg = append(msg, struct {
				Short string
				Long  string
			}{
				Short: v.Short,
				Long:  v.Long,
			})
		}
	}
	common.Top = updateAnnotation()
	var err error
	var annotaData, commonData []byte
	if annotaData, err = json.Marshal(otherType{
		Thanks: false,
		Msg:    msg,
	}); err != nil {
		slog.Error(err.Error())
	}
	if commonData, err = json.Marshal(common); err != nil {
		slog.Error(err.Error())
	}
	if err := otherTmpl.Execute(w, sendType{
		Data:   string(annotaData),
		Common: string(commonData),
	}); err != nil {
		handleError(w, r, err, false)
	}
}

// [main] /thanks
func thanks(w http.ResponseWriter, r *http.Request) {
	config.Headers(w)
	common.Top = updateAnnotation()
	var err error
	var thanksData, commonData []byte
	if thanksData, err = json.Marshal(otherType{
		Thanks: true,
	}); err != nil {
		slog.Error(err.Error())
	}
	if commonData, err = json.Marshal(common); err != nil {
		slog.Error(err.Error())
	}
	if err := otherTmpl.Execute(w, sendType{
		Data:   string(thanksData),
		Common: string(commonData),
	}); err != nil {
		handleError(w, r, err, false)
	}
}

// [admin] /dashboard
func admin(w http.ResponseWriter, r *http.Request) {
	var err error
	if _, err = r.Cookie("auth_token"); err != nil {
		http.Redirect(w, r, "/denied", http.StatusFound)
	}
	config.Headers(w)
	var adminData, commonData []byte
	if adminData, err = json.Marshal(adminType{
		Login: false,
		Main:  content,
	}); err != nil {
		slog.Error(err.Error())
	}
	if commonData, err = json.Marshal(common); err != nil {
		slog.Error(err.Error())
	}
	if err := adminTmpl.Execute(w, sendType{
		Data:   string(adminData),
		Common: string(commonData),
	}); err != nil {
		handleError(w, r, err, false)
	}
}

// [admin] / (login)
func admin_login(w http.ResponseWriter, r *http.Request) {
	var err error
	if r.URL.Path != "/" {
		http.Redirect(w, r, "/404", http.StatusFound)
	}
	config.Headers(w)
	var adminData, commonData []byte
	if adminData, err = json.Marshal(adminType{
		Login: true,
	}); err != nil {
		slog.Error(err.Error())
	}
	if commonData, err = json.Marshal(common); err != nil {
		slog.Error(err.Error())
	}
	if err := adminTmpl.Execute(w, sendType{
		Data:   string(adminData),
		Common: string(commonData),
	}); err != nil {
		handleError(w, r, err, false)
	}
}

// [admin] POST /login
func auth_login(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	if (username != os.Getenv("ADMIN_USERNAME")) || (password != os.Getenv("ADMIN_PASSWORD")) {
		http.Redirect(w, r, "/denied", http.StatusFound)
	} else {
		http.SetCookie(w, &http.Cookie{
			Name:     "auth_token",
			Value:    createToken(),
			Path:     "/",
			Expires:  time.Now().Add(1 * time.Hour),
			MaxAge:   3600,
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteLaxMode,
		})
		http.Redirect(w, r, "/dashboard", http.StatusFound)
	}
}

// [admin] POST /update
func update_data(w http.ResponseWriter, r *http.Request) {
	var err error
	if _, err = r.Cookie("auth_token"); err != nil {
		http.Redirect(w, r, "/denied", http.StatusFound)
	}
	var bodyByte []byte
	if bodyByte, err = io.ReadAll(r.Body); err != nil {
		slog.Error(err.Error())
		return
	}
	config.Set(string(bodyByte))
}

// [admin] POST /login
func auth_logout(w http.ResponseWriter, r *http.Request) {
	var err error
	var cookie *http.Cookie
	if cookie, err = r.Cookie("auth_token"); err != nil {
		slog.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	cookie.Expires = time.Unix(0, 0)
	cookie.MaxAge = -1
	http.SetCookie(w, cookie)
}

// [admin] /denied
func access_denied(w http.ResponseWriter, r *http.Request) {
	config.Headers(w)
	var err error
	var errData, commonData []byte
	if errData, err = json.Marshal(errorType{
		ErrorTitle: fmt.Sprint(http.StatusForbidden),
		ErrorMsg:   "Access Denied!",
	}); err != nil {
		slog.Error(err.Error())
	}
	if commonData, err = json.Marshal(common); err != nil {
		slog.Error(err.Error())
	}
	if err := errorTmpl.Execute(w, sendType{
		Data:   string(errData),
		Common: string(commonData),
	}); err != nil {
		slog.Error(err.Error())
		handleError(w, r, err, true)
	}
}

// [main] POST /chat
func chat(w http.ResponseWriter, r *http.Request) {
	config.Headers(w)
	var data []config.Msg
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	prompt := config.CreatePrompt(data)
	msg := prompt.TalkToAI()
	if msg == "" || msg == "ERROR" {
		if msg == "ERROR" {
			http.Error(w, "No responce", http.StatusNotFound)
			return
		}
		http.Error(w, "Something went wrong!", http.StatusInternalServerError)
		return
	}
	res := map[string]string{"msg": msg}
	if err := json.NewEncoder(w).Encode(res); err != nil {
		handleError(w, r, err, true)
		return
	}
}

// [main] POST /email
func email(w http.ResponseWriter, r *http.Request) {
	config.Headers(w)
	config.Send_email(config.Email{
		From:     r.FormValue("email"),
		Name:     r.FormValue("name"),
		Subject:  r.FormValue("subject"),
		Phone_No: r.FormValue("phone_no"),
		Message:  r.FormValue("message"),
	})
	http.Redirect(w, r, "/thanks", http.StatusFound)
}

// [main] /404
func page_not_found(w http.ResponseWriter, r *http.Request) {
	config.Headers(w)
	var err error
	var errData, commonData []byte
	if errData, err = json.Marshal(errorType{
		ErrorTitle: fmt.Sprint(http.StatusNotFound),
		ErrorMsg:   "Page Not Found!",
	}); err != nil {
		slog.Error(err.Error())
	}
	if commonData, err = json.Marshal(common); err != nil {
		slog.Error(err.Error())
	}
	if err := errorTmpl.Execute(w, sendType{
		Data:   string(errData),
		Common: string(commonData),
	}); err != nil {
		handleError(w, r, err, false)
	}
}

// [main] /error
func error_page(w http.ResponseWriter, r *http.Request) {
	config.Headers(w)
	var err error
	var errData, commonData []byte
	if errData, err = json.Marshal(errorType{
		ErrorTitle: "OOPS!",
		ErrorMsg:   "Something went wrong.",
	}); err != nil {
		slog.Error(err.Error())
	}
	if commonData, err = json.Marshal(common); err != nil {
		slog.Error(err.Error())
	}
	if err := errorTmpl.Execute(w, sendType{
		Data:   string(errData),
		Common: string(commonData),
	}); err != nil {
		slog.Error(err.Error())
		handleError(w, r, err, true)
	}
}
