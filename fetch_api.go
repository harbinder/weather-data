package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	_ "github.com/joho/godotenv"
)

/*
{
  "coord": {
    "lon": 77.2167,
    "lat": 28.6667
  },
  "weather": [
    {
      "id": 804,
      "main": "Clouds",
      "description": "overcast clouds",
      "icon": "04d"
    }
  ],
  "base": "stations",
  "main": {
    "temp": 33.14,
    "feels_like": 40.14,
    "temp_min": 33.14,
    "temp_max": 33.14,
    "pressure": 995,
    "humidity": 62,
    "sea_level": 995,
    "grnd_level": 970
  },
  "visibility": 10000,
  "wind": {
    "speed": 1.84,
    "deg": 182,
    "gust": 1.46
  },
  "clouds": {
    "all": 100
  },
  "dt": 1750770812,
  "sys": {
    "country": "IN",
    "sunrise": 1750722866,
    "sunset": 1750773150
  },
  "timezone": 19800,
  "id": 1273294,
  "name": "Delhi",
  "cod": 200
}
*/

type WeatherResponse struct {
	Main struct {
		Temp      float64 `json:"temp"`
		FeelsLike float64 `json:"feels_like"`
		TempMin   float64 `json:"temp_min"`
		TempMax   float64 `json:"temp_max"`
		Pressure  int     `json:"pressure"`
		Humidity  int     `json:"humidity"`
		SeaLevel  int     `json:"sea_level"`
		GrndLevel int     `json:"grnd_level"`
	} `json:"main"`
	Wind struct {
		Speed float64 `json:"speed"`
		Deg   int     `json:"deg"`
		Gust  float64 `json:"gust"`
	} `json:"wind"`
	/*
	   Clouds struct {
	       All int `json:"all"`
	   } `json:"clouds"`

	   Coord struct {
	       Lon float64 `json:"lon"`
	       Lat float64 `json:"lat"`
	   } `json:"coord"`
	*/
	Sys struct {
		Country string `json:"country"`
		Sunrise int64  `json:"sunrise"`
		Sunset  int64  `json:"sunset"`
	} `json:"sys"`
	//ID int `json:"id"`
	Name     string `json:"name"`
	Timezone int    `json:"timezone"`
	DT       int64  `json:"dt"`
	//Base string `json:"base"`
	// Cod int `json:"cod"` // Not used in this example
	// Visibility is in meters
	Visibility int `json:"visibility"`
	Weather    []struct {
		Description string `json:"description"`
	} `json:"weather"`
}

func fetchWeather(city, apiKey string, writer *csv.Writer) {
	url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?q=%s&appid=%s&units=metric", city, apiKey)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error fetching", city, ":", err)
		return
	}
	defer resp.Body.Close()

	var data WeatherResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		fmt.Println("Error parsing JSON for", city, ":", err)
		return
	}

	ts := time.Now().Format("2006-01-02 15:04:05")
	timestamp := time.Unix(data.DT, 0).Format("2006-01-02 15:04:05")
	Sunrise := time.Unix(data.Sys.Sunrise, 0).Format("2006-01-02 15:04:05")
	Sunset := time.Unix(data.Sys.Sunset, 0).Format("2006-01-02 15:04:05")

	csvData := make([]string, 0, 18)
	csvData = append(csvData,
		ts,
		fmt.Sprintf("%s (%s)", data.Name, data.Sys.Country),
		data.Weather[0].Description,
		fmt.Sprintf("%.2f", data.Main.Temp),
		fmt.Sprintf("%.2f", data.Main.FeelsLike),
		fmt.Sprintf("%.2f", data.Main.TempMin),
		fmt.Sprintf("%.2f", data.Main.TempMax),
		fmt.Sprintf("%d", data.Main.SeaLevel),
		fmt.Sprintf("%d", data.Main.GrndLevel),
		fmt.Sprintf("%d%%", data.Main.Humidity),
		fmt.Sprintf("%d", data.Main.Pressure),
		timestamp, Sunrise, Sunset,
		fmt.Sprintf("%.2f", data.Wind.Speed),
		fmt.Sprintf("%d", data.Wind.Deg),
		fmt.Sprintf("%.2f", data.Wind.Gust),
		fmt.Sprintf("%d", data.Visibility),
	)
	writer.Write(csvData)

}

func main() {
	csvHeader := []string{
		"Timestamp", "City (Country)", "Weather Description",
		"Temperature (°C)", "Feels Like (°C)", "Min(°C)", "Max(°C)",
		"Sea Level", "Ground Level", "Humidity", "Pressure",
		"Timestamp (UTC)", "Sunrise (UTC)", "Sunset (UTC)",
		"Wind Speed (m/s)", "Wind Direction", "Gust", "Visibility (m)",
	}
	apiKey := os.Getenv("WEATHER_API_KEY")
	if apiKey == "" {
		fmt.Println("API_KEY is not set")
		return
	}

	cities := []string{"Delhi", "Mumbai", "Kolkata", "Bangalore"}
	//,"London","New York","Tokyo","Sydney","Paris","Mumbai","Dubai","Beijing","Berlin"}

	isNew := false
	if _, err := os.Stat("weather_data.csv"); os.IsNotExist(err) {
		isNew = true
	}

	file, err := os.OpenFile("weather_data.csv", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	if isNew {
		writer.Write(csvHeader)
	}

	for _, city := range cities {
		fetchWeather(city, apiKey, writer)
	}
}
