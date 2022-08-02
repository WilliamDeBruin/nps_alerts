# nps_alerts

A simple app to use the Twilio SDK to request and send alerts from the National Parks Service related to national parks.

## Local Development 

### Local Environment Setup
`Make run` injects environment variables at runtime using `.env`. See [sample.env](./sample.env) for required variables.

### Run Locally 

Run `make build` to build the docker image and tag it `nps-alerts`

Run `make run` to run the container locally & expose port `8080`

Use the following cURL commands to simulate incoming SMS messages:

> Note: the `from` phone number must be verified via before it can be used as the recepient of an SMS for a trial account

#### Help text 

Use the following cURL to simulate a help sms incoming
```sh
curl --location --request POST 'localhost:8080/incoming-sms' \
    --header 'Content-Type: application/x-www-form-urlencoded' \
    --data-urlencode 'body=help' \
    --data-urlencode 'from=+12407439754'
```

#### Send Alert Message

Use the following cURL to simulate an alert sms incoming 
```sh
curl --location --request POST 'localhost:8080/incoming-sms' \
    --header 'Content-Type: application/x-www-form-urlencoded' \
    --data-urlencode 'body=alerts UT' \
    --data-urlencode 'from=+12407439754'
```

#### Get Application Health
```sh
curl --location --request GET 'localhost:8080/health'
```

## Usage

The app is interacted with via SMS. The following functionality is available:

### Help

Users can text `"help"` to receive help text related to app usage

#### Example
```
> Help!

> Welcome to NPS alerts! Here is a list of commands:
> Help: receive this help text
> Alerts {state}: Text "alerts" followed by the 2-letter state code of the state you would like to see alerts for
```
### alerts {state}

Users can text `"alerts {state}"` where `{state}` is a valid 2-letter state code to see the most recent alert for NPS parks in that state.

#### Example

```
> Alerts CA

> Here is the most recent NPS California alert from Alcatraz Island, released 2022-06-07 17:55:48.0:
> 
>  Face masks are required indoors
> Masks are required for everyone in all NPS buildings and enclosed public transportation, regardless of vaccination status.
>
> For a full list of NPS California alerts, visit https://www.nps.gov/planyourvisit/alerts.htm?s=CA&p=1&v=0
```
