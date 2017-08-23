package post

import (
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/russross/blackfriday"
)

// Post - Overall posts parsed by parts
type Post struct {
	Id            int
	Title         string
	PublishedDate string
	Summary       string
	Body          string
}

// 1. markdown으로 post 작성후 submit 누르면 (파일 업로드?)
// 2. post 받아와서 string으로 변환
// 3. Post struct 형식에 맞게 파싱
// 4. DB 에 저장

func GetPost(w http.ResponseWriter, r *http.Request) {
	// 업로드된 파일 받아오기
	err := r.ParseMultipartForm(1024)
	if err != nil {
		log.Println("Executing function r.ParseMultipartForm() while executing function GetPost() in posts.go...")
		log.Fatalln(err)
	}

	fileHeader := r.MultipartForm.File["post"][0]
	file, err := fileHeader.Open()
	if err != nil {
		log.Println("Executing function fileHeader.Open() while executing function GetPost() in posts.go...")
		log.Fatalln(err)
	}

	// post 받아와서 string으로 변환
	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Println("Executing function ioutil.ReadAll() while executing function GetPost() in posts.go...")
		log.Fatalln(err)
	}

	post := markdownToPost(data)
	err = post.insert()
	if err != nil {
		log.Println("Executing function post.Insert() while executing GetPost() in posts.go...")
		log.Fatalln(err)
	}

	http.Redirect(w, r, "/", 302)
}

func markdownToPost(data []byte) Post {
	rendered := string(blackfriday.MarkdownCommon(data))
	result := Post{}

	// 제목 찾는 인덱스
	startIndex := strings.Index(rendered, "<h1>") + 4
	endIndex := strings.Index(rendered, "</h1>")
	result.Title = rendered[startIndex:endIndex]

	// 요약 찾는 인덱스
	startIndex = strings.Index(rendered, "<p>") + 3
	endIndex = strings.Index(rendered, "</p>")
	result.Summary = rendered[startIndex:endIndex]

	startIndex = endIndex
	result.Body = rendered[startIndex:]
	return result
}
