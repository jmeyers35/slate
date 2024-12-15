package storage

type Team struct {
	ID       string
	TeamCode string
	Name     string
	ESPNID   string
}

type Player struct {
	ID       int
	Name     string
	Position string
	TeamID   int
	ESPNID   string
}
