package main

import (
	"encoding/json"
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
	openSetting()
	file, err := os.OpenFile("gyulog.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalln(err)
	}
	logger = log.New(file, "INFO ", log.Ldate|log.Ltime|log.Lshortfile)
}

func openSetting() {
	file, err := os.Open("settings.json")
	if err != nil {
		fmt.Println("Failed to open file")
		log.Fatal(err)
	}
	decoder := json.NewDecoder(file)
	settings = Setting{}
	err = decoder.Decode(&settings)
	if err != nil {
		fmt.Println("Failed to parse JSON")
		log.Fatal(err)
	}
}

func generateHTML(w http.ResponseWriter, data interface{}, filenames ...string) {
	var files []string
	for _, fname := range filenames {
		files = append(files, fmt.Sprintf("templates/%s.html", fname))
	}
	templates := template.Must(template.ParseFiles(files...))
	templates.ExecuteTemplate(w, "layout", data)
}
