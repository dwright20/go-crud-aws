// Serverless App server functionality lambda
package main

import (
	"hiddenCreds"
	"database/sql"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/gorillamux"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
)

var muxLambda *gorillamux.GorillaMuxAdapter  // initialize mux lambda adapter

// initialize a mux router to handle requests and attach
// it to the lambda adapter
func init() {
	log.Printf("Mux cold start...")
	r := mux.NewRouter()
	r.HandleFunc("/signin", signIn).Methods("POST")
	r.HandleFunc("/createAccount", createAccount).Methods("POST")
	muxLambda = gorillamux.New(r)
}

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
		log.Printf(username + " signed in") //log user sign-in
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
		log.Printf(username + " account created") //log account creation
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

func Handler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// If no name is provided in the HTTP request body, throw an error
	return muxLambda.Proxy(req)
}

// start the lambda mux router
func main() {
	lambda.Start(Handler)
}
