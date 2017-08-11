package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/gyukebox/gyulog/admin"

	"github.com/gyukebox/gyulog/post"
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

func main() {
	defer post.DB.Close()

	mux := http.DefaultServeMux
	files := http.FileServer(http.Dir("./static"))
	server := http.Server{
		Addr:    "127.0.0.1:5000",
		Handler: mux,
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
