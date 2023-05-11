package services

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

const (
	Currency = "USD"
	BaseURL  = "https://data-asg.goldprice.org/dbXRates"
)

type GoldPriceClient struct {
	baseURL  string
	client   *http.Client
	currency string
}

type GoldPrices struct {
	Prices []GoldPrice `json:"items"`
}

type GoldPrice struct {
	Currency      string    `json:"curr"`
	Price         float64   `json:"xauPrice"`
	Change        float64   `json:"chgXau"`
	PreviousClose float64   `json:"xauCLose"`
	Time          time.Time `json:"-"`
}

func NewGoldPriceClient(baseURL string, client *http.Client, currency string) *GoldPriceClient {
	return &GoldPriceClient{
		baseURL:  baseURL,
		client:   client,
		currency: currency,
	}
}

func (g *GoldPriceClient) Get() (*GoldPrice, error) {
	url := fmt.Sprintf("%s/%s", g.baseURL, g.currency)

	response, err := g.client.Get(url)
	if err != nil {
		log.Printf("error requesting %s\n", g.baseURL)

		return nil, err
	}
	defer response.Body.Close()

	goldPrices := GoldPrices{}

	if err := json.NewDecoder(response.Body).Decode(&goldPrices); err != nil {
		log.Printf("error unmarsalling response %s\n", g.baseURL)

		return nil, err
	}

	return &goldPrices.Prices[0], nil
}
