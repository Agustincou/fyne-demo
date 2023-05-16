package tests

import (
	"github.com/Agustincou/fyne-demo/gold-watcher/repository"
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

type RepositoryMock struct {
	mock.Mock
	repository.Repository
}

func (r *RepositoryMock) Migrate() error {
	args := r.Called()

	return args.Error(0)
}

func (r *RepositoryMock) InsertHolding(h repository.Holding) (*repository.Holding, error) {
	args := r.Called(h)

	return args.Get(0).(*repository.Holding), args.Error(1)
}

func (r *RepositoryMock) AllHoldings() ([]repository.Holding, error) {
	args := r.Called()

	return args.Get(0).([]repository.Holding), args.Error(1)
}

func (r *RepositoryMock) GetHoldingByID(id int) (*repository.Holding, error) {
	args := r.Called(id)

	return args.Get(0).(*repository.Holding), args.Error(1)
}

func (r *RepositoryMock) UpdateHolding(id int64, updated repository.Holding) error {
	args := r.Called(id, updated)

	return args.Error(0)
}

func (r *RepositoryMock) DeleteHolding(id int64) error {
	args := r.Called(id)

	return args.Error(0)
}
