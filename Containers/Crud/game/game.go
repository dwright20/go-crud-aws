package game

type Apex struct {
	Username	string	`json:"username"`
	Date		string	`json:"date"`
	Game		string	`json:"game"`
	Result		string	`json:"result"`
	Legend		string	`json:"legend"`
	Kills		string	`json:"kills"`
	Placement	string	`json:"placement"`
	Damage		string	`json:"damage"`
	Time 		string	`json:"time"`
	Teammates	string	`json:"teammates"`
}

type Fort struct {
	Username	string	`json:"username"`
	Date		string	`json:"date"`
	Game		string	`json:"game"`
	Result		string	`json:"result"`
	Kills		string	`json:"kills"`
	Placement	string	`json:"placement"`
	Gamemode	string	`json:"mode"`
	Teammates	string	`json:"teammates"`
}

type Hots struct {
	Username	string	`json:"username"`
	Date		string	`json:"date"`
	Game		string	`json:"game"`
	Result		string	`json:"result"`
	Hero		string	`json:"hero"`
	Kills		string	`json:"kills"`
	Deaths		string	`json:"deaths"`
	Assists		string	`json:"assists"`
	Time		string	`json:"time"`
	Map			string	`json:"map"`
}

func NewApex (username, date, game, result, legend, kills, placement, damage, time, teammates string) Apex {
	apex := Apex{username, date, game, result, legend, kills, placement, damage, time, teammates}
	return apex
}

func NewFort (username, date, game, result, kills, placement, mode, teammates string) Fort {
	fort := Fort{username, date, game, result, kills, placement, mode, teammates}
	return fort
}

func NewHots (username, date, game, result, hero, kills, deaths, assists, time, Map string) Hots {
	hots := Hots{username, date, game, result, hero, kills, deaths, assists, time, Map}
	return hots
}