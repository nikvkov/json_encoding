package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

var database *sql.DB

type API struct {
	Message string "json:message"
}

type User struct {
	ID    int    "json:id"
	Name  string "json:username"
	Email string "json:email"
	First string "json:first"
	Last  string "json:last"
}

type Users struct {
	Users []User `json:"users"`
}

/*********************************/
func Hello(w http.ResponseWriter, r *http.Request) {
	urlParams := r.URL.Query()
	name := urlParams.Get(":name")
	HelloMessage := "Hello, " + name
	message := API{HelloMessage}
	output, err := json.Marshal(message)
	if err != nil {
		fmt.Println("Something went wrong!")
	}
	fmt.Fprintf(w, string(output))
}

func UserCreate(w http.ResponseWriter, r *http.Request) {
	NewUser := User{}
	NewUser.Name = r.FormValue("user")
	NewUser.Email = r.FormValue("email")
	NewUser.First = r.FormValue("first")
	NewUser.Last = r.FormValue("last")
	output, err := json.Marshal(NewUser)
	fmt.Println(string(output))
	if err != nil {
		fmt.Println("Something went wrong!")
	}
	sql := "INSERT INTO users set user_nickname='" + NewUser.Name +
		"', user_first='" + NewUser.First + "', user_last='" +
		NewUser.Last + "', user_email='" + NewUser.Email + "'"
	q, err := database.Exec(sql)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(q)
}

func UsersRetrieve(w http.ResponseWriter, r *http.Request) {
	log.Println("starting retrieval")
	start := 0
	limit := 10
	next := start + limit
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Link", "<http://localhost:85/api/users?start="+string(next)+"; rel=\"next\"")
	rows, _ := database.Query("select * from users LIMIT 10")
	Response := Users{}
	for rows.Next() {
		user := User{}
		rows.Scan(&user.ID, &user.Name, &user.First, &user.Last,
			&user.Email)
		Response.Users = append(Response.Users, user)
	}

	output, _ := json.Marshal(Response)
	fmt.Fprintln(w, string(output))
}

func GetUser(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Pragma", "no-cache")
	urlParams := mux.Vars(r)
	id := urlParams["id"]
	ReadUser := User{}
	err := database.QueryRow("select id, name, lastname, email from persons where id=?", id).Scan(&ReadUser.ID, &ReadUser.Name,
		&ReadUser.Last, &ReadUser.Email)
	switch {
	case err == sql.ErrNoRows:
		fmt.Fprintf(w, "No such user")
	case err != nil:
		log.Fatal(err)
		fmt.Fprintf(w, "Error")
	default:
		output, _ := json.Marshal(ReadUser)
		fmt.Fprintf(w, string(output))
	}
}

/*********************************/

func main() {

	db, err := sql.Open("mysql", "root:developer@/test_go")

	if err != nil {
		log.Println(err)
	}
	database = db
	defer db.Close()

	routes := mux.NewRouter()
	routes.HandleFunc("/api/users", UserCreate).Methods("POST")
	routes.HandleFunc("/api/users", UsersRetrieve).Methods("GET")
	http.Handle("/", routes)
	http.ListenAndServe(":85", nil)
}

/****************Utils*******************************/
