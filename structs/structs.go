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

// StringencyResponse stores the information from the response
type StringencyResponse struct {
	StringencyData StringencyData `json:"stringencyData"`
}

// StringencyData stores the stringency_actual value
type StringencyData struct {
	StringencyActual float32 `json:"stringency_actual"`
	Msg              string  `json:"msg"`
}

// ReturnStringency is used for returning a JSON response for the policy endpoint
type ReturnStringency struct {
	Country    string  `json:"country"`
	Scope      string  `json:"scope"`
	Stringency float32 `json:"stringency"`
	Trend      float32 `json:"trend"`
}

// ReturnDiag is used for returning a JSON response for the diag endpoint
type ReturnDiag struct {
	MMediaGroupApi  string `json:"mmediagroupapi"`
	CovidTrackerAPI string `json:"covidtrackerapi"`
	CountryAPI      string `json:"countryapi"`
	Registered      int    `json:"registered"`
	Version         string `json:"version"`
	Uptime          int    `json:"uptime"`
}

type CountryResponse struct {
	Name       string `json:"name"`
	Alpha3Code string `json:"alpha3Code"`
}

type WebhookRegistration struct {
	ID      string
	Url     string `json:"url"`
	Timeout int    `json:"timeout"`
	Field   string `json:"field"`
	Country string `json:"country"`
	Trigger string `json:"trigger"`
}
