package internal

type WeatherInfo struct {
	Current CurrentWeather `json:"current"`
}

type CurrentWeather struct {
	Temperature float64   `json:"temp_c"`
	Condition   Condition `json:"condition"`
}

type Condition struct {
	Text string `json:"text"`
}
