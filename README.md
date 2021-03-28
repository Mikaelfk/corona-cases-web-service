# Assignment 2

The API has four endpoints: 
```
/corona/v1/country/
/corona/v1/policy/
/corona/v1/notifications/
/corona/v1/diag/
```

## Covid-19 Cases Per Country

This endpoint focuses on the number of confirmed and recovered Covid-19 cases for a given country. The user can also specify a time frame.
When no time frame is specified, the total number of confirmed cases and recovered cases are reported.

### - Request
```
Method: GET
Path: /corona/v1/country/{:country_name}{?scope=begin_date-end_date}
```

```{:country_name}``` is a mandatory parameter, and is the name of the country

```{?scope=begin_date-end_date}``` is an optional parameter, with this you can choose a time frame for the information

Example request: ```corona/v1/country/Norway?scope=2020-12-01-2021-01-31```

## Covid-19 Policy Stringency Trends

This endpoint focuses on the current stringency level of policies regarding Covid-19 for given countries. The endpoint will also report the trend of the stringency level if a time frame is specifed. When no time frame is specifed, only the current stringency level is reported.

### - Request

```
Method: GET
Path: /corona/v1/policy/{:country_name}{?scope=begin_date-end_date}
```

```{:country_name}``` is a mandatory parameter, and is the name of the country

```{?scope=begin_date-end_date}``` is an optional parameter, with this you can choose a time frame for the information

Example request: ```/corona/v1/policy/France?scope=2020-12-01-2021-01-31```

## Diagnostics interface

This endpoint indicates the availability of the services this service depends on. It also reports the version and uptime of the service.

### - Request

```
Method: GET
Path: /corona/v1/diag/
```

## Notification Webhook

Users can register webhooks that are trigged by the service based on specified events related to the stringency of policy or confirmed cases.

### Registrating a Webhook

### - Request

```
Method: POST
Path: /corona/v1/notifications/
```
This request needs a body to be sent with it to register the webhook
This body must contain:
 * The url to be triggered upon event
 * the frequency with which the invocation occurs in seconds (timeout)
 * the information of interest (`stringency` of policy or `confirmed` cases)
 * indication whether notifications should only be sent when information has changed ("ON_CHANGE") - for example, the stringency has changed since the last call and the timeout is reached, or be sent in any case whenever the specified timeout expires ("ON_TIMEOUT")
 * the country for which the trigger applies

 Body (Example):
```
{
   "url": "http://localhost:8081/client/",
   "timeout": 3600,
   "field": "stringency",
   "country": "France",
   "trigger": "ON_CHANGE"
}
```

The request will respond with the id of the webhook. This id can be used to see the details about the webhook or deletion of the webhook.
The format of the id is uuid.

### Deletion of Webhook

### - Request

```
Method: DELETE
Path: /corona/v1/notifications/{id}
```

```{id}``` is the ID of the webhook you wish to delete

### View registered webhook

### - Request

```
Method: GET
Path: /corona/v1/notifications/{id}
```

```{id}``` is the ID of the webhook you wish to get information about.

### View all registered webhooks

### - Request

```
Method: GET
Path: /corona/v1/notifications/
```

### Invoke all webhooks
For testing purposes an endpoint which invokes all the webhooks is included.

### - Request
```
Method: POST
Path: /service
```

## Third Party Libraries Used

Two third party libraries are used in this service:

### uuid library
https://github.com/satori/go.uuid

The reason for using this library is to easily generate random uuid's, if compiling the program on linux, this library would not be needed.
To allow for compiling on Windows and MacOS I have chosen to use this library.

### Mux router
https://github.com/gorilla/mux

This router is used instead of the default go router mainly for simplicity.
