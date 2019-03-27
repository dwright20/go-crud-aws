// Serverless App server functionality lambda
package main

import (
	"game"
	"hiddenCreds"
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/awslabs/aws-lambda-go-api-proxy/gorillamux"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"html/template"
	"log"
	"net/http"
	"reflect"
	"time"
)

var muxLambda *gorillamux.GorillaMuxAdapter  // initialize mux lambda adapter

// initialize a mux router to handle requests and attach
// it to the lambda adapter
func init() {
	log.Printf("Mux cold start...")
	r := mux.NewRouter()
	r.HandleFunc("/signin", signIn).Methods("POST")
	r.HandleFunc("/createAccount", createAccount).Methods("POST")
	r.HandleFunc("/submit/{user}", submit).Methods("POST")
	r.HandleFunc("/view/{game}/{user}", view).Methods("GET")
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

// parses the form in the request, creates correct game, and
// uploads game to the DynamoDB
func submit(w http.ResponseWriter, r *http.Request){
	params := mux.Vars(r)
	r.ParseForm()

	log.Print("starting session...")

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
	})

	if err != nil{
		log.Print(err.Error())
	}

	// create dynamodb client
	svc := dynamodb.New(sess)

	if r.FormValue("game") == "apex" {
		game := game.NewApex(params["user"], time.Now().Format(time.RFC822), r.FormValue("game"),r.FormValue("result"),r.FormValue("legend"),r.FormValue("kills"),r.FormValue("placement"),r.FormValue("damage"),r.FormValue("time"),r.FormValue("teammates"))

		log.Print(game)//log created game

		av, err := dynamodbattribute.MarshalMap(&game)

		input := &dynamodb.PutItemInput{
			Item: av,
			TableName: aws.String("results-apex"),
		}

		_, err = svc.PutItem(input) //put item in db

		//log error if applicable
		if err != nil {
			log.Println("Got error calling PutItem:")
			log.Println(err.Error())
		}
		log.Println("Added to Dynamodb")

	} else if r.FormValue("game") == "fort" {
		game := game.NewFort(params["user"], time.Now().Format(time.RFC822),r.FormValue("game"),r.FormValue("result"),r.FormValue("kills"),r.FormValue("placement"),r.FormValue("mode"), r.FormValue("teammates"))

		log.Print(game)//log created game

		av, err := dynamodbattribute.MarshalMap(&game)

		input := &dynamodb.PutItemInput{
			Item: av,
			TableName: aws.String("results-fort"),
		}

		_, err = svc.PutItem(input)

		if err != nil {
			log.Println("Got error calling PutItem:")
			log.Println(err.Error())
		}
		log.Println("Added to Dynamodb")

	} else {
		game := game.NewHots(params["user"], time.Now().Format(time.RFC822), r.FormValue("game"),r.FormValue("result"),r.FormValue("hero"),r.FormValue("kills"),r.FormValue("deaths"),r.FormValue("assists"),r.FormValue("time"),r.FormValue("map"))

		log.Print(game)//log created game

		av, err := dynamodbattribute.MarshalMap(&game)

		input := &dynamodb.PutItemInput{
			Item: av,
			TableName: aws.String("results-hots"),
		}

		_, err = svc.PutItem(input)

		if err != nil {
			log.Println("Got error calling PutItem:")
			log.Println(err.Error())
		}
		log.Println("Added to Dynamodb")

	}
	w.WriteHeader(200)
}

// read requesting users game results, generate html table
// of the results, and encode results and send them in
// response body
func view(w http.ResponseWriter, r *http.Request)  {
	params := mux.Vars(r)
	var templateFuncs = template.FuncMap{"rangeStruct": RangeStructer}

	//HTML template where generated content will go
	var htmlTemplate = `<!DOCTYPE html>

<html>
<head>
</head>
<body>
<div id="main">
    <table style="width: 100%">
    {{range .}}<tr>
    {{range rangeStruct .}} <td>{{.}}</td>
    {{end}}</tr>
    {{end}}
    </table>
</div>
</body>
</html>`

	var tpl bytes.Buffer

	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1")},
	)

	svc := dynamodb.New(sess)

	if params["game"] == "apex" {
		input := &dynamodb.QueryInput{
			TableName:	aws.String("results-apex"),
			KeyConditions: map[string]*dynamodb.Condition{
				"username": {
					ComparisonOperator: aws.String("EQ"),
					AttributeValueList: []*dynamodb.AttributeValue{
						{
							S: aws.String(params["user"]),
						},
					},
				},
			},
		}

		var resp, err = svc.Query(input)
		if err != nil {
			log.Println(err)
		}

		var games []game.Apex

		gms := []game.Apex{}

		dynamodbattribute.UnmarshalListOfMaps(resp.Items, &gms)

		headers := game.NewApex("User", "Date", "Game", "W/L", "Legend", "Kills", "Place", "Damage", "Time", "Team")

		games = append(games, headers)

		games = append(games, gms...)

		t := template.New("t").Funcs(templateFuncs) //create template with function to generate content

		t, err = t.Parse(htmlTemplate)
		if err != nil {
			panic(err)
		}

		err = t.Execute(&tpl, games) //execute template and pass slice of results into template function
		if err != nil {
			panic(err)
		}

		results := tpl.String() //convert generated html content into string

		//encode string of html content and write to response
		b := new(bytes.Buffer)
		encoder := json.NewEncoder(b)
		encoder.SetEscapeHTML(false)
		encoder.Encode(results)

		b.WriteTo(w)
		log.Println(params["user"] + " Apex data retrieved.")
	} else if params["game"] == "fort" {
		input := &dynamodb.QueryInput{
			TableName:	aws.String("results-fort"),
			KeyConditions: map[string]*dynamodb.Condition{
				"username": {
					ComparisonOperator: aws.String("EQ"),
					AttributeValueList: []*dynamodb.AttributeValue{
						{
							S: aws.String(params["user"]),
						},
					},
				},
			},
		}

		var resp, err = svc.Query(input)
		if err != nil {
			log.Println(err)
		}

		var games []game.Fort

		gms := []game.Fort{}

		dynamodbattribute.UnmarshalListOfMaps(resp.Items, &gms)

		headers := game.NewFort("User", "Date", "Game", "W/L", "Kills", "Place", "Mode", "Team")

		games = append(games, headers)

		games = append(games, gms...)

		t := template.New("t").Funcs(templateFuncs)

		t, err = t.Parse(htmlTemplate)
		if err != nil {
			panic(err)
		}

		err = t.Execute(&tpl, games)
		if err != nil {
			panic(err)
		}

		results := tpl.String()

		b := new(bytes.Buffer)
		encoder := json.NewEncoder(b)
		encoder.SetEscapeHTML(false)
		encoder.Encode(results)

		b.WriteTo(w)
		log.Println(params["user"] + " Fort data retrieved.")
	} else if params["game"] == "hots"{
		input := &dynamodb.QueryInput{
			TableName:	aws.String("results-hots"),
			KeyConditions: map[string]*dynamodb.Condition{
				"username": {
					ComparisonOperator: aws.String("EQ"),
					AttributeValueList: []*dynamodb.AttributeValue{
						{
							S: aws.String(params["user"]),
						},
					},
				},
			},
		}

		var resp, err = svc.Query(input)
		if err != nil {
			log.Println(err)
		}

		var games []game.Hots

		gms := []game.Hots{}

		dynamodbattribute.UnmarshalListOfMaps(resp.Items, &gms)

		headers := game.NewHots("User", "Date", "Game", "W/L", "Hero", "Kills", "Deaths", "Assists", "Time", "Map")

		games = append(games, headers)

		games = append(games, gms...)

		t := template.New("t").Funcs(templateFuncs)

		t, err = t.Parse(htmlTemplate)
		if err != nil {
			panic(err)
		}

		err = t.Execute(&tpl, games)
		if err != nil {
			panic(err)
		}

		results := tpl.String()
		b := new(bytes.Buffer)
		encoder := json.NewEncoder(b)
		encoder.SetEscapeHTML(false)
		encoder.Encode(results)

		b.WriteTo(w)
		log.Println(params["user"] + " Hots data retrieved.")
	}

}

// function to iterate through range of game results and
// fill html table template
func RangeStructer(args ...interface{}) []interface{} {
	if len(args) == 0 {
		return nil
	}

	v := reflect.ValueOf(args[0])
	if v.Kind() != reflect.Struct {
		return nil
	}

	out := make([]interface{}, v.NumField())
	for i := 0; i < v.NumField(); i++ {
		out[i] = v.Field(i).Interface()
	}

	return out
}

func Handler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// If no name is provided in the HTTP request body, throw an error
	return muxLambda.Proxy(req)
}

// start the lambda mux router
func main() {
	lambda.Start(Handler)
}
