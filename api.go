package api

import (
	"database/sql"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

var Database *sql.DB

type Users struct {
	Users []User `json:"users"`
}

type User struct {
	ID    int    "json:id"
	Name  string "json:username"
	Email string "json:email"
	First string "json:first"
	Last  string "json:last"
}

func StartServer() {
	db, err := sql.Open("mysql", "root@/social_network")
	if err != nil {
	}
	Database = db
	routes := mux.NewRouter()
	http.Handle("/", routes)
	http.ListenAndServe(":85", nil)
}
