package utils

import (
	"html/template"
	"net/http"
	"time"
)

type Renderer interface {
	GetTemplateName() string
}

func FormatCreationTime(timeStr string) string {
	creationTime, _ := time.Parse(time.RFC3339, timeStr)
	return creationTime.Format("2006-01-02 15:04")
}

func RenderHTML(w http.ResponseWriter, data Renderer) {
	tmpl, err := template.ParseFiles("templates/" + data.GetTemplateName())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
