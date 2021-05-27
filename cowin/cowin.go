package cowin

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type VaccineSlotLocator struct {
	ctx                    context.Context
	date                   string
	pincodes               []string
	endpoint               string
	minAgeLimit            int
	minAvailableCapacity   int
	availableCapacityDose1 int
	availableCapacityDose2 int
	ResultChan             chan LocatorResult
}

type LocatorResult struct {
	Pincode                int
	CenterName             string
	Date                   string
	AvailableCapacity      int
	MinAgeLimit            int
	Vaccine                string
	AvailableCapacityDose1 int
	AvailableCapacityDose2 int
}

type CowinVacSearchResult struct {
	Centers []struct {
		CenterID     int    `json:"center_id"`
		Name         string `json:"name"`
		Address      string `json:"address"`
		StateName    string `json:"state_name"`
		DistrictName string `json:"district_name"`
		BlockName    string `json:"block_name"`
		Pincode      int    `json:"pincode"`
		Lat          int    `json:"lat"`
		Long         int    `json:"long"`
		From         string `json:"from"`
		To           string `json:"to"`
		FeeType      string `json:"fee_type"`
		Sessions     []struct {
			SessionID              string   `json:"session_id"`
			Date                   string   `json:"date"`
			AvailableCapacity      int      `json:"available_capacity"`
			MinAgeLimit            int      `json:"min_age_limit"`
			Vaccine                string   `json:"vaccine"`
			Slots                  []string `json:"slots"`
			AvailableCapacityDose1 int      `json:"available_capacity_dose1"`
			AvailableCapacityDose2 int      `json:"available_capacity_dose2"`
		} `json:"sessions"`
		VaccineFees []struct {
			Vaccine string `json:"vaccine"`
			Fee     string `json:"fee"`
		} `json:"vaccine_fees"`
	} `json:"centers"`
}

func NewVacSlotLocator(ctx context.Context, date string, pincodes []string, minAgeLimit, minAvailableCapacity, availableCapacityDose1, availableCapacityDose2 int) *VaccineSlotLocator {
	resultChan := make(chan LocatorResult)
	endpoint := "https://cdn-api.co-vin.in/api/v2/appointment/sessions/public/calendarByPin"
	return &VaccineSlotLocator{
		endpoint:               endpoint,
		pincodes:               pincodes,
		date:                   date,
		minAgeLimit:            minAgeLimit,
		minAvailableCapacity:   minAvailableCapacity,
		availableCapacityDose1: availableCapacityDose1,
		availableCapacityDose2: availableCapacityDose2,
		ResultChan:             resultChan,
		ctx:                    ctx,
	}
}

func (v VaccineSlotLocator) Locate(interval int) {
	for {
		select {
		case <-v.ctx.Done():
			return
		default:
			for _, pincode := range v.pincodes {
				go v.executeLocator(pincode)
			}
			log.Printf("Locating Vaccination slots in %d Seconds", interval)
			time.Sleep(time.Duration(interval) * time.Second)
		}
	}
}

func (v VaccineSlotLocator) executeLocator(pincode string) {
	result, err := v.getVacSlotResultFromCowin(pincode)
	if err != nil {
		log.Println("Error : ", err.Error())
	}
	for _, center := range result.Centers {
		for _, session := range center.Sessions {
			if session.MinAgeLimit <= v.minAgeLimit && session.AvailableCapacity >= v.minAvailableCapacity && (session.AvailableCapacityDose1 >= v.availableCapacityDose1 || session.AvailableCapacityDose2 >= v.availableCapacityDose2) {
				m := LocatorResult{
					CenterName:             center.Name,
					Pincode:                center.Pincode,
					Date:                   session.Date,
					AvailableCapacity:      session.AvailableCapacity,
					MinAgeLimit:            session.MinAgeLimit,
					Vaccine:                session.Vaccine,
					AvailableCapacityDose1: session.AvailableCapacityDose1,
					AvailableCapacityDose2: session.AvailableCapacityDose2,
				}
				v.ResultChan <- m
			}
		}
	}
}

func (v VaccineSlotLocator) getVacSlotResultFromCowin(pincode string) (CowinVacSearchResult, error) {
	url := fmt.Sprintf("%s?pincode=%s&date=%s", v.endpoint, pincode, v.date)
	method := "GET"
	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	req.Header.Set("User-Agent", "PostmanRuntime/7.28.0")
	var result CowinVacSearchResult
	if err != nil {
		return result, err
	}
	res, err := client.Do(req)
	if err != nil {
		return result, err
	}
	if res.StatusCode != http.StatusOK {
		log.Println("Status not ok ", res.StatusCode)
		return result, fmt.Errorf("status not ok")
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return result, err
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return result, err
	}
	return result, nil
}
