// CRUD Server read functionality lambda
package main

import (
	"game"
	"bytes"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/awslabs/aws-lambda-go-api-proxy/gorillamux"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
	"reflect"
)

var muxLambda *gorillamux.GorillaMuxAdapter  // initialize mux lambda adapter

func init() {
	log.Printf("Create game Mux start...")
	r := mux.NewRouter()
	r.HandleFunc("/view/{game}/{user}", readEntry).Methods("GET")  // if request is from web tier
	r.HandleFunc("/read/{game}/{user}", readEntry).Methods("GET")  // if request is from app tier
	muxLambda = gorillamux.New(r)
}

// read requesting users game results, generate html table
// of the results, and encode results and send them in
// response body
func readEntry(w http.ResponseWriter, r *http.Request)  {
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

// pass request into mux proxy
func Handler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return muxLambda.Proxy(req)
}

// start the lambda mux router
func main() {
	lambda.Start(Handler)
}
