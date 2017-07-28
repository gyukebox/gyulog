package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/gyukebox/gyulog/post"
)

// temporary
func index(w http.ResponseWriter, r *http.Request) {
	var index int
	end, err := url.QueryUnescape(r.URL.RawQuery)
	if err != nil {
		index = 0
	} else {
		index, _ = strconv.Atoi(end)
	}
	if err != nil {
		fmt.Print("At Handler, ")
		log.Fatalln(err)
	}
	posts, err := post.GetFivePosts(index)
	data := map[string]interface{}{
		"Post":     posts,
		"Previous": index - 1,
		"Next":     index + 1,
	}
	generateHTML(w, data, "index", "layout", "mobile", "navbar", "sidebar")
}

// main page 에서 필요한 데이터
// 1. post list
// 2. index, 앞뒤 index

// post page 에서 필요한 데이터
// 1. single post info

func postDetail(w http.ResponseWriter, r *http.Request) {
	var id int
	query, err := url.QueryUnescape(r.URL.RawQuery)
	if err != nil {
		fmt.Print("At parsing url, ")
		log.Fatalln(err)
	} else {
		id, err = strconv.Atoi(query)
		if err != nil {
			fmt.Print("At dealing with query, ")
			log.Fatalln(err)
		}
	}
	//post := post.GetPostByTitle(title)
	post := post.GetPostById(id)
	data := map[string]interface{}{
		"Post": map[string]interface{}{
			"Title":         post.Title,
			"Body":          template.HTML(post.Body),
			"PublishedDate": post.PublishedDate,
			"Summary":       template.HTML(post.Summary),
			"Id":            post.Id,
			"Previous":      post.Id - 1,
			"Next":          post.Id + 1,
		},
	}
	generateHTML(w, data, "post", "layout", "mobile", "navbar", "sidebar")
}
