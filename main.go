package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"time"
)

type Status struct {
	Water int `json:"water"`
	Wind  int `json:"wind"`
}

type Weather struct {
	Status Status `json:"status"`
}

func main() {
	rand.Seed(time.Now().Unix())
	go generateWeather()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		data, err := ioutil.ReadFile("weather.json")
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		var weatherData Weather
		err = json.Unmarshal(data, &weatherData)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		var message string
		if weatherData.Status.Water < 5 {
			message = "Aman"
		} else if weatherData.Status.Water > 5 && weatherData.Status.Water < 9 {
			message = "Siaga" 
		} else if weatherData.Status.Water > 8 {
			message = "Bahaya"
		} else if weatherData.Status.Wind < 6 {
			message = "Aman"
		} else if weatherData.Status.Wind > 6 && weatherData.Status.Wind < 16 {
			message = "Siaga"
		} else if weatherData.Status.Wind > 15 {
			message = "Bahaya"
		}

		htmlResponse := fmt.Sprintf("<html><body><h1>Weather Report</h1><p>Temperature: %d meter</p><p>Condition: %d meter/detik</p><p>%s</p></body></html>", weatherData.Status.Water, weatherData.Status.Wind, message)

		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		fmt.Fprint(w, htmlResponse)
	})
	http.ListenAndServe(":8080", nil)
}

func generateWeather(){
	for {
		weather := Weather{
			Status: Status{
				Water: rand.Intn(100)+1,
				Wind: rand.Intn(100)+1,
			},
		}
		data, err := json.Marshal(weather)

		if err != nil {
			panic(err)
		}

		_ = os.WriteFile("weather.json", data, 0644)

		time.Sleep(15*time.Second)
	}
}

