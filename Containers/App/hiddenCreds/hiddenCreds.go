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

	creds.Host     = ""
	creds.Port     = 
	creds.User     = ""
	creds.Password = ""
	creds.Dbname   = ""

	return creds
}

