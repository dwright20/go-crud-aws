// CRUD Server create functionality in lambda
package main

import (
	"game"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/awslabs/aws-lambda-go-api-proxy/gorillamux"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

var muxLambda *gorillamux.GorillaMuxAdapter  // initialize mux lambda adapter

func init() {
	log.Printf("Create game Mux start...")
	r := mux.NewRouter()
	r.HandleFunc("/create/{game}", createEntryApp).Methods("POST")  // if request is from app tier
	r.HandleFunc("/submit/{user}", createEntry).Methods("POST")  // if request is from web tier
	muxLambda = gorillamux.New(r)
}

// create appropriate game result entry and upload to db
// for requests from app tier
func createEntryApp(_ http.ResponseWriter, r *http.Request)  {
	params := mux.Vars(r)

	log.Print("starting session...")

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
	})

	if err != nil{
		log.Print(err.Error())
	}

	// create dynamodb client
	svc := dynamodb.New(sess)

	if params["game"] == "apex" {
		var game game.Apex
		_ = json.NewDecoder(r.Body).Decode(&game) //decode request contents into game

		log.Println(game) //log newly created game

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

	} else if params["game"] == "fort" {
		var game game.Fort
		_ = json.NewDecoder(r.Body).Decode(&game)

		log.Println(game)

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

	} else if params["game"] == "hots"{
		var game game.Hots
		_ = json.NewDecoder(r.Body).Decode(&game)

		log.Println(game)

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
}

// create appropriate game result entry and upload to db
// for requests from web tier
func createEntry(w http.ResponseWriter, r *http.Request)  {
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

// pass request into mux proxy
func Handler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return muxLambda.Proxy(req)
}

// start the lambda mux router
func main() {
	lambda.Start(Handler)
}
