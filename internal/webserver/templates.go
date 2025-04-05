package webserver

import (
	"fmt"
	"log"
	"net/http"
	"text/template"
)

type Templates struct {
	templates *template.Template
}

func newTemplate(template_glob_path string) *Templates {
	return &Templates{
		templates: template.Must(template.ParseGlob(template_glob_path)),
	}
}

func (t *Templates) Render(w http.ResponseWriter, name string, status int, data interface{}) error {
	w.Header().Add("Content-Type", "text/html")
	w.WriteHeader(status)
	fileName := fmt.Sprintf("%s.html", name)
	err := t.templates.ExecuteTemplate(w, fileName, data)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
