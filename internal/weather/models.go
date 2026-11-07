package weather

type Current struct {
	Temperature         float64 `json:"temperature"`
	RelativeHumidity    int     `json:"relativeHumidity"`
	ApparentTemperature float64 `json:"apparentTemperature"`
	IsDay               bool    `json:"isDay"`
	Precipitation       float64 `json:"precipitation"`
	WeatherCode         int     `json:"weatherCode"`
	CloudCover          int     `json:"cloudCover"`
	Pressure            float64 `json:"pressure"`
	WindSpeed           float64 `json:"windSpeed"`
	WindDirection       int     `json:"windDirection"`
	WindGusts           float64 `json:"windGusts"`
}
type Hour struct {
	Time                     string  `json:"time"`
	Temperature              float64 `json:"temperature"`
	ApparentTemperature      float64 `json:"apparentTemperature"`
	Humidity                 int     `json:"humidity"`
	PrecipitationProbability int     `json:"precipitationProbability"`
	Precipitation            float64 `json:"precipitation"`
	WeatherCode              int     `json:"weatherCode"`
	WindSpeed                float64 `json:"windSpeed"`
	Pressure                 float64 `json:"pressure"`
	UVIndex                  float64 `json:"uvIndex"`
	Visibility               float64 `json:"visibility"`
	DewPoint                 float64 `json:"dewPoint"`
}
type Day struct {
	Date                        string  `json:"date"`
	WeatherCode                 int     `json:"weatherCode"`
	TemperatureMax              float64 `json:"temperatureMax"`
	TemperatureMin              float64 `json:"temperatureMin"`
	Sunrise                     string  `json:"sunrise"`
	Sunset                      string  `json:"sunset"`
	DaylightDuration            float64 `json:"daylightDuration"`
	UVIndexMax                  float64 `json:"uvIndexMax"`
	PrecipitationSum            float64 `json:"precipitationSum"`
	PrecipitationProbabilityMax int     `json:"precipitationProbabilityMax"`
	WindSpeedMax                float64 `json:"windSpeedMax"`
}
type Overview struct {
	Current  Current `json:"current"`
	Hourly   []Hour  `json:"hourly"`
	Daily    []Day   `json:"daily"`
	Timezone string  `json:"timezone"`
	Updated  string  `json:"updated"`
}
type City struct {
	Name      string  `json:"name"`
	Country   string  `json:"country"`
	Region    string  `json:"region,omitempty"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Timezone  string  `json:"timezone"`
}
type APIError struct {
	Error string `json:"error"`
}
