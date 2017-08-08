package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
)

type Setting struct {
	Address      string
	ReadTimeout  int64
	WriteTimeout int64
	Static       string
}

var settings Setting
var logger *log.Logger

func init() {
	file, err := os.OpenFile("gyulog.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalln(err)
	}
	logger = log.New(file, "INFO ", log.Ldate|log.Ltime|log.Lshortfile)
}

func generateHTML(w http.ResponseWriter, data interface{}, filenames ...string) {
	var files []string
	for _, fname := range filenames {
		files = append(files, fmt.Sprintf("templates/%s.html", fname))
	}
	templates := template.Must(template.ParseFiles(files...))
	templates.ExecuteTemplate(w, "layout", data)
}
