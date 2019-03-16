package internal

import (
	"encoding/json"
	"github.com/thoeni/google-homebase/pkg/apple"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_FindByName(t *testing.T) {
	b, err := ioutil.ReadFile("testdata/fmip-success.json")
	if err != nil {
		t.FailNow()
	}

	r := apple.Response{}
	err = json.Unmarshal(b, &r)
	assert.NoError(t, err)

	t.Run("can unmarshal", func(t *testing.T) {
		assert.Equal(t, 2, len(r.Devices))
		assert.Equal(t, "iPhone X", r.Devices[0].Name)
	})

	t.Run("can find device", func(t *testing.T) {
		name := "iPhone X"
		d, err := r.Get(name)

		assert.NoError(t, err)
		assert.Equal(t, name, d.Name)
	})
}

func TestDevice_isHome(t *testing.T) {

	os.Setenv("LAT", "51.529977")
	os.Setenv("LNG", "-0.185478")

	type fields struct {
		Lat float64
		Lng float64
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{"Exactly same coordinates", fields{51.529977, -0.185478}, true},
		{"Coordinates near home (about 25 meters away west)", fields{51.529884, -0.185780}, true},
		{"Coordinates near home (about 25 meters away north)", fields{51.530216, -0.185648}, true},
		{"Coordinates 60 meters away", fields{51.530556, -0.185467}, false},
		{"Coordinates far away", fields{51.513137, -0.158908}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := apple.Device{
				Location: apple.Location{Lat: tt.fields.Lat, Long: tt.fields.Lng},
			}
			if got := IsHome(d); got != tt.want {
				t.Errorf("Device.isHome() = %v, want %v", got, tt.want)
			}
		})
	}
}
