package nps

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

//go:embed state_codes.json
var stateCodesContent []byte

//go:embed parks.json
var parksDetailsContent []byte

const (
	baseURL = "https://developer.nps.gov/api/v1/alerts"

	alertsUrl = "https://www.nps.gov/planyourvisit/alerts.htm?s=%s&p=1&v=0"
)

type fetcher struct {
	apiKey     string
	httpClient *http.Client
	stateCodes map[string]string
	parks      *[]parkDetails
}

type Client interface {
	GetAlert(stateCode string) (*AlertDetails, error)
	SetTransport(http.RoundTripper)
}

type alertResponse struct {
	Total string     `json:"total,omitempty"`
	Limit string     `json:"limit,omitempty"`
	Start string     `json:"start,omitempty"`
	Data  []npsAlert `json:"data,omitempty"`
}

type npsAlert struct {
	ID              string `json:"id,omitempty"`
	URL             string `json:"url,omitempty"`
	Title           string `json:"title,omitempty"`
	ParkCode        string `json:"parkCode,omitempty"`
	Description     string `json:"description,omitempty"`
	Category        string `json:"category,omitempty"`
	LastIndexedDate string `json:"lastIndexedDate,omitempty"`
}

type parkDetails struct {
	UnitName        string   `json:"unitName,omitempty"`
	UnitCode        string   `json:"unitCode,omitempty"`
	UnitDesignation string   `json:"unitDesignation,omitempty"`
	State           []string `json:"state,omitempty"`
	EstDate         string   `json:"estDate,omitempty"`
}

type AlertDetails struct {
	FullStateName   string
	FullParkName    string
	RecentAlertDate string
	AlertHeader     string
	AlertMessage    string
	URL             string
}

func NewClient(apiKey string) (Client, error) {

	if apiKey == "" {
		return nil, fmt.Errorf("apiKey cannot be empty")
	}

	c := &http.Client{
		Timeout: time.Duration(1) * time.Second,
	}

	var stateCodes map[string]string
	err := json.Unmarshal(stateCodesContent, &stateCodes)
	if err != nil {
		return nil, err
	}

	var parksDetails *[]parkDetails = &[]parkDetails{}
	err = json.Unmarshal(parksDetailsContent, parksDetails)
	if err != nil {
		return nil, err
	}

	return &fetcher{
		apiKey:     apiKey,
		httpClient: c,
		stateCodes: stateCodes,
		parks:      parksDetails,
	}, nil
}

func (f *fetcher) GetAlert(stateCode string) (*AlertDetails, error) {

	fullStateName, err := f.stateCodeToState(strings.ToUpper(stateCode))

	if err != nil {
		return nil, fmt.Errorf("state code %s is not a valid state code", stateCode)
	}

	req, _ := http.NewRequest("GET", baseURL, nil)

	q := url.Values{}
	q.Add("stateCode", stateCode)
	req.URL.RawQuery = q.Encode()

	req.Header.Add("x-api-key", f.apiKey)

	res, err := f.httpClient.Do(req)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	alertResponse := &alertResponse{}

	err = json.NewDecoder(res.Body).Decode(alertResponse)

	if err != nil {
		return nil, err
	}

	fullParkName, err := f.parkCodeToFullParkName(alertResponse.Data[0].ParkCode)

	if err != nil {
		return nil, fmt.Errorf("cannot find details for park code %s", alertResponse.Data[0].ParkCode)
	}

	return &AlertDetails{
		FullStateName:   fullStateName,
		FullParkName:    fullParkName,
		RecentAlertDate: alertResponse.Data[0].LastIndexedDate,
		AlertHeader:     alertResponse.Data[0].Title,
		AlertMessage:    alertResponse.Data[0].Description,
		URL:             fmt.Sprintf(alertsUrl, stateCode),
	}, nil
}

func (f *fetcher) stateCodeToState(stateCode string) (string, error) {
	if stateName, ok := f.stateCodes[stateCode]; ok {
		return stateName, nil
	}
	return "", fmt.Errorf("cannot find state code %s in list", stateCode)
}

func (f *fetcher) parkCodeToFullParkName(parkCode string) (string, error) {
	for _, v := range *f.parks {
		if v.UnitCode == parkCode {
			return v.UnitName, nil
		}
	}
	return "", fmt.Errorf("cannot find park code %s in list", parkCode)
}

func (f *fetcher) SetTransport(transport http.RoundTripper) {
	f.httpClient.Transport = transport
}
