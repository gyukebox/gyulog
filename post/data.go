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

func (p *Post) rowToPost(row *sql.Row) (err error) {
	err = row.Scan(&p.Id, &p.Title, &p.PublishedDate, &p.Summary, &p.Body)
	if err != nil {
		fmt.Print("At scanning row(s) : ")
		fmt.Println(err)
		log.Fatalln(err)
	}
	return
}

func (p *Post) rowsToPost(rows *sql.Rows) (err error) {
	err = rows.Scan(&p.Id, &p.Title, &p.PublishedDate, &p.Summary, &p.Body)
	if err != nil {
		fmt.Print("At scanning row(s) : ")
		fmt.Println(err)
		log.Fatalln(err)
	}
	return
}

func GetFivePosts(offset int) (posts []*Post, err error) {
	queryString := "select * from post order by id desc limit ?, 5"
	rows, err := DB.Query(queryString, offset*5)
	if err != nil {
		fmt.Print("At executing query : ")
		fmt.Println(err)
		log.Fatalln(err)
	}
	defer rows.Close()
	for rows.Next() {
		result := Post{}
		result.rowsToPost(rows)
		posts = append(posts, &result)
	}
	return
}

func GetAllPosts() (posts []*Post, err error) {
	queryString := "select * from post order by id desc"
	rows, err := DB.Query(queryString)
	if err != nil {
		fmt.Print("At Executing Query : ")
		fmt.Println(err)
		log.Fatalln(err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		result := Post{}
		result.rowsToPost(rows)
		posts = append(posts, &result)
	}
	return
}

func GetPostByTitle(title string) (post Post) {
	queryString := "select * from post where title = ?"
	row := DB.QueryRow(queryString, title)
	post = Post{}
	post.rowToPost(row)
	return
}

func GetPostById(id int) (post Post) {
	queryString := "select * from post where id = ?"
	row := DB.QueryRow(queryString, id)
	post = Post{}
	post.rowToPost(row)
	return
}
