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
- Data store for sign-in credentials
- Runs on PostgreSQL
### Results DB
- Data store for game results
- Individual table for each game 
## Packages Used
### AWS
### HTTP
## TODO
- [ ] Encrypt Passwords for data store
- [ ] Implement a cookie for smoother results viewing
- [ ] Setup RR scheme & auto scaling policy for CRUD servers
- [ ] Incorporate more games

