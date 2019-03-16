package apple

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
