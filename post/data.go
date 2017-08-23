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
		log.Println("Executing function DB.Query while executing function init() at data.go...")
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
	DB, err = sql.Open("mysql", "username:password@tcp(dbaddress:port)/dbname")
	if err != nil {
		log.Println("Executing function sql.Open() while executing ConnectDB() at data.go...")
		log.Fatalln(err)
	}
}

func (p *Post) insert() (err error) {
	queryString := "insert into post (title, summary, body) values (?, ?, ?)"
	rows, err := DB.Exec(queryString, p.Title, p.Summary, p.Body)
	if err != nil {
		log.Println("Executing function DB.Exec() while executing insert() in data.go...")
		log.Fatalln(err)
		return
	}

	n, err := rows.RowsAffected()
	if err != nil {
		log.Println("Executing function rows.RowsAffected() while executing insert() in data.go...")
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
		log.Println("Executing function row.Scan() while executing rowToPost() in data.go...")
		log.Fatalln(err)
	}
	return
}

func (p *Post) rowsToPost(rows *sql.Rows) (err error) {
	err = rows.Scan(&p.Id, &p.Title, &p.PublishedDate, &p.Summary, &p.Body)
	if err != nil {
		log.Println("Executing function rows.Scan() while executing rowsToPost() in data.go...")
		log.Fatalln(err)
	}
	return
}

func GetFivePosts(offset int) (posts []*Post, err error) {
	queryString := "select * from post order by id desc limit ?, 5"
	rows, err := DB.Query(queryString, offset*5)
	if err != nil {
		log.Println("Executing function DB.Query() while executing GetFivePosts() in data.go...")
		log.Fatalln(err)
	}

	defer rows.Close()
	for rows.Next() {
		result := Post{}
		err = result.rowsToPost(rows)
		if err != nil {
			log.Println("Executing function result.rowsToPost() while executing GetFivePosts() in data.go...")
			log.Fatalln(err)
		}
		posts = append(posts, &result)
	}
	return
}

func GetPostByTitle(title string) (post Post) {
	queryString := "select * from post where title = ?"
	row := DB.QueryRow(queryString, title)
	post = Post{}
	err := post.rowToPost(row)
	if err != nil {
		log.Println("Executing function post.rowToPost() while executing GetPostByTitle() in data.go...")
		log.Fatalln(err)
	}
	return
}

func GetPostById(id int) (post Post, err error) {
	queryString := "select * from post where id = ?"
	row := DB.QueryRow(queryString, id)
	post = Post{}
	err = post.rowToPost(row)
	if err != nil {
		log.Println("Executing function post.rowToPost() while executing GetPostById() in data.go...")
		log.Fatalln(err)
	}
	return
}
