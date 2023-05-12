package principal

import (
	"testing"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"github.com/Agustincou/fyne-demo/gold-watcher/services"
	"github.com/Agustincou/fyne-demo/gold-watcher/tests"
	"github.com/stretchr/testify/assert"
)

func Test_getPriceText(t *testing.T) {
	//setup
	app.New()

	//given
	goldPriceClientMock := tests.GoldPriceClientMock{}
	contentToTest := NewContent()

	goldPriceClientMock.On("Get").Return(&services.GoldPrice{
		Currency:      services.Currency,
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
