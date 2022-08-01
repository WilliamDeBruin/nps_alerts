# nps_alerts

A simple app to use the Twilio SDK to request and send alerts from the National Parks Service related to national parks.

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

> There are 0 new alerts today for California parks. Here is the most recent NPS California alert from Alcatraz Island, released 2022-06-07 17:55:48.0:
> 
>  Face masks are required indoors
> Masks are required for everyone in all NPS buildings and enclosed public transportation, regardless of vaccination status.
>
> For a full list of NPS California alerts, visit https://www.nps.gov/planyourvisit/alerts.htm?s=CA&p=1&v=0
```

# Local Development 

## Setup
