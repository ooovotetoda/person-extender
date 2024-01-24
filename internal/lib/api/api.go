package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type PersonExtends struct {
	Age     int64  `json:"age"`
	Gender  string `json:"gender"`
	Country string `json:"country"`
}

type AgeResponse struct {
	Age int64 `json:"age"`
}

type GenderResponse struct {
	Gender string `json:"gender"`
}

type CountryResponse struct {
	Country []Country `json:"country"`
}

type Country struct {
	CountryId   string  `json:"country_id"`
	Probability float64 `json:"probability"`
}

func GetPersonExtends(name string) (*PersonExtends, error) {
	ageURL := fmt.Sprintf("%s?name=%s", os.Getenv("API_AGIFY_URL"), name)

	res, err := http.Get(ageURL)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var ageResponse AgeResponse
	if err := json.Unmarshal(body, &ageResponse); err != nil {
		return nil, err
	}

	genderURL := fmt.Sprintf("%s?name=%s", os.Getenv("API_GENDERIZE_URL"), name)

	res, err = http.Get(genderURL)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err = io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var genderResponse GenderResponse
	if err := json.Unmarshal(body, &genderResponse); err != nil {
		return nil, err
	}

	countryURL := fmt.Sprintf("%s?name=%s", os.Getenv("API_NATIONALIZE_URL"), name)

	res, err = http.Get(countryURL)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err = io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var countryResponse CountryResponse
	if err := json.Unmarshal(body, &countryResponse); err != nil {
		return nil, err
	}

	personExtends := &PersonExtends{
		Age:     ageResponse.Age,
		Gender:  genderResponse.Gender,
		Country: countryResponse.Country[0].CountryId,
	}

	return personExtends, nil
}
