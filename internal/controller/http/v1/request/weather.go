package request

type Weather struct {
	Lat string `json:"lat"       validate:"required"  example:"69.240562"`
	Lon string `json:"lon"  validate:"required"  example:"41.311081"`
}
