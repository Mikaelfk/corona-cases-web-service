package utils

import "net/http"

const CovidCasesAPI = "https://covid-api.mmediagroup.fr/v1"
const DataAPI = "https://covidtrackerapi.bsg.ox.ac.uk/api/v2"
const CountryAPI = "https://restcountries.eu/rest/v2"

const BadRequest = http.StatusBadRequest
const InternalServerError = http.StatusInternalServerError
