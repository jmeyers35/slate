package storage

type Team struct {
	ID       string
	TeamCode string
	Name     string
	ESPNID   string
}

type Player struct {
	ID       string
	Name     string
	Position string
	TeamID   string
	ESPNID   string
}
