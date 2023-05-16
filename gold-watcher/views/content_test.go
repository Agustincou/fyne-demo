package views

import (
	"fmt"
	"testing"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/test"
	"github.com/Agustincou/fyne-demo/gold-watcher/repository"
	"github.com/Agustincou/fyne-demo/gold-watcher/services"
	"github.com/Agustincou/fyne-demo/gold-watcher/tests"
	"github.com/stretchr/testify/assert"
)

func Test_refreshGoldPrices(t *testing.T) {
	//setup
	app.New()

	//given
	goldPriceClientMock := tests.GoldPriceClientMock{}
	contentToTest := NewContent()

	goldPriceClientMock.On("Get").Return(&services.GoldPrice{
		Currency:      services.PreferredCurrency,
		Price:         123,
		Change:        -0.5,
		PreviousClose: 120,
		Time:          time.Now(),
	}, nil)
	contentToTest.GoldPriceService = &goldPriceClientMock

	//when
	contentToTest.refreshGoldPrices()

	//then
	assert.Equal(t, "Open: $120.0000 USD", contentToTest.openPrice.Text)
	assert.Equal(t, "Current: $123.0000 USD", contentToTest.currentPrice.Text)
	assert.Equal(t, "Change: $-0.5000 USD", contentToTest.changePrice.Text)

	assert.Equal(t, 4, len(contentToTest.ToolBar.Items))

	//tear down
	fyne.CurrentApp().Quit()
}

func Test_getAllHoldings(t *testing.T) {
	//setup
	app.New()

	//given
	repoMock := tests.RepositoryMock{}
	contentToTest := NewContent()

	contentToTest.Repository = &repoMock

	repoMock.On("AllHoldings").Return([]repository.Holding{{
		ID:            1,
		Amount:        2,
		PurchaseDate:  time.Time{},
		PurchasePrice: 123,
	}}, nil)

	//when
	holdingTableItems := contentToTest.getHoldingsTableItems()

	//then
	assert.Equal(t, 2, len(holdingTableItems))
	assert.Equal(t, "[ID Amount Price Date Delete?]", fmt.Sprintf("%v", holdingTableItems[0]))
	assert.Contains(t, fmt.Sprintf("%v", holdingTableItems[1]), "[1 2 toz $1.23 0001-01-01")

	//tear down
	fyne.CurrentApp().Quit()
}

func Test_addHoldingsDialog(t *testing.T) {
	//setup
	app.NewWithID("asdgfsdg")
	fyne.CurrentApp().NewWindow("test")

	contentToTest := NewContent()

	test.Type(contentToTest.AddHoldingsPurchaseAmountEntry, "1")
	test.Type(contentToTest.AddHoldingsPurchasePriceEntry, "1000")
	test.Type(contentToTest.AddHoldingsPurchaseDateEntry, "2022-12-12")

	assert.Equal(t, "2022-12-12", contentToTest.AddHoldingsPurchaseDateEntry.Text)
	assert.Equal(t, "1000", contentToTest.AddHoldingsPurchasePriceEntry.Text)
	assert.Equal(t, "1", contentToTest.AddHoldingsPurchaseAmountEntry.Text)

	//tear down
	fyne.CurrentApp().Quit()
}
