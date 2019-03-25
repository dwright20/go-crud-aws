// App server
package main

import (
	"game"
	"hiddenCreds"
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"time"
)

const CrudServer = ""

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
		fmt.Println(username + " signed in") //log user sign-in
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
		fmt.Println(username + " account created") //log account creation
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
		fmt.Println("No rows were returned!")
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
func createUser(user_username, user_password string) bool {
	creds := getCreds()
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		creds.Host, creds.Port, creds.User, creds.Password, creds.Dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Println(err)
		return false
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user_password), 8)

	sqlStatement := `
INSERT INTO postgres.public.users (user_name, user_pass)
VALUES ($1, $2)`
	_, err = db.Exec(sqlStatement, user_username, hashedPassword)
	if err != nil {
		log.Println(err)
		return false
	}

	db.Close()
	return true
}

// parses the form in the request, creates correct game, and
// POSTs the game to the CRUD server
func submit(_ http.ResponseWriter, r *http.Request){
	params := mux.Vars(r)
	r.ParseForm()

	if r.FormValue("game") == "apex" {
		game := game.NewApex(params["user"], time.Now().Format(time.RFC822), r.FormValue("game"),r.FormValue("result"),r.FormValue("legend"),r.FormValue("kills"),r.FormValue("placement"),r.FormValue("damage"),r.FormValue("time"),r.FormValue("teammates"))

		fmt.Println(game)//log created game

		b := new(bytes.Buffer)
		json.NewEncoder(b).Encode(game)
		//Post encoded game to CRUD server
		_, err := http.Post(CrudServer + "/create/apex", "application/json", b)

		if err !=nil{
			log.Print("Posting error: ", err)
		}
	} else if r.FormValue("game") == "fort" {
		game := game.NewFort(params["user"], time.Now().Format(time.RFC822),r.FormValue("game"),r.FormValue("result"),r.FormValue("kills"),r.FormValue("placement"),r.FormValue("mode"), r.FormValue("teammates"))

		fmt.Println(game)//log created game

		b := new(bytes.Buffer)
		json.NewEncoder(b).Encode(game)
		//Post encoded game to CRUD server
		_, err := http.Post(CrudServer + "/create/fort", "application/json", b)

		if err !=nil{
			log.Print("Posting error: ", err)
		}
	} else {
		game := game.NewHots(params["user"], time.Now().Format(time.RFC822), r.FormValue("game"),r.FormValue("result"),r.FormValue("hero"),r.FormValue("kills"),r.FormValue("deaths"),r.FormValue("assists"),r.FormValue("time"),r.FormValue("map"))

		fmt.Println(game)//log created game

		b := new(bytes.Buffer)
		json.NewEncoder(b).Encode(game)
		//Post encoded game to CRUD server
		_, err := http.Post(CrudServer + "/create/hots", "application/json", b)

		if err !=nil{
			log.Print("Posting error: ", err)
		}
	}
}

// sends a get request to the CRUD server to retrieve the
// specified user's game results and sends GET request contents
// back to the initial requesting server
func view(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	req := CrudServer + "/read/" + params["game"] + "/" + params["user"]
	fmt.Println("Reading " + params["user"] + "-" + params["game"])
	resp, _ := http.Get(req)
	resp.Write(w)
}

// create mux router to listen on port 8000 and handle
// various requests
func main() {
	r := mux.NewRouter()
	r.HandleFunc("/signin", signIn).Methods("POST")
	r.HandleFunc("/createAccount", createAccount).Methods("POST")
	r.HandleFunc("/submit/{user}", submit).Methods("POST")
	r.HandleFunc("/view/{game}/{user}", view).Methods("GET")
	log.Fatal(http.ListenAndServe(":8000", r))
}
