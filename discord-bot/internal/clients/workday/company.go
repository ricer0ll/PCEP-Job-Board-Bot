package workday

import (
	"encoding/json"
	"errors"
	"os"
)

type Company struct {
	Name                           string
	WorkdayBaseURL                 string
	WorkdayRequestURL              string
	JobFamily                      []string
	JobFamilyGroup                 []string
	Locations                      []string
	LocationCountry                []string
	Location_Country               []string
	Location_Region_State_Province []string
}

type Config struct {
	Companies []Company `json:"companies"`
}

func loadCompanies(path string) ([]Company, error) {
	fileData, err := os.ReadFile(path)
	if err != nil {
		return []Company{}, errors.New("unable to read json file")
	}

	var config Config
	err = json.Unmarshal(fileData, &config)
	if err != nil {
		return []Company{}, errors.New("unable to unmarsal json file.")
	}

	return config.Companies, nil
}
