package functions

import (
	"embed"
	"html/template"
	"net/http"
)

//go:embed pages/*
var pages embed.FS

const path = "pages/"

func Index(w http.ResponseWriter, r *http.Request) {
	var files []string
	for _, i := range []string{
		"layout.html",
		"index.html",
		"home.html",
		"about.html",
		"services.html",
		"contact.html",
		"other.html",
	}{
		files = append(files, path+i)
	}
	temp, err := template.ParseFS(
		pages, files...,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := temp.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
