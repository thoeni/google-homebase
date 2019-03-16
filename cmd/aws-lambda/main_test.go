package main_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/golang/geo/s2"
	"github.com/stretchr/testify/assert"
)

func Test_Distance(t *testing.T) {
	phoneLocation := s2.LatLngFromDegrees(51.529991, 0.185867)
	homeLocation := s2.LatLngFromDegrees(51.5301, -0.18556933)

	fmt.Println(phoneLocation.Distance(homeLocation).Abs())
}

func Test_Timestamp(t *testing.T) {
	ts := time.Unix(1532177632132/1000, 0)

	assert.Equal(t, 21, ts.Day())
	assert.Equal(t, time.July, ts.Month())
}
