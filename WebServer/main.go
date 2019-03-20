//Web server
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
)

var appServer =  ""

//serves all static html to the web server by taking
//in the request, reading the request url, and serving
//the correct html based on the switch cases
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
		fileToServe = "gameSelect.html"
	case "/apexForm":
		fileToServe = "apexForm.html"
	case "/fortniteForm":
		fileToServe = "forniteForm.html"
	case "/hotsForm":
		fileToServe = "hotsForm.html"
	case "/apexSelect":
		fileToServe = "apexSelect.html"
	case "/fortniteSelect":
		fileToServe = "fortniteSelect.html"
	case "/hotsSelect":
		fileToServe = "hotsSelect.html"
	case "/viewApex":
		fileToServe = "viewApex.html"
	case "/viewFortnite":
		fileToServe = "viewFortnite.html"
	case "/viewHots":
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

//recursively parses html and returns all content that
//is within the tbody element
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

//takes in html node and returns string format
func renderNode(n *html.Node) string {
    var buf bytes.Buffer
    w := io.Writer(&buf)
    html.Render(w, n)
    return buf.String()
}

//parses request form, gets the requested user's game results,
//parses html table generated from GET call, and executes
//table into template for response
func Results(w http.ResponseWriter, r *http.Request)  {
	r.ParseForm()

	resp, _ := http.Get(appServer + "/view/" + r.FormValue( "game") + "/" + r.FormValue("user_name"))

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

//POSTs the request to the appServer, parses the request to,
//find game type, and responds with the game types select
//page
func Submit (w http.ResponseWriter, r *http.Request) {
	var fileToServe string

	_, _ = http.Post(appServer + "/submit", "application/x-www-form-urlencoded", r.Body)

	r.ParseForm()
	game := r.FormValue("game")
	switch game {
	case "apex":
		fileToServe = "apexSelect.html"
	case "fort":
		fileToServe = "fortniteSelect.html"
	case "hots":
		fileToServe = "hotsSelect.html"
	default:
		fileToServe = "gameSelect.html"
	}

	//redirect to game's select screen
	t, err := template.ParseFiles(fileToServe)
	if err != nil {
		log.Print("template parsing error: ", err)
	}
	err = t.Execute(w, nil)
	if err != nil {
		log.Print("template executing error: ", err)
	}
}

//create mux router to listen on port 80, handle
//all user interaction, and generate web content
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
	r.HandleFunc("/submit", Submit)

	log.Fatal(http.ListenAndServe(":80", r))
}
