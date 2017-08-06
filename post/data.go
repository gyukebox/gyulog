package post

import (
	"database/sql"
	"fmt"
	"log"

	// driver import
	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB
var TotalPosts int

// connect to db and get all posts at start of the server
func init() {
	ConnectDB()
	queryString := "select * from post order by id desc"
	rows, err := DB.Query(queryString)
	if err != nil {
		fmt.Print("At collecting data, ")
		log.Fatalln(err)
	}
	TotalPosts = 0
	defer rows.Close()
	for rows.Next() {
		TotalPosts++
	}
}

func ConnectDB() {
	var err error
	DB, err = sql.Open("mysql", "user:pw@tcp(127.0.0.1:3306)/gyulog")
	if err != nil {
		fmt.Print("At connecting DB, ")
		log.Fatalln(err)
	}
}

func (p *Post) insert() (err error) {
	queryString := "insert into post (title, summary, body) values (?, ?, ?)"
	rows, err := DB.Exec(queryString, p.Title, p.Summary, p.Body)
	if err != nil {
		log.Fatalln(err)
		return
	}
	n, err := rows.RowsAffected()
	if err != nil {
		log.Fatalln(err)
		return
	}
	TotalPosts++
	fmt.Printf("%d rows affected\n", n)
	return
}

func (p *Post) rowToPost(row *sql.Row) (err error) {
	err = row.Scan(&p.Id, &p.Title, &p.PublishedDate, &p.Summary, &p.Body)
	if err != nil {
		fmt.Print("At scanning row(s) : ")
		log.Fatalln(err)
	}
	return
}

func (p *Post) rowsToPost(rows *sql.Rows) (err error) {
	err = rows.Scan(&p.Id, &p.Title, &p.PublishedDate, &p.Summary, &p.Body)
	if err != nil {
		fmt.Print("At scanning row(s) : ")
		log.Fatalln(err)
	}
	return
}

func GetFivePosts(offset int) (posts []*Post, err error) {
	queryString := "select * from post order by id desc limit ?, 5"
	rows, err := DB.Query(queryString, offset*5)
	if err != nil {
		fmt.Print("At executing query, ")
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

func GetPostByTitle(title string) (post Post) {
	queryString := "select * from post where title = ?"
	row := DB.QueryRow(queryString, title)
	post = Post{}
	post.rowToPost(row)
	return
}

func GetPostById(id int) (post Post, err error) {
	queryString := "select * from post where id = ?"
	row := DB.QueryRow(queryString, id)
	post = Post{}
	err = post.rowToPost(row)
	return
}
