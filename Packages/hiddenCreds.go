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

	creds.Host     = //Host
	creds.Port     = //Port
	creds.User     = //User
	creds.Password = //Password
	creds.Dbname   = //DB Type

	return creds
}

