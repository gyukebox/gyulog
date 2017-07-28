package post

import (
	"database/sql"
	"fmt"
	"log"

	// driver import
	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func init() {
	var err error
	DB, err = sql.Open("mysql", "root:biss9541@tcp(127.0.0.1:3306)/gyulog")
	if err != nil {
		fmt.Println(err)
		log.Fatalln(err)
	} else {
		fmt.Println("connected to database!")
	}
	// defer DB.Close()
}

func (p *Post) insert() (err error) {
	queryString := "insert into post (title, summary, body) values (?, ?, ?)"
	result, err := DB.Exec(queryString, p.Title, p.Summary, p.Body)
	if err != nil {
		fmt.Println(err)
		log.Fatalln(err)
		return
	}
	n, err := result.RowsAffected()
	if err != nil {
		fmt.Println(err)
		log.Fatalln(err)
		return
	}

	fmt.Printf("%d rows affected\n", n)
	return
}

func GetAllPosts() (posts []*Post, err error) {
	queryString := "select * from post"
	rows, err := DB.Query(queryString)
	if err != nil {
		fmt.Print("At Executing Query : ")
		fmt.Println(err)
		log.Fatalln(err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var (
			id      int
			title   string
			summary string
			body    string
			date    string
		)

		err = rows.Scan(&id, &title, &date, &summary, &body)
		if err != nil {
			fmt.Print("At scanning resultset : ")
			fmt.Println(err)
			log.Fatalln(err)
			return
		}

		result := Post{}
		result.Body = body
		result.PublishedDate = date
		result.Title = title
		result.Summary = summary

		posts = append(posts, &result)
	}
	return
}

func GetPostByTitle(title string) (post Post) {
	queryString := "select * from post where title = ?"
	row := DB.QueryRow(queryString, title)

	post = Post{}

	var (
		id        int
		postTitle string
		body      string
		date      string
		summary   string
	)

	row.Scan(&id, &postTitle, &date, &summary, &body)

	post.Body = body
	post.PublishedDate = date
	post.Summary = summary
	post.Title = title

	return
}
