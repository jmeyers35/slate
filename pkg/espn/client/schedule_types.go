package client

type ScheduleResponse struct {
	Events    []Event    `json:"events"`
	Week      WeekDetail `json:"week"`
	Season    Season     `json:"season"`
	Leagues   []League   `json:"leagues"`
}

type Event struct {
	ID        string    `json:"id"`
	Date      Time      `json:"date"`
	Name      string    `json:"name"`
	ShortName string    `json:"shortName"`
	Status    Status    `json:"status"`
	Venue     Venue     `json:"venue"`
	Competitions []Competition `json:"competitions"`
}

type Competition struct {
	ID          string         `json:"id"`
	Date        Time           `json:"date"`
	Competitors []Competitor   `json:"competitors"`
	Venue       Venue         `json:"venue"`
	Status      Status        `json:"status"`
}

type Competitor struct {
	ID         string    `json:"id"`
	HomeAway   string    `json:"homeAway"`
	Team       TeamInfo  `json:"team"`
	Score      string    `json:"score"`
}

type Venue struct {
	ID       string `json:"id"`
	FullName string `json:"fullName"`
	Indoor   bool   `json:"indoor"`
}

type WeekDetail struct {
	Number int    `json:"number"`
	Type   string `json:"type"` // e.g., "REG", "PRE", "POST"
}