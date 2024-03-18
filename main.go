package main

import (
	"bytes"
	"encoding/json"

	"html/template"
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

		tmpl, err := template.ParseFiles("weather.html")
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		dataMap := map[string]interface{}{
			"water": weatherData.Status.Water,
			"wind":   weatherData.Status.Wind,
			"Message":     message,
		}

		// Create a new file to save the populated template
		reportFile, err := ioutil.TempFile("", "weather_*.html")
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		defer reportFile.Close()

		var resultBuffer bytes.Buffer
		err = tmpl.Execute(&resultBuffer, dataMap)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Respond to the client with the generated HTML content
		w.Header().Set("Content-Type", "text/html")
		w.Write(resultBuffer.Bytes())
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

