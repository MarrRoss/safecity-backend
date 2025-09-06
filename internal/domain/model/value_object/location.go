package value_object

type Location struct {
	Latitude  Latitude  // Y
	Longitude Longitude // X
}

func NewLocation(lat Latitude, long Longitude) Location {
	return Location{Latitude: lat, Longitude: long}
}
