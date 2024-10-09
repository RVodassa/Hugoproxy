package server

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

func HandleGeocode(w http.ResponseWriter, r *http.Request) {
	apiKey, secretKey := os.Getenv("ApiKey"), os.Getenv("SecretKey")

	geoService := NewGeoService(apiKey, secretKey)

	var data GeocodeRequest

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		log.Printf("ERROR: Не удалось прочитать JSON данные: %v\n", err)
		http.Error(w, "Неверные данные", http.StatusBadRequest)
		return
	}

	log.Printf("INFO: Запрос на геокодирование: lat=%s, lng=%s\n", data.Lat, data.Lng)

	geoRes, err := geoService.GeoCode(data.Lat, data.Lng)
	if err != nil {
		log.Println(err)
	}

	res, err := json.MarshalIndent(geoRes, "", " ")
	if err != nil {
		log.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func HandleSearch(w http.ResponseWriter, r *http.Request) {

	apiKey, secretKey := os.Getenv("ApiKey"), os.Getenv("SecretKey")

	geoService := NewGeoService(apiKey, secretKey)

	var data SearchRequest

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		log.Printf("ERROR: Не удалось прочитать JSON данные: %v\n", err)
		http.Error(w, "Неверные данные", http.StatusBadRequest)
		return
	}

	log.Printf("INFO: Запрос на поиск: query=%s", data.Query)

	geoRes, err := geoService.AddressSearch(data.Query)
	if err != nil {
		log.Println(err)
	}

	res, err := json.MarshalIndent(geoRes, "", " ")
	if err != nil {
		log.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}
