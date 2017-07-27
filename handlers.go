package main

import (
	"fmt"
	"net/http"

	"github.com/gyukebox/gyulog/post"
)

// temporary
func index(w http.ResponseWriter, r *http.Request) {
	posts, err := post.GetAllPosts()
	if err != nil {
		fmt.Print("At Handler : ")
		fmt.Println(err)
	}
	//test
	fmt.Println(posts)
	generateHTML(w, posts, "index", "layout", "mobile", "navbar", "pager", "sidebar")
}
