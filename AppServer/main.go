//App server
package main

import (
	"awesomeProject/first_project/game"
	"awesomeProject/first_project/hiddenCreds"
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

//get DB credentials
func getCreds() hiddenCreds.Creds{
	return hiddenCreds.GetCreds()
}

//parses the form in the request, checks if user is in
//the db and credentials are correct, and redirects to
//an error page or main screen
func signIn(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	data := r.PostForm
	username := data.Get("user_name")
	password := data.Get("user_pass")

	res := checkPassword(username, password)
	if res == true{
		fmt.Println(username + " signed in") //log user sign-in
		http.Redirect(w, r, "http://ec2-3-95-30-158.compute-1.amazonaws.com/gameSelect", http.StatusSeeOther)
	}else {
		http.Redirect(w, r, "http://ec2-3-95-30-158.compute-1.amazonaws.com/signinError", http.StatusSeeOther)
	}
}

//parses the form in the request, checks if user exists
//already, and redirects to an error page or main screen
func createAccount(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	data := r.PostForm
	username := data.Get("user_name")
	password := data.Get("user_pass")

	res := createUser(username, password)
	if res == true{
		fmt.Println(username + " account created") //log account creation
		http.Redirect(w, r, "http://ec2-3-95-30-158.compute-1.amazonaws.com/gameSelect", http.StatusSeeOther)
	}else {
		http.Redirect(w, r, "http://ec2-3-95-30-158.compute-1.amazonaws.com/createError", http.StatusSeeOther)
	}
}

//take in a username and password and return a bool of
//the validation
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

//take in a username and password and return a bool of
//the validation of creation
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

//parses the form in the request, creates correct game, and
//POSTs the game to the CRUD server
func submit(w http.ResponseWriter, r *http.Request){
	r.ParseForm()

	if r.FormValue("game") == "apex" {
		game := game.NewApex(r.FormValue("user_name"), time.Now().Format(time.RFC822), r.FormValue("game"),r.FormValue("result"),r.FormValue("legend"),r.FormValue("kills"),r.FormValue("placement"),r.FormValue("damage"),r.FormValue("time"),r.FormValue("teammates"))

		fmt.Println(game)//log created game

		b := new(bytes.Buffer)
		json.NewEncoder(b).Encode(game)
		//Post encoded game to CRUD server
		_, _ = http.Post("http://ec2-174-129-90-38.compute-1.amazonaws.com:8000/create/apex", "application/json", b)

		//redirect to game's select screen
		http.Redirect(w, r, "http://ec2-3-95-30-158.compute-1.amazonaws.com/apexSelect", http.StatusSeeOther)
	} else if r.FormValue("game") == "fort" {
		game := game.NewFort(r.FormValue("user_name"), time.Now().Format(time.RFC822),r.FormValue("game"),r.FormValue("result"),r.FormValue("kills"),r.FormValue("placement"),r.FormValue("mode"), r.FormValue("teammates"))

		fmt.Println(game)//log created game

		b := new(bytes.Buffer)
		json.NewEncoder(b).Encode(game)
		//Post encoded game to CRUD server
		_, _ = http.Post("http://ec2-174-129-90-38.compute-1.amazonaws.com:8000/create/fort", "application/json", b)

		//redirect to game's select screen
		http.Redirect(w, r, "http://ec2-3-95-30-158.compute-1.amazonaws.com/fortniteSelect", http.StatusSeeOther)
	} else {
		game := game.NewHots(r.FormValue("user_name"), time.Now().Format(time.RFC822), r.FormValue("game"),r.FormValue("result"),r.FormValue("hero"),r.FormValue("kills"),r.FormValue("deaths"),r.FormValue("assists"),r.FormValue("time"),r.FormValue("map"))

		fmt.Println(game)//log created game

		b := new(bytes.Buffer)
		json.NewEncoder(b).Encode(game)
		//Post encoded game to CRUD server
		_, _ = http.Post("http://ec2-174-129-90-38.compute-1.amazonaws.com:8000/create/hots", "application/json", b)

		//redirect to game's select screen
		http.Redirect(w, r, "http://ec2-3-95-30-158.compute-1.amazonaws.com/hotsSelect", http.StatusSeeOther)
	}
}

//sends a get request to the CRUD server to retrieve the
//specified user's game results and sends GET request contents
//back to the initial requesting server
func view(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	req := "http://ec2-174-129-90-38.compute-1.amazonaws.com:8000/read/" + params["game"] + "/" + params["user"]
	fmt.Println("Reading " + params["user"] + "-" + params["game"])
	resp, _ := http.Get(req)
	resp.Write(w)
}

//create mux router to listen on port 8000 and handle
//various requests
func main() {
	r := mux.NewRouter()
	r.HandleFunc("/signin", signIn).Methods("POST")
	r.HandleFunc("/createAccount", createAccount).Methods("POST")
	r.HandleFunc("/submit", submit).Methods("POST")
	r.HandleFunc("/view/{game}/{user}", view).Methods("GET")
	log.Fatal(http.ListenAndServe(":8000", r))
}
