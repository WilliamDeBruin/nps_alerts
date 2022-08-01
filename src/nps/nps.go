package nps

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"
)

//go:embed state_codes.json
var stateCodesContent []byte

//go:embed parks.json
var parksDetailsContent []byte

const (
	apiKeyEnvKey        = "NPS_API_KEY"
	baseURL             = "https://developer.nps.gov/api/v1/alerts"
	stateCodeFileName   = "nps/state_codes.json"
	parkDetailsFileName = "nps/parks.json"
)

type fetcher struct {
	apiKey     string
	httpClient *http.Client
	stateCodes map[string]string
	parks      *[]parkDetails
}

type Client interface {
	GetAlert(stateCode string) (*alertResponse, error)
	StateCodeToState(stateCode string) (string, error)
	ParkCodeToFullParkName(parkCode string) (string, error)
}

type alertResponse struct {
	Total string `json:"total,omitempty"`
	Limit string `json:"limit,omitempty"`
	Start string `json:"start,omitempty"`
	Data  []struct {
		ID              string `json:"id,omitempty"`
		URL             string `json:"url,omitempty"`
		Title           string `json:"title,omitempty"`
		ParkCode        string `json:"parkCode,omitempty"`
		Description     string `json:"description,omitempty"`
		Category        string `json:"category,omitempty"`
		LastIndexedDate string `json:"lastIndexedDate,omitempty"`
	} `json:"data,omitempty"`
}

type parkDetails struct {
	UnitName        string   `json:"unitName,omitempty"`
	UnitCode        string   `json:"unitCode,omitempty"`
	UnitDesignation string   `json:"unitDesignation,omitempty"`
	State           []string `json:"state,omitempty"`
	EstDate         string   `json:"estDate,omitempty"`
}

func NewClient() (Client, error) {

	apiKey, ok := os.LookupEnv(apiKeyEnvKey)

	if !ok {
		return nil, fmt.Errorf("env key missing: %s", apiKeyEnvKey)
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

func (f *fetcher) GetAlert(stateCode string) (*alertResponse, error) {
	if !f.stateCodeIsValid(stateCode) {
		return nil, fmt.Errorf("state code %s is not a valid state code", stateCode)
	}

	req, err := http.NewRequest("GET", baseURL, nil)

	if err != nil {
		return nil, err
	}

	q := url.Values{}
	q.Add("stateCode", stateCode)

	req.Header.Add("x-api-key", f.apiKey)

	res, err := f.httpClient.Do(req)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	var alertResponse *alertResponse

	err = json.NewDecoder(res.Body).Decode(alertResponse)

	if err != nil {
		return nil, err
	}

	return alertResponse, nil
}

func (f *fetcher) stateCodeIsValid(stateCode string) (ok bool) {
	_, ok = f.stateCodes[stateCode]
	return
}

func (f *fetcher) StateCodeToState(stateCode string) (string, error) {
	if stateName, ok := f.stateCodes[stateCode]; ok {
		return stateName, nil
	}
	return "", fmt.Errorf("cannot find state code %s in list", stateCode)
}

func (f *fetcher) ParkCodeToFullParkName(parkCode string) (string, error) {
	for _, v := range *f.parks {
		if v.UnitCode == parkCode {
			return v.UnitName, nil
		}
	}
	return "", fmt.Errorf("cannot find park code %s in list", parkCode)
}