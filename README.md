# Go, CRUD, AWS
Go, CRUD, AWS is a project I am working on to get a better understanding of Go(Golang), CRUD functionality, and HTTP while hosting it all on AWS. The system works as a storage system for a user's game results in Apex Legends, Fortnite, and Heroes of the Storm.  The user can retrieve their uploaded results and produce a table for each specific game.
## Design
![Architecture Diagram](https://github.com/dwright20/go-crud-aws/blob/master/Images/ArchitectureDiagram.jpg)
System is hosted within an AWS default VPC with 3 EC2 instances, an RDS instance, & DynamoDB  tables.
### Web Server
- Handles all User/Client interaction
- Interacts with Client & App server
- Runs on Go
- Stores all HTML/CSS files
### App Server
- Handles all requests from Web server
- Interacts with Credentials DB (RDS Server) CRUD server, & Web server
- Runs on Go
### CRUD Server
- Handles all requests from App server
- Interacts with App server & Results DB
- Only handles POST & GET requests
- Runs on Go
### Credentials DB
- Database for sign-in credentials
- Runs on PostgreSQL
### Results DB
- Database for game results
- Individual table for each game 
## Key Packages Used
- [AWS](https://github.com/aws/aws-sdk-go) - [Session](https://github.com/aws/aws-sdk-go/aws/session), [DynamoDB](https://github.com/aws/aws-sdk-go/service/dynamodb), [DynamoDB Attribute](https://github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute): Used for all AWS interactions/API calls
- [Creds](https://github.com/dwright20/go-crud-aws/blob/master/Packages/hiddenCreds.go): Package created to store DB credentials and generate them when needed
- [Game](https://github.com/dwright20/go-crud-aws/blob/master/Packages/game.go): Package created to contain structs for different games and a generator for each
- [JSON](https://golang.org/pkg/encoding/json/) & [Bytes](https://golang.org/pkg/bytes/): For passing data between servers 
- [Mux Router](https://github.com/gorilla/mux) & [HTTP](https://golang.org/pkg/net/http/): Powered all HTTP interactions & routing on servers
- [Template](https://golang.org/pkg/html/template/): Used to serve all Web server content
## TODO
- [x] Encrypt Passwords for data store
* Used [bcrypt](https://godoc.org/golang.org/x/crypto/bcrypt#GenerateFromPassword) to salt and hash passwords for storage in DB
* Is not encrypting the password, but solves the issue of storing passwords in plain text
- [ ] Implement a cookie for smoother results viewing
* Added a cookie to be created when user signs in that expires after 24 hours
* Implemented requirement for cookie to be present to access pages past sign-in/creation
- [ ] Error handling & edge cases
- [ ] Setup RR scheme & auto scaling policy for CRUD servers
- [ ] Incorporate more games
## Acknowledgements
Some resources that I found very helpful:
* [Requests/JSON/Forms](http://polyglot.ninja/golang-making-http-requests/)
* [Go on AWS](https://hackernoon.com/deploying-a-go-application-on-aws-ec2-76390c09c2c5)
* [HTML Parsing](https://stackoverflow.com/questions/30109061/golang-parse-html-extract-all-content-with-body-body-tags)
* [Generate HTML content](https://stackoverflow.com/questions/19991124/go-template-html-iteration-to-generate-table-from-struct)
