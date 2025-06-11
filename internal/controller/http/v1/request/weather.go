package request

type Weather struct {
	Lat float64 `json:"lat"       validate:"required"  example:"69.240562"`
	Lon float64 `json:"lon"  validate:"required"  example:"41.311081"`
}
