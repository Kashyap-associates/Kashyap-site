package functions

import (
	"html/template"
	"net/http"
)

var errPath = "error/"

func Error(w http.ResponseWriter, r *http.Request) {
	var files []string
	for _, i := range []string{
		layout, errPath + "error.html",
	} {
		files = append(files, path+i)
	}
	temp, err := template.ParseFS(
		pages, files...,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
  w.WriteHeader(http.StatusInternalServerError)
	if err := temp.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func NotFound(w http.ResponseWriter, r *http.Request) {
	var files []string
	for _, i := range []string{
		layout, errPath + "404.html",
	} {
		files = append(files, path+i)
	}
	temp, err := template.ParseFS(
		pages, files...,
	)
	if err != nil {
    http.Redirect(w, r, "/error", http.StatusSeeOther)
		return
	}
	if err := temp.Execute(w, nil); err != nil {
    http.Redirect(w, r, "/error", http.StatusSeeOther)
		return
	}
}
