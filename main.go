package main

import (
	"net/http"
	"time"
)

func main() {

	mux := http.DefaultServeMux
	files := http.FileServer(http.Dir(settings.Static))
	server := http.Server{
		Addr:         settings.Address,
		Handler:      mux,
		ReadTimeout:  time.Duration(settings.ReadTimeout * int64(time.Second)),
		WriteTimeout: time.Duration(settings.WriteTimeout * int64(time.Second)),
	}

	//add handler for serving static files
	mux.Handle("/static/", http.StripPrefix("/static/", files))

	//match proper handlers
	mux.HandleFunc("/", index)

	server.ListenAndServe()
}
