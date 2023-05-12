package tests

import (
	"github.com/Agustincou/fyne-demo/gold-watcher/services"
	"github.com/stretchr/testify/mock"
)

type GoldPriceClientMock struct {
	mock.Mock
	services.GoldPriceClient
}
func (g GoldPriceClientMock) Get() (*services.GoldPrice, error) {
	args := g.Called()

	return args.Get(0).(*services.GoldPrice), args.Error(1)
}
