package weather

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Структура для текущей погоды
type CurrentWeather struct {
	Temperature float64 `json:"temperature"`
	Windspeed   float64 `json:"windspeed"`
}

// Структура для часовых данных
type HourlyWeather struct {
	Time                []string  `json:"time"`
	ApparentTemperature []float64 `json:"apparent_temperature"`
}

// Вся структура ответа API
type Forecast struct {
	CurrentWeather CurrentWeather `json:"current_weather"`
	Hourly         HourlyWeather  `json:"hourly"`
}

// Функция для вычисления среднего значения
func average(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	var sum float64
	for _, v := range values {
		sum += v
	}
	return sum / float64(len(values))
}

// Функция для фильтрации данных за последние N часов
func filterLastHours(times []string, temps []float64, hours int) []float64 {
	now := time.Now()
	var filteredTemps []float64

	for i, t := range times {
		parsedTime, err := time.Parse("2006-01-02T15:04", t)
		if err != nil {
			fmt.Println("Ошибка при разборе времени:", err)
			continue
		}

		if now.Sub(parsedTime).Hours() <= float64(hours) {
			filteredTemps = append(filteredTemps, temps[i])
		}
	}
	return filteredTemps
}

// Функция для получения данных о погоде
func GetWeatherData(url string, hours int) (float64, CurrentWeather, error) {
	// Отправляем GET-запрос
	resp, err := http.Get(url)
	if err != nil {
		return 0, CurrentWeather{}, fmt.Errorf("ошибка запроса: %w", err)
	}
	defer resp.Body.Close()

	// Читаем тело ответа
	body, _ := io.ReadAll(resp.Body)

	// Декодируем JSON
	var forecast Forecast
	err = json.Unmarshal(body, &forecast)
	if err != nil {
		return 0, CurrentWeather{}, fmt.Errorf("ошибка при разборе JSON: %w", err)
	}

	// Фильтруем данные за последние N часов
	lastHoursTemps := filterLastHours(forecast.Hourly.Time, forecast.Hourly.ApparentTemperature, hours)

	// Вычисляем среднюю ощущаемую температуру
	avgApparentTemp := average(lastHoursTemps)

	return avgApparentTemp, forecast.CurrentWeather, nil
}
