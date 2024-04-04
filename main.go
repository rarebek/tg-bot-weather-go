package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

const openWeatherMapURL = "http://api.openweathermap.org/data/2.5/forecast"
const openWeatherMapAPIKey = "bd100876d2fcf85b2ff5e4ab44875d13"

func main() {
	bot, err := tgbotapi.NewBotAPI("7081870658:AAGv1AcRYBq6dGWdlPA_15uqjn7cotD8da0")
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		if update.Message.IsCommand() {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Hello! I'm your weather bot. You can ask me about the weather by typing /weather followed by the city name.")
			bot.Send(msg)
			continue
		}

		if strings.HasPrefix(update.Message.Text, "/weather") {
			city := strings.TrimSpace(strings.TrimPrefix(update.Message.Text, "/weather"))

			weatherInfo, err := getWeather(city)
			if err != nil {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Error: Failed to retrieve weather information.")
				bot.Send(msg)
				continue
			}

			response := tgbotapi.NewMessage(update.Message.Chat.ID, "")
			response.Text = fmt.Sprintf("Weather in %s:\n%s\nTemperature: %.2f°C\nHumidity: %d%%",
				city, weatherInfo.Weather[0].Description, weatherInfo.Main.Temp-273.15, weatherInfo.Main.Humidity)

			bot.Send(response)
		}
	}
}

type weatherResponse struct {
	Weather []struct {
		Description string `json:"description"`
	} `json:"weather"`
	Main struct {
		Temp     float64 `json:"temp"`
		Humidity int     `json:"humidity"`
	} `json:"main"`
}

func getWeather(city string) (*weatherResponse, error) {
	url := fmt.Sprintf("%s?q=%s&appid=%s", openWeatherMapURL, city, openWeatherMapAPIKey)
	fmt.Println(url)
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Error fetching weather: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	var weatherInfo weatherResponse
	if err := json.NewDecoder(resp.Body).Decode(&weatherInfo); err != nil {
		return nil, err
	}

	return &weatherInfo, nil
}
