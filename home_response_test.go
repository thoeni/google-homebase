package main

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

func Test_Template(t *testing.T) {
	b, err := ioutil.ReadFile("testdata/fmip-success.json")
	if err != nil {
		t.FailNow()
	}

	r := AppleResponse{}
	err = json.Unmarshal(b, &r)
	d, _ := r.Get("iPhone X")

	expected_it := "John molto probabilmente si trova a casa.\nUltimo aggiornamento delle 11 e 49 del giorno 21 Luglio.\nStato batteria del suo cellulare: 91 percento.\n"
	assert.Equal(t, expected_it, it(r.UserInfo.FirstName, d))

	expected_en := "John is likely to be at home right now.\nLast updated at 10 49 of Saturday, 21 of July.\nHis phone battery is charged at 91 percent.\n"
	assert.Equal(t, expected_en, en(r.UserInfo.FirstName, d))
}
