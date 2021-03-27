# Assignment 2

The API has four endpoints: 
```
/corona/v1/country/
/corona/v1/policy/
/corona/v1/diag/
/corona/v1/notifications/
```
## How to use

### Cases Per Country

#### Request
```
Method: GET
Path: corona/v1/country/{:country_name}{?scope=begin_date-end_date}
```

```{:country_name}``` is a mandatory parameter, and is the name of the country

```{?scope=begin_date-end_date}``` is an optional parameter, with this you can choose a time frame for the information

Example request: ```corona/v1/country/Norway?scope=2020-12-31-2021-01-31
