package internal

import (
	"fmt"
	"os"
	"strconv"

	"github.com/golang/geo/s2"
	"github.com/thoeni/google-homebase/pkg/apple"
)

func IsHome(d apple.Device) bool {
	homeLat, err := strconv.ParseFloat(os.Getenv("LAT"), 64)
	if err != nil {
		fmt.Println("cannot convert to float")
	}
	homeLong, err := strconv.ParseFloat(os.Getenv("LNG"), 64)
	if err != nil {
		fmt.Println("cannot convert to float")
	}

	phoneLocation := s2.LatLngFromDegrees(d.Location.Lat, d.Location.Long)
	homeLocation := s2.LatLngFromDegrees(homeLat, homeLong)

	var R = 6371e3 // earth radius in meters
	distance := R * phoneLocation.Distance(homeLocation).Abs().Radians()

	return distance <= 30 // threshold in meters
}
