package internal

import (
	"encoding/json"
	"github.com/thoeni/google-homebase/pkg/apple"
	"io/ioutil"
	"testing"

	"os"

	"github.com/stretchr/testify/assert"
)

func Test_Template(t *testing.T) {
	b, err := ioutil.ReadFile("testdata/fmip-success.json")
	if err != nil {
		t.FailNow()
	}

	os.Setenv("LAT", "51.5301")
	os.Setenv("LNG", "-0.18556933")

	r := apple.Response{}
	err = json.Unmarshal(b, &r)
	assert.NoError(t, err)

	d, _ := r.Get("iPhone X")

	expectedIt := "John molto probabilmente si trova a casa.\nUltimo aggiornamento delle 11 e 49 del giorno 21 Luglio.\nStato batteria del suo cellulare: 91 percento.\n"
	assert.Equal(t, expectedIt, it(r.UserInfo.FirstName, d))

	expectedEn := "John is likely to be at home right now.\nLast updated at 10 49 of Saturday, 21 of July.\nHis phone battery is charged at 91 percent.\n"
	assert.Equal(t, expectedEn, en(r.UserInfo.FirstName, d))
}
