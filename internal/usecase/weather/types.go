package weather

type Weather struct {
	Sunrise     int     `json:"sunrise"`
	Sunset      int     `json:"sunset"`
	Temperature float64 `json:"temperature"`
	FeelsLike   float64 `json:"feels_like"`
	Humidity    int     `json:"humidity"`
	WindSpeed   float64 `json:"wind_speed"`
	Description string  `json:"description"`

	County       string `json:"county"`
	Country      string `json:"country"`
	City         string `json:"city"`
	PostCode     string `json:"postcode"`
	Neighborhood string `json:"neighbourhood"`
}
