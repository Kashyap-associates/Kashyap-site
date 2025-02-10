package functions

import (
	"embed"
	"html/template"
	"net/http"
)

//go:embed pages/*
var pages embed.FS

const (
	path     = "pages/"
	mainPath = "main/"
	layout   = "layout.html"
)

func Index(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Redirect(w, r, "/404", http.StatusSeeOther)
		return
	}
	var files []string
	for _, i := range []string{
		layout,
		mainPath + "index.html",
		mainPath + "home.html",
		mainPath + "about.html",
		mainPath + "services.html",
		mainPath + "contact.html",
		mainPath + "other.html",
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
