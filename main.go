package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

const infoUrl = "https://data.gov.sg/api/action/datastore_search?resource_id=139a3035-e624-4f56-b63f-89ae28d4ae4c&q="
const availUrl = "https://api.data.gov.sg/v1/transport/carpark-availability?date_time="

func main() {
	port := os.Getenv("PORT")
	mux := http.NewServeMux()
	mux.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.Dir("./static"))))
	mux.Handle("/.well-known/ai-plugin.json", cors(manifest))
	mux.Handle("/openapi.yaml", cors(specs))
	mux.Handle("/info", cors(info))
	mux.Handle("/availability", cors(availability))

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
	server.ListenAndServe()
}

// middleware to set CORS for handlers
func cors(h http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET")
		w.Header().Set("Access-Control-Allow-Headers", "Origin,Content-Type,OpenAI-Conversation-ID,OpenAI-Ephemeral-User-ID")
		h.ServeHTTP(w, r)
	})
}

// returns openai.yaml file
func specs(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "openapi.yaml")
}

// returns .well-known/ai-plugin.json file
func manifest(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, ".well-known/ai-plugin.json")
}

// handler function
// get information on car park
func info(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	records := getCarParkRecords(query)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(records)
}

// handler function
// get availability on car park
func availability(w http.ResponseWriter, r *http.Request) {
	// get the car parks requested
	query := r.URL.Query().Get("q")
	log.Println("q:", query)
	log.Println("query:", query[:len(query)-1])
	records := getCarParkRecords(query)

	// get the availability of car parks in Singapore now
	now := time.Now()
	dt := now.Format(time.RFC3339)
	availability := new(CarParkAvailability)
	carparkdata := []CarParkData{}

	resp, err := http.Get(availUrl + url.QueryEscape(dt))
	log.Println("url:", availUrl+url.QueryEscape(dt))
	if err != nil {
		log.Println("Cannot call url", availUrl+dt, err)
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	log.Println("data", string(data))
	if err != nil {
		log.Println("Cannot read JSON file", err)
	}
	err = json.Unmarshal(data, availability)
	if err != nil {
		log.Println("Error unmarshalling json data:", err)
	}
	log.Println("availability:", availability)
	// filter off those that we want to look for
	for _, item := range availability.Items {
		for _, data := range item.Data {
			for _, record := range records {
				if data.CarParkNumber == record.CarParkNo {
					log.Println("data:", data)
					carparkdata = append(carparkdata, data)
				}
			}
		}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(carparkdata)
}

// get card park records based on a query
func getCarParkRecords(query string) []CarParkRecord {
	carparks := new(CarParks)
	resp, err := http.Get(infoUrl + query)
	if err != nil {
		log.Println("Cannot call url", infoUrl+query, err)
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Cannot read JSON file", err)
	}
	err = json.Unmarshal(data, carparks)
	if err != nil {
		log.Println("Error unmarshalling json data:", err)
	}
	// log.Println("records:", carparks.Result.Records)
	return carparks.Result.Records
}
