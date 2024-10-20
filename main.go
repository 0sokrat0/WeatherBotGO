package main

import (
	"TGBOT/weather"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

func main() {
	bot, err := tgbotapi.NewBotAPI("7893146217:AAEHa4WabeS2FFxko92TiPBBLvRCO0yqjB4")
	if err != nil {
		log.Panic(err)

	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			var msg tgbotapi.MessageConfig

			switch update.Message.Command() {
			case "start":
				avgTemp, currentWeather, err := req()
				if err != nil {
					msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Ошибка получения данных о погоде.")
				} else {
					// Форматируем сообщение о погоде
					weatherMsg := formatWeatherMessage(currentWeather.Temperature, currentWeather.Windspeed, avgTemp)

					// Создаем сообщение с поддержкой Markdown
					msg = tgbotapi.NewMessage(update.Message.Chat.ID, weatherMsg)
					msg.ParseMode = "Markdown"
				}
			case "help":
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, "I can help you with the following commands:\n/start - Start the bot\n/help - Display this help message")
			default:
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, "I don't know that command")
			}

			bot.Send(msg)
		}
	}
}

func formatWeatherMessage(temperature, windspeed, avgTemp float64) string {
	return fmt.Sprintf(
		"🌡 *Температура*: %.f°C\n🍃 *Ветер*: %.f м/с\n☁️ *Средняя ощущаемая*: %.f°C",
		temperature, windspeed, avgTemp,
	)
}

func req() (float64, weather.CurrentWeather, error) {
	// URL для запроса данных
	url := "https://api.open-meteo.com/v1/forecast?latitude=55.7522&longitude=37.6156&current_weather=true&hourly=apparent_temperature&forecast_days=1"

	// Получаем данные о погоде
	avgTemp, currentWeather, err := weather.GetWeatherData(url, 6)
	if err != nil {
		return 0, weather.CurrentWeather{}, fmt.Errorf("ошибка: %w", err)
	}

	// Возвращаем данные, если всё успешно
	return avgTemp, currentWeather, nil
}
