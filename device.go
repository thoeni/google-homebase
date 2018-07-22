package main

import (
	"os"
	"strconv"

	"fmt"
	"github.com/golang/geo/s2"
)

// Device represent a subset of the Device information available on the iCloud
// API
type Device struct {
	Name         string   `json:"deviceDisplayName"`
	BatteryLevel float32  `json:"batteryLevel"`
	Location     Location `json:"location"`
}

// Location represents a subset of the Location information available on iCloud
type Location struct {
	Outdated   bool    `json:"isOld"`
	Inaccurate bool    `json:"isInaccurate"`
	Timestamp  int64   `json:"timestamp"`
	Lat        float64 `json:"latitude"`
	Long       float64 `json:"longitude"`
}

func (d Device) isHome() bool {
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
	return phoneLocation.Distance(homeLocation).Abs() <= 0.25
}
