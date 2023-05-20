package main

import "time"

type CarParkData struct {
	CarParkInfo []struct {
		TotalLots     string `json:"total_lots"`
		LotType       string `json:"lot_type"`
		LotsAvailable string `json:"lots_available"`
	} `json:"carpark_info"`
	CarParkNumber  string `json:"carpark_number"`
	UpdateDatetime string `json:"update_datetime"`
}

type CarParkAvailability struct {
	Items []struct {
		Timestamp time.Time     `json:"timestamp"`
		Data      []CarParkData `json:"carpark_data"`
	} `json:"items"`
}

type CarParkRecord struct {
	FullCount           string  `json:"_full_count"`
	ShortTermParking    string  `json:"short_term_parking"`
	CarParkType         string  `json:"car_park_type"`
	YCoord              string  `json:"y_coord"`
	XCoord              string  `json:"x_coord"`
	Rank                float64 `json:"rank"`
	FreeParking         string  `json:"free_parking"`
	GantryHeight        string  `json:"gantry_height"`
	CarParkBasement     string  `json:"car_park_basement"`
	NightParking        string  `json:"night_parking"`
	Address             string  `json:"address"`
	CarParkDecks        string  `json:"car_park_decks"`
	ID                  int     `json:"_id"`
	CarParkNo           string  `json:"car_park_no"`
	TypeOfParkingSystem string  `json:"type_of_parking_system"`
}

type CarParks struct {
	Help    string `json:"help"`
	Success bool   `json:"success"`
	Result  struct {
		ResourceID string `json:"resource_id"`
		Fields     []struct {
			Type string `json:"type"`
			ID   string `json:"id"`
		} `json:"fields"`
		Q       string          `json:"q"`
		Records []CarParkRecord `json:"records"`
		Links   struct {
			Start string `json:"start"`
			Next  string `json:"next"`
		} `json:"_links"`
		Total int `json:"total"`
	} `json:"result"`
}
