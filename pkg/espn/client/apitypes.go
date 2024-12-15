package client

import (
	"fmt"
	"strings"
	"time"
)

type TeamsResponse struct {
	Sports []SportTeams `json:"sports"`
}

type SportTeams struct {
	ID      string              `json:"id"`
	UID     string              `json:"uid"`
	Name    string              `json:"name"`
	Slug    string              `json:"slug"`
	Leagues []SportsLeagueTeams `json:"leagues"`
}

type SportsLeagueTeams struct {
	ID           string `json:"id"`
	UID          string `json:"uid"`
	Name         string `json:"name"`
	Abbreviation string `json:"abbreviation"`
	ShortName    string `json:"shortName"`
	Slug         string `json:"slug"`
	Teams        []Team `json:"teams"`
}

type Team struct {
	Team TeamInfo `json:"team"`
}

type TeamInfo struct {
	ID               string `json:"id"`
	UID              string `json:"uid"`
	Slug             string `json:"slug"`
	Abbreviation     string `json:"abbreviation"`
	DisplayName      string `json:"displayName"`
	ShortDisplayName string `json:"shortDisplayName"`
	Name             string `json:"name"`
	Nickname         string `json:"nickname"`
	Location         string `json:"location"`
	Color            string `json:"color"`
}

type RosterResponse struct {
	Timestamp string     `json:"timestamp"`
	Status    string     `json:"status"`
	Season    Season     `json:"season"`
	Athletes  []Athletes `json:"athletes"`
}

type Season struct {
	Year        int    `json:"year"`
	DisplayName string `json:"displayName"`
	// TODO: figure out what this enum means
	Type int    `json:"type"`
	Name string `json:"name"`
}

type Athletes struct {
	Position string    `json:"position"`
	Athletes []Athlete `json:"items"`
}

type Athlete struct {
	ID            string   `json:"id"`
	UID           string   `json:"uid"`
	GUID          string   `json:"guid"`
	FirstName     string   `json:"firstName"`
	LastName      string   `json:"lastName"`
	FullName      string   `json:"fullName"`
	DisplayName   string   `json:"displayName"`
	ShortName     string   `json:"shortName"`
	Weight        float32  `json:"weight"`
	DisplayWeight string   `json:"displayWeight"`
	Height        float32  `json:"height"`
	DisplayHeight string   `json:"displayHeight"`
	Age           int      `json:"age"`
	DateOfBirth   Time     `json:"dateOfBirth"`
	Slug          string   `json:"slug"`
	Jersey        string   `json:"jersey"`
	Position      Position `json:"position"`
	Injuries      []Injury `json:"injuries"`
}

type Position struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	DisplayName  string `json:"displayName"`
	Abbreviation string `json:"abbreviation"`
	// TODO: figure out what this means
	Leaf bool `json:"leaf"`
}

type College struct {
	ID           string `json:"id"`
	Mascot       string `json:"mascot"`
	Name         string `json:"name"`
	ShortName    string `json:"shortName"`
	Abbreviation string `json:"abbreviation"`
	// TODO: add logos when i feel less lazy
}

type Injury struct {
	// Status can be strings like: Out, Injured Reserve, Questionable, etc
	Status string `json:"status"`
	Date   Time   `json:"date"`
}

type Status struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Type         string `json:"type"`
	Abbreviation string `json:"abbreviation"`
}

type Time struct {
	time.Time
}

func (ct *Time) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")

	// Try different formats
	formats := []string{
		"2006-01-02T07:00Z",
		"2006-01-02T15:04Z",
		"2006-01-02T15:04:05Z",
	}

	var lastErr error
	for _, format := range formats {
		t, err := time.Parse(format, s)
		if err == nil {
			ct.Time = t
			return nil
		}
		lastErr = err
	}

	return fmt.Errorf("could not parse time %q with any known format: %v", s, lastErr)
}
