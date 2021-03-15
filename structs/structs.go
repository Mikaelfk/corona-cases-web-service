package structs

type CovidApiResponse struct {
	All EntireCountryInformation
}

type EntireCountryInformation struct {
	Confirmed  int
	Recovered  int
	Country    string
	Continent  string
	Population int
}

type ReturnConfirmedCases struct {
	Country              string  `json:"country"`
	Continent            string  `json:"continent"`
	Scope                string  `json:"scope"`
	Confirmed            int     `json:"confirmed"`
	Recovered            int     `json:"recovered"`
	PopulationPercentage float32 `json:"population_percentage"`
}
