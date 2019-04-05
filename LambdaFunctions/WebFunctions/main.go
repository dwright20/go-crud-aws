// Server-less web server
package main

import (
	"bytes"
	"errors"
	"github.com/GeertJohan/go.rice"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/gorillamux"
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

const server = ""  // API server address
const site = ""  // Front-end gateway
var muxLambda *gorillamux.GorillaMuxAdapter  // initialize mux lambda adapter
var staticBox *rice.Box  // initialize box that will store static web content

// initialize a mux router to handle requests and attach
// it to the lambda adapter
func init() {
	log.Printf("Web server starting...")

	staticBox = rice.MustFindBox("websites")

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

	muxLambda = gorillamux.New(r)
}

// serves all static html to the web server by taking
// in the request, reading the request url, and serving
// the correct html based on the switch cases.  Ensures
// that the request is coming from a source that is validated
func ServeStaticHTML(w http.ResponseWriter, r *http.Request) {
	var fileToServe string

	path := r.URL.String()
	path = strings.TrimPrefix(path, "https://aws-serverless-go-api.com")  // remove prefix

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
	default:
		fileToServe = "index.html"
	}

	log.Printf("serving webpage: %s", fileToServe)

	templateData, _ := staticBox.String(fileToServe)  // get content from box

	t, err := template.New("t").Parse(templateData)
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

	cookie := GetCookieValue(r)  // pull username

	params := mux.Vars(r)  // pull game

	log.Println("getting results...")

	resp, _ := http.Get(server + "/view/" + params["game"] + "/" + cookie)

	log.Println("processing results...")

	body, _ := ioutil.ReadAll(resp.Body)

	res := string(body)

	doc, err := html.Parse(strings.NewReader(res))
	bn, err := getBody(doc)
	if err != nil {
		log.Println(err)
	}
	bod := renderNode(bn)

	//struct to pass generated content into template
	dataTable := struct {
		Table template.HTML
	}{
		Table: template.HTML(bod),
	}

	templateData, _ := staticBox.String("table.html")

	t, err := template.New("t").Parse(templateData)
	if err != nil {
		log.Print("template parsing error: ", err)
	}
	err = t.Execute(w, dataTable)
	if err != nil {
		log.Print("template executing error: ", err)
	}
}

// POSTs the request to the api server, parses
// the request to find game type, and redirects to the
// games select page
func Submit (w http.ResponseWriter, r *http.Request) {
	user := GetCookieValue(r)
	log.Println("submitting results...")
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

	http.Redirect(w, r, site + redirect, 301)
}

// POSTs the request to the api server, reads the
// response status code to determine if sign-in was good
// or bad, if good, sets a cookie that expires in 24hrs
// with username, and redirects to appropriate web page
func Signin (w http.ResponseWriter, r *http.Request) {
	log.Println("trying to sign in...")
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
		log.Printf("%s signed in...", user)
		http.SetCookie(w, &cookie)
		http.Redirect(w, r, site + "/gameSelect", 301)
	}else{
		log.Println("signin failed...")
		http.Redirect(w, r, site + "/signinError", 301)
	}
}

// POSTs the request to the api server, reads the
// response status code to determine if creation was good
// or bad, if good, sets a cookie that expires in 24hrs with
// username, and redirects to appropriate web page
func CreateAccount (w http.ResponseWriter, r *http.Request) {
	log.Println("trying to create account...")
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
		log.Printf("%s signed in...", user)
		http.SetCookie(w, &cookie)
		http.Redirect(w, r, site + "/gameSelect", 301)
	}else{
		log.Println("signin failed...")
		http.Redirect(w, r, site + "/createError", 301)
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

	http.Redirect(w, r, site + "/", 301)
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
	log.Println("catching-all...")
	if GetCookieValue(r) == "" {
		http.Redirect(w, r, site + "/", 301)
	} else {
		http.Redirect(w, r, site + "/gameSelect", 301)
	}
}

func Handler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return muxLambda.Proxy(req)
}

// start the lambda mux router
func main() {
	lambda.Start(Handler)
}

