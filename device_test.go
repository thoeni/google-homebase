package main

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_FindByName(t *testing.T) {
	b, err := ioutil.ReadFile("testdata/fmip-success.json")
	if err != nil {
		t.FailNow()
	}

	r := AppleResponse{}
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
