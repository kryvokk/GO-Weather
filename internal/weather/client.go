package weather

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type Client struct {
	http *http.Client
	base string
}

func NewClient() *Client {
	return &Client{http: &http.Client{Timeout: 12 * time.Second}, base: "https://api.open-meteo.com/v1/forecast"}
}

type forecast struct {
	Timezone string `json:"timezone"`
	Current  struct {
		Temperature   float64 `json:"temperature_2m"`
		Humidity      int     `json:"relative_humidity_2m"`
		Apparent      float64 `json:"apparent_temperature"`
		IsDay         int     `json:"is_day"`
		Precipitation float64 `json:"precipitation"`
		Code          int     `json:"weather_code"`
		Clouds        int     `json:"cloud_cover"`
		Pressure      float64 `json:"pressure_msl"`
		Wind          float64 `json:"wind_speed_10m"`
		Direction     int     `json:"wind_direction_10m"`
		Gusts         float64 `json:"wind_gusts_10m"`
	} `json:"current"`
	Hourly struct {
		Time       []string  `json:"time"`
		Temp       []float64 `json:"temperature_2m"`
		Apparent   []float64 `json:"apparent_temperature"`
		Humidity   []int     `json:"relative_humidity_2m"`
		PrecipProb []int     `json:"precipitation_probability"`
		Code       []int     `json:"weather_code"`
		Precip     []float64 `json:"precipitation"`
		Wind       []float64 `json:"wind_speed_10m"`
		Pressure   []float64 `json:"pressure_msl"`
		UV         []float64 `json:"uv_index"`
		Visibility []float64 `json:"visibility"`
		DewPoint   []float64 `json:"dew_point_2m"`
	} `json:"hourly"`
	Daily struct {
		Time     []string  `json:"time"`
		Code     []int     `json:"weather_code"`
		MaxTemp  []float64 `json:"temperature_2m_max"`
		MinTemp  []float64 `json:"temperature_2m_min"`
		Sunrise  []string  `json:"sunrise"`
		Sunset   []string  `json:"sunset"`
		Daylight []float64 `json:"daylight_duration"`
		UV       []float64 `json:"uv_index_max"`
		Precip   []float64 `json:"precipitation_sum"`
		Prob     []int     `json:"precipitation_probability_max"`
		Wind     []float64 `json:"wind_speed_10m_max"`
	} `json:"daily"`
}

func (c *Client) Fetch(ctx context.Context, lat, lon float64) (Overview, error) {
	q := url.Values{"latitude": {fmt.Sprint(lat)}, "longitude": {fmt.Sprint(lon)}, "timezone": {"auto"}, "forecast_days": {"7"}, "temperature_unit": {"celsius"}, "wind_speed_unit": {"kmh"}, "current": {"temperature_2m,relative_humidity_2m,apparent_temperature,is_day,precipitation,weather_code,cloud_cover,pressure_msl,wind_speed_10m,wind_direction_10m,wind_gusts_10m"}, "hourly": {"temperature_2m,relative_humidity_2m,apparent_temperature,precipitation_probability,precipitation,weather_code,wind_speed_10m,pressure_msl,uv_index,visibility,dew_point_2m"}, "daily": {"weather_code,temperature_2m_max,temperature_2m_min,sunrise,sunset,daylight_duration,uv_index_max,precipitation_sum,precipitation_probability_max,wind_speed_10m_max"}}
	req, e := http.NewRequestWithContext(ctx, "GET", c.base+"?"+q.Encode(), nil)
	if e != nil {
		return Overview{}, e
	}
	res, e := c.http.Do(req)
	if e != nil {
		return Overview{}, e
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return Overview{}, fmt.Errorf("weather provider returned HTTP %d", res.StatusCode)
	}
	var f forecast
	if e = json.NewDecoder(res.Body).Decode(&f); e != nil {
		return Overview{}, e
	}
	o := Overview{Timezone: f.Timezone, Updated: time.Now().UTC().Format(time.RFC3339), Current: Current{f.Current.Temperature, f.Current.Humidity, f.Current.Apparent, f.Current.IsDay == 1, f.Current.Precipitation, f.Current.Code, f.Current.Clouds, f.Current.Pressure, f.Current.Wind, f.Current.Direction, f.Current.Gusts}}
	for i, t := range f.Hourly.Time {
		o.Hourly = append(o.Hourly, Hour{t, getf(f.Hourly.Temp, i), getf(f.Hourly.Apparent, i), geti(f.Hourly.Humidity, i), geti(f.Hourly.PrecipProb, i), getf(f.Hourly.Precip, i), geti(f.Hourly.Code, i), getf(f.Hourly.Wind, i), getf(f.Hourly.Pressure, i), getf(f.Hourly.UV, i), getf(f.Hourly.Visibility, i), getf(f.Hourly.DewPoint, i)})
	}
	for i, t := range f.Daily.Time {
		o.Daily = append(o.Daily, Day{t, geti(f.Daily.Code, i), getf(f.Daily.MaxTemp, i), getf(f.Daily.MinTemp, i), gets(f.Daily.Sunrise, i), gets(f.Daily.Sunset, i), getf(f.Daily.Daylight, i), getf(f.Daily.UV, i), getf(f.Daily.Precip, i), geti(f.Daily.Prob, i), getf(f.Daily.Wind, i)})
	}
	return o, nil
}
func getf(a []float64, i int) float64 {
	if i < len(a) {
		return a[i]
	}
	return 0
}
func geti(a []int, i int) int {
	if i < len(a) {
		return a[i]
	}
	return 0
}
func gets(a []string, i int) string {
	if i < len(a) {
		return a[i]
	}
	return ""
}
