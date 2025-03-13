package functions

import (
	"html/template"
	"net/http"
)

func Admin(w http.ResponseWriter, r *http.Request) {
	var files []string
	for _, i := range []string{
		layout, "admin.html",
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
