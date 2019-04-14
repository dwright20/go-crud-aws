// Web server
package main

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"golang.org/x/net/html"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

const AppServer = ""  // App Server address
const CrudServer = ""  //CRUD Server address
const FailOverApi = ""  // Fail-over server address

var server string  // server that will be routed too

// serves all static html to the web server by taking
// in the request, reading the request url, and serving
// the correct html based on the switch cases.  Ensures
// that the request is coming from a source that is validated
func ServeStaticHTML(w http.ResponseWriter, r *http.Request) {
	fileToServe := "static/"

	path := r.URL.String()
	switch path{
	case "/create":
		fileToServe += "create.html"
	case "/createError":
		fileToServe += "createError.html"
	case "/signinError":
		fileToServe += "signinError.html"
	case "/gameSelect":
		ValidateCookie(w, r)  // ensure user is validated to access
		fileToServe += "gameSelect.html"
	case "/apexForm":
		ValidateCookie(w, r)
		fileToServe += "apexForm.html"
	case "/fortniteForm":
		ValidateCookie(w, r)
		fileToServe += "fortniteForm.html"
	case "/hotsForm":
		ValidateCookie(w, r)
		fileToServe += "hotsForm.html"
	case "/apexSelect":
		ValidateCookie(w, r)
		fileToServe += "apexSelect.html"
	case "/fortniteSelect":
		ValidateCookie(w, r)
		fileToServe += "fortniteSelect.html"
	case "/hotsSelect":
		ValidateCookie(w, r)
		fileToServe += "hotsSelect.html"
	default:
		fileToServe += "index.html"
	}
	t, err := template.ParseFiles(fileToServe)
	if err != nil {
		log.Print("template parsing error: ", err)
	}
	err = t.Execute(w, nil)
	if err != nil {
		log.Print("template executing error: ", err)
	}
}

// recursively parses html and returns all content that
// is within the tbody element
func getBody(doc *html.Node) (*html.Node, error) {
	var b *html.Node
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "tbody" {
			b = n
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
	if b != nil {
		return b, nil
	}
	return nil, errors.New("Missing <tbody> in the node tree")
}

// takes in html node and returns string format
func renderNode(n *html.Node) string {
	var buf bytes.Buffer
	w := io.Writer(&buf)
	html.Render(w, n)
	return buf.String()
}

// parses request form, gets the requested user's game results,
// parses html table generated from GET call, and executes
// table into template for response
func Results(w http.ResponseWriter, r *http.Request)  {
	ValidateCookie(w, r)  // ensure user is validated to access

	healthCheck(false)  // determine which server to use

	cookie := GetCookieValue(r)  // pull username

	params := mux.Vars(r)  // pull game

	resp, _ := http.Get(server + "/view/" + params["game"] + "/" + cookie)

	body, _ := ioutil.ReadAll(resp.Body)

	res := string(body)

	doc, err := html.Parse(strings.NewReader(res))
	bn, err := getBody(doc)
	if err != nil {
		fmt.Println(err)
	}
	bod := renderNode(bn)

	//struct to pass generated content into template
	dataTable := struct {
		Table template.HTML
	}{
		Table: template.HTML(bod),
	}

	t, err := template.ParseFiles("static/table.html")
	if err != nil {
		log.Print("template parsing error: ", err)
	}
	err = t.Execute(w, dataTable)
	if err != nil {
		log.Print("template executing error: ", err)
	}
}

// POSTs the request to the appropriate server, parses
// the request to, find game type, and redirects to the
// games select page
func Submit (w http.ResponseWriter, r *http.Request) {
	user := GetCookieValue(r)
	healthCheck(false)  // determine which server to use
	_, _ = http.Post(server + "/submit/" + user, "application/x-www-form-urlencoded", r.Body)

	r.ParseForm()
	game := r.FormValue("game")
	var redirect string
	switch game {
	case "apex":
		redirect = "/apexSelect"
	case "fort":
		redirect = "/fortniteSelect"
	case "hots":
		redirect = "/hotsSelect"
	default:
		redirect = "/gameSelect"
	}
	http.Redirect(w, r, redirect, 301)
}

// POSTs the request to the appropriate server, reads the
// response status code to determine if sign-in was good
// or bad, if good, sets a cookie that expires in 24hrs
// with username, and redirects to appropriate web page
func Signin (w http.ResponseWriter, r *http.Request) {
	healthCheck(true)  // determine which server to use
	resp, err := http.Post(server + "/signin", "application/x-www-form-urlencoded", r.Body)

	if err != nil{
		log.Println(err)
	}
	if resp.StatusCode == 200 && err == nil{
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil{
			log.Println(err)
		}
		user := string(body)
		expire := time.Now().Add(time.Hour * 24)
		cookie := http.Cookie{
			Name:		"user",
			Value:		user,
			Expires:	expire,
		}
		http.SetCookie(w, &cookie)
		http.Redirect(w, r, "/gameSelect", 301)
	}else{
		http.Redirect(w, r, "/signinError", 301)
	}
}

// POSTs the request to the appropriate server, reads the
// response status code to determine if creation was good
// or bad, if good, sets a cookie that expires in 24hrs with
// username, and redirects to appropriate web page
func CreateAccount (w http.ResponseWriter, r *http.Request) {
	healthCheck(true)  // determine which server to use
	resp, err := http.Post(server + "/createAccount", "application/x-www-form-urlencoded", r.Body)

	if err != nil{
		log.Println(err)
	}
	if resp.StatusCode == 200 && err == nil{
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil{
			log.Println(err)
		}
		user := string(body)
		expire := time.Now().Add(time.Hour * 24)
		cookie := http.Cookie{
			Name:		"user",
			Value:		user,
			Expires:	expire,
		}
		http.SetCookie(w, &cookie)
		http.Redirect(w, r, "/gameSelect", 301)
	}else{
		http.Redirect(w, r, "/createError", 301)
	}
}

// takes in a request and response writer and determines if
// the requesting client is already validated by looking at
// the cookies. If client does not have web server's cookie,
// it is redirected to sign-in screen.  If it does, function
// returns and lets handler continue processing
func ValidateCookie (w http.ResponseWriter,  r *http.Request) {
	for _, cookie := range r.Cookies() {
		if cookie.Name == "user" {
			return
		}
	}

	http.Redirect(w, r, "/", 301)
}

// takes in a request and reads the cookie content to find
// the cookie value of the cookie created upon sign-in and
// returns it as a string
func GetCookieValue (r *http.Request) string {
	var val string

	for _, cookie := range r.Cookies() {
		if cookie.Name == "user" {
			val = cookie.Value
		}
	}
	return val
}

// takes in all requests that are for a path not specified
// and redirects them to the game select page if they
// already have a cookie or the sign-in page if they do not
func NotFound (w http.ResponseWriter, r *http.Request) {
	if GetCookieValue(r) == "" {
		http.Redirect(w, r, "/", 301)
	} else {
		http.Redirect(w, r, "/gameSelect", 301)
	}
}

// sends a GET request to the primary server to determine if it
// is running.  If it receives an error or anything other
// than a 200 status code, it sets the server to the failover,
// otherwise it sets it to the primary server.
func healthCheck (app bool){
	// determine which server to check
	if app {
		resp, err := http.Get(AppServer + "/health")

		if err != nil || resp.StatusCode != 200 {
			log.Print("App Server down!\n")
			server = FailOverApi
		} else {
			server = AppServer
		}
	}else{
		resp, err := http.Get(CrudServer + "/health")

		if err != nil || resp.StatusCode != 200 {
			log.Print("Crud Server down!\n")
			server = FailOverApi
		} else {
			server = CrudServer
		}
	}
}

// create mux router to listen on port 80, handle
// all user interaction, and generate web content
func main() {
	r := mux.NewRouter() //create router

	// html pages
	r.HandleFunc("/", ServeStaticHTML)
	r.HandleFunc("/create", ServeStaticHTML)
	r.HandleFunc("/createError", ServeStaticHTML)
	r.HandleFunc("/signinError", ServeStaticHTML)
	r.HandleFunc("/gameSelect", ServeStaticHTML)
	r.HandleFunc("/apexForm", ServeStaticHTML)
	r.HandleFunc("/fortniteForm", ServeStaticHTML)
	r.HandleFunc("/hotsForm", ServeStaticHTML)
	r.HandleFunc("/apexSelect", ServeStaticHTML)
	r.HandleFunc("/fortniteSelect", ServeStaticHTML)
	r.HandleFunc("/hotsSelect", ServeStaticHTML)
	r.HandleFunc("/results/{game}", Results)

	// form handling
	r.HandleFunc("/signin", Signin)
	r.HandleFunc("/createAccount", CreateAccount)
	r.HandleFunc("/submit", Submit)

	// catch-all
	r.NotFoundHandler = http.HandlerFunc(NotFound)

	srv := &http.Server{
		Handler: 		r,
		Addr:			":80",
		WriteTimeout: 	15 * time.Second,
		ReadTimeout:	15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
