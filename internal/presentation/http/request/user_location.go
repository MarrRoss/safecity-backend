package request

type AddUserLocationRequest struct {
	Latitude  float64 `json:"latitude" format:"float64"`
	Longitude float64 `json:"longitude" format:"float64"`
	Battery   int     `json:"battery" format:"int"`
} // @name AddUserLocationRequest

// Time time.Time `json:"time"`
