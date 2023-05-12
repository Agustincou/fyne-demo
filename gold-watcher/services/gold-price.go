package services

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

const (
	Currency            = "USD"
	GoldPriceOrgBaseURL = "https://data-asg.goldprice.org/dbXRates"
)

type GoldPriceClient interface {
	Get() (*GoldPrice, error)
}

type HTTPGoldPriceClient struct {
	GoldPriceClient
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

func NewHTTPGoldPriceClient(baseURL string, client *http.Client, currency string) GoldPriceClient {
	return &HTTPGoldPriceClient{
		baseURL:  baseURL,
		client:   client,
		currency: currency,
	}
}

func (g *HTTPGoldPriceClient) Get() (*GoldPrice, error) {
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
