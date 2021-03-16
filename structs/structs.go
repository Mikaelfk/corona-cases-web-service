package structs

// CovidAPIResponse is for storing the information from the response
type CovidAPIResponse struct {
	All EntireCountryInformation
}

// EntireCountryInformation stores all information about a certain country
type EntireCountryInformation struct {
	Confirmed  int
	Recovered  int
	Country    string
	Continent  string
	Population int
	Dates      map[string]int
}

// ReturnConfirmedCases is used for returning a JSON response
type ReturnConfirmedCases struct {
	Country              string  `json:"country"`
	Continent            string  `json:"continent"`
	Scope                string  `json:"scope"`
	Confirmed            int     `json:"confirmed"`
	Recovered            int     `json:"recovered"`
	PopulationPercentage float32 `json:"population_percentage"`
}

// ReturnDiag is used for returning a JSON response for the diag endpoint
type ReturnDiag struct {
	MMediaGroupApi  string `json:"mmediagroupapi"`
	CovidTrackerAPI string `json:"covidtrackerapi"`
	Registered      int    `json:"registered"`
	Version         string `json:"version"`
	Uptime          int    `json:"uptime"`
}
