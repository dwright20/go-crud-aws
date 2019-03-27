# Go, CRUD, AWS
Go, CRUD, AWS is a project I am working on to get a better understanding of Go(Golang), CRUD functionality, and HTTP while hosting it all on AWS. The system works as a storage system for a user's game results in Apex Legends, Fortnite, and Heroes of the Storm.  The user can retrieve their uploaded results and produce a table for each specific game.
## Design
### [Design 1](https://github.com/dwright20/go-crud-aws/blob/master/Images/ArchitectureDiagram.jpg)
![Architecture Diagram](https://github.com/dwright20/go-crud-aws/blob/master/Images/ArchitectureDiagram2.jpeg)

System is hosted within on AWS  with 3 EC2 instances, an API gateway, various Lambda functions, an RDS instance, & DynamoDB  tables.
### Web Server
- Handles all User/Client interaction
- Interacts with Client & App server
- Runs on Go
- Stores all HTML/CSS files
- Handles cookie
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
### Back-up Gateway
- AWS API Gateway
- Proxies all requests to appropriate Lambda function
- If request comes from Web Server, can directly call CRUD functions if needed
### App Functions
- Same functionality as App Server
- Runs on Go
### CRUD Functions
- Same functionality as CRUD Server
- Runs on GO
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
- [x] Implement a cookie for smoother results viewing
* Added a cookie that is created when user signs in that expires after 24 hours
* Cookie required to access pages past sign-in/creation
* Streamlined process for viewing game results by leveraging a cookie
* Now skips a webpage that requests user's username to retrieve results
* Ensures user can only see their own results
- [x] Setup fail-over 
* Setup fail-over to be done at the server level by the go applications
* Server will check if primary server is up, if it is not, it will send request to the fail-over API Gateway backed by Lambda
* If CRUD Server is down, requests will still go to App Server prior to fail-over gateway
* If App Server is down, all requests will go to fail-over gateway and will not reach CRUD server even if it is up
- [ ] Setup RR scheme & auto scaling policy for CRUD servers
- [ ] Error handling & edge cases
- [ ] Incorporate more games
## Acknowledgements
Some resources that I found very helpful:
* [Requests/JSON/Forms](http://polyglot.ninja/golang-making-http-requests/)
* [Go on AWS](https://hackernoon.com/deploying-a-go-application-on-aws-ec2-76390c09c2c5)
* [HTML Parsing](https://stackoverflow.com/questions/30109061/golang-parse-html-extract-all-content-with-body-body-tags)
* [Generate HTML content](https://stackoverflow.com/questions/19991124/go-template-html-iteration-to-generate-table-from-struct)
