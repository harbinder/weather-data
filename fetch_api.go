package main

import (
    "encoding/csv"
    "encoding/json"
    "fmt"
    "net/http"
    "os"
    "time"
)

type WeatherResponse struct {
    Main struct {
        Temp float64 `json:"temp"`
    } `json:"main"`
    Weather []struct {
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
    writer.Write([]string{ts, city, fmt.Sprintf("%.2f", data.Main.Temp), data.Weather[0].Description})
}

func main() {
    apiKey := os.Getenv("API_KEY")
    if apiKey == "" {
        fmt.Println("API_KEY is not set")
        return
    }

    cities := []string{"Delhi","London","New York","Tokyo","Sydney","Paris","Mumbai","Dubai","Beijing","Berlin"}

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
        writer.Write([]string{"Timestamp","City","Temperature (°C)","Weather"})
    }

    for _, city := range cities {
        fetchWeather(city, apiKey, writer)
    }
}
