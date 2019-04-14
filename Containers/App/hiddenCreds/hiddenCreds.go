package hiddenCreds

type Creds struct {
	Host string
	Port int
	User string
	Password string
	Dbname string
}

func GetCreds () Creds{
	creds := Creds{}

	creds.Host     = "dbinstance1.com7aocj6hq1.us-east-1.rds.amazonaws.com"
	creds.Port     = 5432
	creds.User     = "dbinstance1admin"
	creds.Password = "_dbinstance1Admin_"
	creds.Dbname   = "postgres"

	return creds
}

