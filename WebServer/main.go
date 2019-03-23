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

const AppServer =  ""

// serves all static html to the web server by taking
// in the request, reading the request url, and serving
// the correct html based on the switch cases.  Ensures
// that the request is coming from a source that is validated
func ServeStaticHTML(w http.ResponseWriter, r *http.Request) {
	var fileToServe string

	path := r.URL.String()
	switch path{
	case "/create":
		fileToServe = "create.html"
	case "/createError":
		fileToServe = "createError.html"
	case "/signinError":
		fileToServe = "signinError.html"
	case "/gameSelect":
		ValidateCookie(w, r)  // ensure user is validated to access
		fileToServe = "gameSelect.html"
	case "/apexForm":
		ValidateCookie(w, r)
		fileToServe = "apexForm.html"
	case "/fortniteForm":
		ValidateCookie(w, r)
		fileToServe = "fortniteForm.html"
	case "/hotsForm":
		ValidateCookie(w, r)
		fileToServe = "hotsForm.html"
	case "/apexSelect":
		ValidateCookie(w, r)
		fileToServe = "apexSelect.html"
	case "/fortniteSelect":
		ValidateCookie(w, r)
		fileToServe = "fortniteSelect.html"
	case "/hotsSelect":
		ValidateCookie(w, r)
		fileToServe = "hotsSelect.html"
	case "/viewApex":
		ValidateCookie(w, r)
		fileToServe = "viewApex.html"
	case "/viewFortnite":
		ValidateCookie(w, r)
		fileToServe = "viewFortnite.html"
	case "/viewHots":
		ValidateCookie(w, r)
		fileToServe = "viewHots.html"
	default:
		fileToServe = "index.html"
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

func Style(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "style.css")
}

func ApexLogo (w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "apex_legends_logo.jpg")
}

func FortLogo (w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "fortnite_logo.jpg")
}

func HotsLogo (w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "hots_logo.jpg")
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

	r.ParseForm()

	resp, _ := http.Get(AppServer + "/view/" + r.FormValue( "game") + "/" + r.FormValue("user_name"))

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

	t, err := template.ParseFiles("table.html")
	if err != nil {
		log.Print("template parsing error: ", err)
	}
	err = t.Execute(w, dataTable)
	if err != nil {
		log.Print("template executing error: ", err)
	}
}

// POSTs the request to the AppServer, parses the request to,
// find game type, and redirects to the games select page
func Submit (w http.ResponseWriter, r *http.Request) {
	_, _ = http.Post(AppServer + "/submit", "application/x-www-form-urlencoded", r.Body)

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

// POSTs the request to the AppServer, reads the response status
// code to determine if sign-in was good or bad, if good, sets a
// cookie that expires in 24hrs with username, and redirects
// to appropriate web page
func Signin (w http.ResponseWriter, r *http.Request) {
	resp, err := http.Post(AppServer + "/signin", "application/x-www-form-urlencoded", r.Body)

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

// POSTs the request to the AppServer, reads the response status
// code to determine if creation was good or bad, if good, sets a
// cookie that expires in 24hrs with username, and redirects
// to appropriate web page
func CreateAccount (w http.ResponseWriter, r *http.Request) {
	resp, err := http.Post(AppServer + "/createAccount", "application/x-www-form-urlencoded", r.Body)

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
	var valid bool

	for _, cookie := range r.Cookies() {
		if cookie.Name == "user" {
			valid = true
			return
		}
	}

	if !valid{
		http.Redirect(w, r, "/", 301)
	}
}

// create mux router to listen on port 80, handle
// all user interaction, and generate web content
func main() {
	r := mux.NewRouter() //create router

	//html pages
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
	r.HandleFunc("/viewApex", ServeStaticHTML)
	r.HandleFunc("/viewFortnite", ServeStaticHTML)
	r.HandleFunc("/viewHots", ServeStaticHTML)
	r.HandleFunc("/results", Results)

	//style sheet
	r.HandleFunc("/style.css", Style)

	//images
	r.HandleFunc("/hots_logo.jpg", HotsLogo)
	r.HandleFunc("/apex_legends_logo.jpg", ApexLogo)
	r.HandleFunc("/fortnite_logo.jpg", FortLogo)

	//form handling
	r.HandleFunc("/signin", Signin)
	r.HandleFunc("/createAccount", CreateAccount)
	r.HandleFunc("/submit", Submit)

	log.Fatal(http.ListenAndServe(":80", r))
}
