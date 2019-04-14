// App server
package main

import (
	"web_project/ecs_app/hiddenCreds"
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"time"
)

// get DB credentials
func getCreds() hiddenCreds.Creds{
	return hiddenCreds.GetCreds()
}

// parses the form in the request, checks if user is in
// the db and credentials are correct, and responds with
// a good status code and username or a bad status code
func signIn(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	data := r.PostForm
	username := data.Get("user_name")
	password := data.Get("user_pass")

	res := checkPassword(username, password)
	if res == true{
		log.Println(username + " signed in") //log user sign-in
		w.WriteHeader(200)
		w.Write([]byte(username))
	}else {
		w.WriteHeader(400)
	}
}

// parses the form in the request, checks if user exists
// already, and responds with a good status code and username
// or a bad status code
func createAccount(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	data := r.PostForm
	username := data.Get("user_name")
	password := data.Get("user_pass")

	res := createUser(username, password)
	if res == true{
		log.Println(username + " account created") //log account creation
		w.WriteHeader(200)
		w.Write([]byte(username))
	}else {
		w.WriteHeader(400)
	}
}

// take in a username and password and return a bool of
// the validation
func checkPassword(username, password string) bool {
	var dbPass string
	creds := getCreds()
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		creds.Host, creds.Port, creds.User, creds.Password, creds.Dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Println(err)
		return false
	}

	sqlStatement := `SELECT user_pass FROM postgres.public.users WHERE user_name=$1;`
	row := db.QueryRow(sqlStatement, username)

	switch err := row.Scan(&dbPass); err {
	case sql.ErrNoRows:
		log.Println("No rows were returned!")
	case nil:
		hashAndPass := bcrypt.CompareHashAndPassword([]byte(dbPass), []byte(password))
		if hashAndPass == nil{
			db.Close()
			return true
		} else {
			return false
		}
	default:
		log.Println(err)
	}

	return false
}

// take in a username and password and return a bool of
// the validation of creation
func createUser(username, password string) bool {
	creds := getCreds()
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		creds.Host, creds.Port, creds.User, creds.Password, creds.Dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Println(err)
		return false
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), 8)

	sqlStatement := `
INSERT INTO postgres.public.users (user_name, user_pass)
VALUES ($1, $2)`
	_, err = db.Exec(sqlStatement, username, hashedPassword)
	if err != nil {
		log.Println(err)
		return false
	}

	db.Close()
	return true
}

// responds to health check request with a good status code
func healthStatus(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(200)
}

// create mux router to listen on port 8000 and handle
// various requests
func main() {
	r := mux.NewRouter()
	r.HandleFunc("/signin", signIn).Methods("POST")
	r.HandleFunc("/createAccount", createAccount).Methods("POST")
	r.HandleFunc("/health", healthStatus).Methods("GET")

	srv := &http.Server{
		Handler: 		r,
		Addr:			":80",
		WriteTimeout: 	15 * time.Second,
		ReadTimeout:	15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
