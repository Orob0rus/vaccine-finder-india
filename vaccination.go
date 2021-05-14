package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Centers struct {
	Centers []Center `json:"centers,omitempty"`
}

type Session struct {
	SessionID         string   `json:"session_id,omitempty"`
	Date              string   `json:"date,omitempty"`
	AvailableCapacity int      `json:"available_capacity,omitempty"`
	MinAgeLimit       int      `json:"min_age_limit,omitempty"`
	Vaccine           string   `json:"vaccine,omitempty"`
	Slots             []string `json:"slots,omitempty"`
}

type Center struct {
	CenterID     int           `json:"center_id,omitempty"`
	Name         string        `json:"name,omitempty"`
	Address      string        `json:"address,omitempty"`
	StateName    string        `json:"state_name,omitempty"`
	DistrictName string        `json:"district_name,omitempty"`
	BlockName    string        `json:"block_name,omitempty"`
	Pincode      int           `json:"pincode,omitempty"`
	Latitude     float64       `json:"lat,omitempty"`
	Longitude    float64       `json:"long,omitempty"`
	FeeType      string        `json:"fee_type,omitempty"`
	Sessions     []Session     `json:"sessions,omitempty"`
	To           string        `json:"to,omitempty"`
	From         string        `json:"from,omitempty"`
	VaccineFees  []VaccineFees `json:"vaccine_fees,omitempty"`
}

type VaccineFees struct {
	Vaccine string `json:"vaccine,omitempty"`
	Fee     string `json:"fee,omitempty"`
}

type CenterResult struct {
	Name              string
	Address           string
	BlockName         string
	DistrictName      string
	StateName         string
	Pincode           int
	To                string
	From              string
	Latitude          float64
	Longitude         float64
	Vaccine           string
	AvailableCapacity int
	Date              string
}

type Result struct {
	results []CenterResult
}

func findVaccinationSlots(date string, district_id, age int) {
	results := Result{[]CenterResult{}}
	client := &http.Client{}
	req, _ := http.NewRequest("GET", fmt.Sprintf("https://cdn-api.co-vin.in/api/v2/appointment/sessions/public/calendarByDistrict?district_id=%d&date=%s", district_id, date), nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:88.0) Gecko/20100101 Firefox/88.0")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("%s", err.Error())
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("%s", err.Error())
		return
	}
	var centers = new(Centers)
	err = json.Unmarshal(body, centers)
	if err != nil {
		fmt.Printf("%s", err.Error())
		return
	}
	for _, center := range centers.Centers {
		for _, session := range center.Sessions {
			if session.MinAgeLimit <= age {
				centerResult := CenterResult{
					Name:              center.Name,
					Address:           center.Address,
					BlockName:         center.BlockName,
					DistrictName:      center.DistrictName,
					StateName:         center.StateName,
					Pincode:           center.Pincode,
					To:                center.To,
					From:              center.From,
					Vaccine:           session.Vaccine,
					AvailableCapacity: session.AvailableCapacity,
					Date:              session.Date,
				}
				results.results = append(results.results, centerResult)
			}
		}
	}
	for _, result := range results.results {
		fmt.Printf("%#v\n", result)
		fmt.Println("------------------------------------")
	}
}

func main() {
	findVaccinationSlots("14-05-2021", 506, 24)
	return
}
