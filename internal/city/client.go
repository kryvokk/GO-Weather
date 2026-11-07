package city

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/kryvna/go-weather/internal/weather"
	"net/http"
	"net/url"
	"time"
)

type Client struct{ http *http.Client }

func NewClient() *Client { return &Client{&http.Client{Timeout: 8 * time.Second}} }
func (c *Client) Search(ctx context.Context, name string) ([]weather.City, error) {
	q := url.Values{"name": {name}, "count": {"10"}, "language": {"en"}, "format": {"json"}}
	req, e := http.NewRequestWithContext(ctx, "GET", "https://geocoding-api.open-meteo.com/v1/search?"+q.Encode(), nil)
	if e != nil {
		return nil, e
	}
	res, e := c.http.Do(req)
	if e != nil {
		return nil, e
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("geocoding provider returned HTTP %d", res.StatusCode)
	}
	var data struct {
		Results []struct {
			Name, Country, Admin1, Timezone string
			Latitude, Longitude             float64
		} `json:"results"`
	}
	if e = json.NewDecoder(res.Body).Decode(&data); e != nil {
		return nil, e
	}
	out := make([]weather.City, 0, len(data.Results))
	for _, r := range data.Results {
		out = append(out, weather.City{Name: r.Name, Country: r.Country, Region: r.Admin1, Latitude: r.Latitude, Longitude: r.Longitude, Timezone: r.Timezone})
	}
	return out, nil
}
