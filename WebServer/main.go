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

func HomePage(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("index.html") 
	if err != nil { 
		log.Print("template parsing error: ", err) 
	}
	err = t.Execute(w, nil) 
	if err != nil { 
		log.Print("template executing error: ", err) 
	}
}

func Create(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("create.html") 
	if err != nil { 
		log.Print("template parsing error: ", err) 
	}
	err = t.Execute(w, nil) 
	if err != nil { 
		log.Print("template executing error: ", err) 
	}
}

func CreateError(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("createError.html") 
	if err != nil { 
		log.Print("template parsing error: ", err) 
	}
	err = t.Execute(w, nil) 
	if err != nil { 
		log.Print("template executing error: ", err) 
	}
}

func SigninError(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("signinError.html") 
	if err != nil { 
		log.Print("template parsing error: ", err) 
	}
	err = t.Execute(w, nil) 
	if err != nil { 
		log.Print("template executing error: ", err) 
	}
}

func GameSelect(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("gameSelect.html") 
	if err != nil { 
		log.Print("template parsing error: ", err) 
	}
	err = t.Execute(w, nil) 
	if err != nil { 
		log.Print("template executing error: ", err) 
	}
}

func ApexForm(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("apexForm.html") 
	if err != nil { 
		log.Print("template parsing error: ", err) 
	}
	err = t.Execute(w, nil) 
	if err != nil { 
		log.Print("template executing error: ", err) 
	}
}

func FortniteForm(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("fortniteForm.html") 
	if err != nil { 
		log.Print("template parsing error: ", err) 
	}
	err = t.Execute(w, nil) 
	if err != nil { 
		log.Print("template executing error: ", err) 
	}
}

func HotsForm(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("hotsForm.html") 
	if err != nil { 
		log.Print("template parsing error: ", err) 
	}
	err = t.Execute(w, nil) 
	if err != nil { 
		log.Print("template executing error: ", err) 
	}
}

func ApexSelect(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("apexSelect.html") 
	if err != nil { 
		log.Print("template parsing error: ", err) 
	}
	err = t.Execute(w, nil) 
	if err != nil { 
		log.Print("template executing error: ", err) 
	}
}

func FortniteSelect(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("fortniteSelect.html") 
	if err != nil { 
		log.Print("template parsing error: ", err) 
	}
	err = t.Execute(w, nil) 
	if err != nil { 
		log.Print("template executing error: ", err) 
	}
}

func HotsSelect(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("hotsSelect.html") 
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

func ViewApex (w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("viewApex.html") 
	if err != nil { 
		log.Print("template parsing error: ", err) 
	}
	err = t.Execute(w, nil) 
	if err != nil { 
		log.Print("template executing error: ", err) 
	}
}

func ViewFortnite (w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("viewFortnite.html") 
	if err != nil { 
		log.Print("template parsing error: ", err) 
	}
	err = t.Execute(w, nil) 
	if err != nil { 
		log.Print("template executing error: ", err) 
	}
}

func ViewHots (w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("viewHots.html") 
	if err != nil { 
		log.Print("template parsing error: ", err) 
	}
	err = t.Execute(w, nil) 
	if err != nil { 
		log.Print("template executing error: ", err) 
	}
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
//parses html table generated from get call, and executes
//table into template for response
func Results(w http.ResponseWriter, r *http.Request)  {
	r.ParseForm()

	resp, _ := http.Get("http://ec2-3-92-133-225.compute-1.amazonaws.com:8000/view/" + r.FormValue( "game") + "/" + r.FormValue("user_name"))

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

//create mux router to listen on port 80, handle
//all user interaction, and generate web content
func main() {
	r := mux.NewRouter() //create router

	//html pages
	r.HandleFunc("/", HomePage)
	r.HandleFunc("/create", Create)
	r.HandleFunc("/createError", CreateError)
	r.HandleFunc("/signinError", SigninError)
	r.HandleFunc("/gameSelect", GameSelect)
	r.HandleFunc("/apexForm", ApexForm)
	r.HandleFunc("/fortniteForm", FortniteForm)
	r.HandleFunc("/hotsForm", HotsForm)
	r.HandleFunc("/apexSelect", ApexSelect)
	r.HandleFunc("/fortniteSelect", FortniteSelect)
	r.HandleFunc("/hotsSelect", HotsSelect)
	r.HandleFunc("/viewApex", ViewApex)
	r.HandleFunc("/viewFortnite", ViewFortnite)
	r.HandleFunc("/viewHots", ViewHots)
	r.HandleFunc("/results", Results)

	//style sheet
	r.HandleFunc("/style.css", Style)

	//images
	r.HandleFunc("/hots_logo.jpg", HotsLogo)
	r.HandleFunc("/apex_legends_logo.jpg", ApexLogo)
	r.HandleFunc("/fortnite_logo.jpg", FortLogo)

	log.Fatal(http.ListenAndServe(":80", r))
}
