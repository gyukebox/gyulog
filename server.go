package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/gyukebox/gyulog/admin"
	"github.com/gyukebox/gyulog/post"
)

func index(w http.ResponseWriter, r *http.Request) {
	var index int
	end, err := url.QueryUnescape(r.URL.RawQuery)
	if err != nil {
		index = 0
	} else {
		index, _ = strconv.Atoi(end)
	}
	posts, err := post.GetFivePosts(index)
	if err != nil {
		fmt.Print("At Handler, ")
		log.Fatalln(err)
	}
	data := map[string]interface{}{
		"Post":     posts,
		"Previous": index - 1,
		"Next":     index + 1,
		"Limit":    (post.TotalPosts - 1) / 5,
	}
	generateHTML(w, data, "index", "layout", "navbar")
}

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

	result, _ := post.GetPostById(id)

	data := map[string]interface{}{
		"Post": map[string]interface{}{
			"Title":         result.Title,
			"Body":          template.HTML(result.Body),
			"PublishedDate": result.PublishedDate,
			"Summary":       template.HTML(result.Summary),
			"Id":            result.Id,
		},
		"Previous": result.Id - 1,
		"Next":     result.Id + 1,
		"Total":    post.TotalPosts,
	}
	generateHTML(w, data, "post", "layout", "navbar")
}

func adminPage(w http.ResponseWriter, r *http.Request) {
	// if logged in - redirect to admin page
	// if not - redirect to login page
	if admin.GlobalSession.Id != "" && admin.GlobalSession.Pw != "" {
		http.Redirect(w, r, "static/adminpage/manage.html", 302)
	} else {
		http.Redirect(w, r, "static/adminpage/login.html", 302)
	}
}
