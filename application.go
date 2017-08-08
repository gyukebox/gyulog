package main

import (
	"net/http"
	"time"

	"github.com/gyukebox/gyulog/admin"

	"github.com/gyukebox/gyulog/post"
)

func main() {
	defer post.DB.Close()

	mux := http.DefaultServeMux
	files := http.FileServer(http.Dir("./static"))
	server := http.Server{
		Addr:         "172.31.12.72:8080",
		Handler:      mux,
		ReadTimeout:  time.Duration(10 * int64(time.Second)),
		WriteTimeout: time.Duration(600 * int64(time.Second)),
	}

	//add handler for serving static files
	mux.Handle("/static/", http.StripPrefix("/static/", files))

	//match proper handlers
	mux.HandleFunc("/", index)
	mux.HandleFunc("/upload", post.GetPost)
	mux.HandleFunc("/post", postDetail)
	mux.HandleFunc("/admin", adminPage)
	mux.HandleFunc("/authenticate", admin.Authenticate)
	mux.HandleFunc("/logout", admin.Logout)

	server.ListenAndServe()
}
