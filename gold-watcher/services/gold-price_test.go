package services

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Get(t *testing.T) {
	//given
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != fmt.Sprintf("/%s", Currency) {
			t.Errorf("Expected to request '/%s', got: %s", Currency, r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"ts": 1683826006122,
			"tsj": 1683825996338,
			"date": "May 11th 2023, 01:26:36 pm NY",
			"items": [
				{
					"curr": "USD",
					"xauPrice": 2015.4175,
					"xagPrice": 24.2402,
					"chgXau": -16.1675,
					"chgXag": -1.1503,
					"pcXau": -0.7958,
					"pcXag": -4.5304,
					"xauClose": 2031.585,
					"xagClose": 25.3905
				}
			]
		}`))
	}))
	defer server.Close()

	//when
	price, err := NewGoldPriceClient(server.URL, http.DefaultClient, Currency).Get()

	//then
	assert.Nil(t, err)
	assert.Equal(t, 2015.4175, price.Price)
}
