package client

type ScheduleResponse struct {
	Events  []Event       `json:"events"`
	Week    WeekDetail    `json:"week"`
	Season  Season        `json:"season"`  // Using original Season type which has Type as int
	Leagues []LeagueDetail `json:"leagues"`
}

type LeagueDetail struct {
	ID            string          `json:"id"`
	UID           string          `json:"uid"`
	Name          string          `json:"name"`
	Abbreviation  string          `json:"abbreviation"`
	Season        LeagueSeason    `json:"season"`
	CalendarType  string          `json:"calendarType"`
	Calendar      []CalendarEntry `json:"calendar"`
}

type LeagueSeason struct {
	Year        int               `json:"year"`
	Type        ScheduleSeasonType `json:"type"`
	StartDate   string            `json:"startDate"`
	EndDate     string            `json:"endDate"`
}

type ScheduleSeasonType struct {
	ID           string `json:"id"`
	Type         int    `json:"type"`
	Name         string `json:"name"`
	Abbreviation string `json:"abbreviation"`
}

type CalendarEntry struct {
	Label      string               `json:"label"`
	Value      string               `json:"value"`
	StartDate  string               `json:"startDate"`
	EndDate    string               `json:"endDate"`
	Entries    []CalendarDetailEntry `json:"entries,omitempty"`
}

type CalendarDetailEntry struct {
	Label          string `json:"label"`
	AlternateLabel string `json:"alternateLabel"`
	Detail         string `json:"detail"`
	Value          string `json:"value"`
	StartDate      string `json:"startDate"`
	EndDate        string `json:"endDate"`
}

type Event struct {
	ID           string        `json:"id"`
	Date         Time          `json:"date"`
	Name         string        `json:"name"`
	ShortName    string        `json:"shortName"`
	Status       EventStatus   `json:"status"`  // Changed to EventStatus
	Venue        Venue         `json:"venue"`
	Competitions []Competition `json:"competitions"`
}

type Competition struct {
	ID          string       `json:"id"`
	Date        Time         `json:"date"`
	Competitors []Competitor `json:"competitors"`
	Venue       Venue        `json:"venue"`
	Status      EventStatus  `json:"status"`  // Changed to EventStatus
}

type EventStatus struct {
	Clock        float64     `json:"clock"`
	DisplayClock string      `json:"displayClock"`
	Period      int         `json:"period"`
	Type        StatusType   `json:"type"`
}

type StatusType struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	State        string `json:"state"`
	Completed    bool   `json:"completed"`
	Description  string `json:"description"`
	Detail       string `json:"detail"`
	ShortDetail  string `json:"shortDetail"`
}

type Competitor struct {
	ID       string   `json:"id"`
	HomeAway string   `json:"homeAway"`
	Team     TeamInfo `json:"team"`
	Score    string   `json:"score"`
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