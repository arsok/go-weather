package internal

import (
	json2 "encoding/json"
	"errors"
	"fmt"
	internal "go-weather/internal/model"
	"io"
	"math"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"
)

var apiKey = os.Getenv("WEATHER_API_KEY")
var placeRegex = regexp.MustCompile(`^[a-zA-Z\p{L}\s\-]{2,}$`)
var httpClient = http.Client{
	Timeout: 2 * time.Second,
}

func GetWeather(inputPlace string) (string, error) {
	place, ok := sanitizePlace(inputPlace)
	if !ok {
		return "", fmt.Errorf("incorrect place name '%s'", inputPlace)
	}

	weatherJson, rawResponse, err := getWeatherJson(place)
	if rawResponse != nil {
		defer rawResponse.Body.Close()
	}
	if err != nil {
		return "", err
	}

	var weatherInfo internal.WeatherInfo
	err = json2.Unmarshal(weatherJson, &weatherInfo)
	if err != nil {
		return "", fmt.Errorf("error decoding weather data for '%s': %v", inputPlace, err)
	}

	if !hasWeatherCondition(weatherInfo) {
		return "", errors.New(fmt.Sprintf("no weather information available for '%s'", inputPlace))
	}

	temperature := formatTemperature(weatherInfo.Current.Temperature)
	weatherCondition := strings.ToLower(weatherInfo.Current.Condition.Text)
	return fmt.Sprintf("It's %s°C, %s in %s", temperature, weatherCondition, inputPlace), nil
}

func sanitizePlace(input string) (string, bool) {
	if !placeRegex.MatchString(input) {
		return "", false
	}

	return url.QueryEscape(input), true
}

func getWeatherJson(place string) ([]byte, *http.Response, error) {
	requestUrl := fmt.Sprintf("http://api.weatherapi.com/v1/current.json?q=%s&key=%s", place, apiKey)

	response, err := httpClient.Get(requestUrl)
	if err != nil {
		return nil, nil, fmt.Errorf("error retrieving weather for '%s': %v", place, err)
	}

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, nil, fmt.Errorf("error reading response body: %v", err)
	}
	return responseBody, response, nil
}

func hasWeatherCondition(weatherInfo internal.WeatherInfo) bool {
	return len(weatherInfo.Current.Condition.Text) > 0
}

func formatTemperature(temperature float64) string {
	if temperature == math.Trunc(temperature) {
		return fmt.Sprintf("%.0f", temperature)
	}
	return fmt.Sprintf("%.1f", temperature)
}
