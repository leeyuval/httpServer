package utils

import (
	"html/template"
	"net/http"
	"time"
)

func FormatCreationTime(timeStr string) string {
	creationTime, _ := time.Parse(time.RFC3339, timeStr)
	return creationTime.Format("2006-01-02 15:04")
}

func RenderHTMLTemplate(w http.ResponseWriter, tmplPath string, data interface{}) error {
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "text/html")
	err = tmpl.Execute(w, data)
	if err != nil {
		return err
	}

	return nil
}
