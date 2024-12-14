package client

type Sport string

const (
	SportFootball Sport = "football"
)

type League string

const (
	LeagueNFL League = "nfl"
)

type ClientConfiguration struct {
	Sport  Sport
	League League
}
