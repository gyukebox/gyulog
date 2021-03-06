package admin

import (
	"encoding/base64"
	"log"
	"net/http"

	"github.com/gyukebox/gyulog/post"
)

type Session struct {
	Id string
	Pw string
}

var GlobalSession Session

func Authenticate(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var id, pw string
	rows, err := post.DB.Query("select * from admin")
	defer rows.Close()
	if err != nil {
		log.Println("Executing function post.DB.Query() while executing Authenticate() in admin.go...")
		log.Fatalln(err)
	}
	rows.Next()
	rows.Scan(&id, &pw)

	idResult, err := base64.URLEncoding.DecodeString(id)
	if err != nil {
		log.Println("Executing function base64.URLEncoding.DecodeString() while executing Authenticate() in admin.go...")
		log.Fatalln(err)
	}
	pwResult, err := base64.URLEncoding.DecodeString(pw)
	if err != nil {
		log.Println("Executing function base64.URLEncoding.DecodeString() while executing Authenticate() in admin.go...")
		log.Fatalln(err)
	}

	if string(idResult) != r.FormValue("id") || string(pwResult) != r.FormValue("pw") {
		// login fail
		log.Println("Error : Incorrect ID or Password")
		http.Redirect(w, r, "../static/adminpage/login.html", 302)
	} else {
		// session store
		GlobalSession = Session{}
		GlobalSession.Id = base64.URLEncoding.EncodeToString([]byte(id))
		GlobalSession.Pw = base64.URLEncoding.EncodeToString([]byte(pw))
		http.Redirect(w, r, "/admin", 302)
	}
}

func Logout(w http.ResponseWriter, r *http.Request) {
	GlobalSession.Id = ""
	GlobalSession.Pw = ""
	http.Redirect(w, r, "/", 302)
}
