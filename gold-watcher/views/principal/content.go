package principal

import (
	"fmt"
	"image/color"
	"log"
	"net/http"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/Agustincou/fyne-demo/gold-watcher/services"
	"github.com/Agustincou/fyne-demo/gold-watcher/views/images"
)

var (
	_grey  = color.NRGBA{R: 155, G: 155, B: 155, A: 255}
	_green = color.NRGBA{R: 0, G: 180, B: 0, A: 255}
	_red   = color.NRGBA{R: 180, G: 0, B: 0, A: 255}
)

type Content struct {
	//App fyne.App -> if needed
	InfoLog  *log.Logger
	ErrorLog *log.Logger

	//Services
	GoldPriceService       services.GoldPriceClient
	ImageDownloaderService services.ImageDownloader

	ToolBar *widget.Toolbar

	//Containers
	priceContainer                       *fyne.Container
	openPrice, currentPrice, changePrice *canvas.Text

	imageContainer *fyne.Container
	imageGraph     *canvas.Image
}

func NewContent() *Content {
	var c Content

	c.InfoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	c.ErrorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	c.GoldPriceService = services.NewHTTPGoldPriceClient(services.GoldPriceOrgBaseURL, http.DefaultClient, services.Currency)
	c.ImageDownloaderService = services.NewHTTPImageDownloader(http.DefaultClient, services.GoldPriceOrgChartURL, services.DownloadedFileName)

	c.openPrice = canvas.NewText("Open: Unreachable", _grey)
	c.currentPrice = canvas.NewText("Current: Unreachable", _grey)
	c.changePrice = canvas.NewText("Change: Unreachable", _grey)

	c.priceContainer = container.NewGridWithColumns(3,
		c.openPrice,
		c.currentPrice,
		c.changePrice,
	)

	c.imageGraph = canvas.NewImageFromResource(images.ResourceUnreachablePng)
	c.imageContainer = container.NewVBox(c.imageGraph)

	c.ToolBar = widget.NewToolbar(
		widget.NewToolbarSpacer(),
		widget.NewToolbarAction(theme.DocumentCreateIcon(), func() {}),
		widget.NewToolbarAction(theme.ViewRefreshIcon(), c.refreshAll),
		widget.NewToolbarAction(theme.SettingsIcon(), func() {}),
	)

	return &c
}

func (c *Content) GetPriceContainer() *fyne.Container {
	c.refreshGoldPrices()

	c.priceContainer = container.NewGridWithColumns(3,
		c.openPrice,
		c.currentPrice,
		c.changePrice,
	)

	return c.priceContainer
}

func (c *Content) GetGraphTabContainer() *fyne.Container {
	c.refreshImage()

	//Previous refresh all don't render the initial image for unknown reasons. It only works with NewVBox
	c.imageContainer = container.NewVBox(c.imageGraph)

	return c.imageContainer
}

func (c *Content) refreshGoldPrices() {
	goldPrice, err := c.GoldPriceService.Get()

	if err != nil {
		c.openPrice.Color = _grey
		c.currentPrice.Color = _grey
		c.changePrice.Color = _green
	} else {
		colorToDisplay := _green

		if goldPrice.Price < goldPrice.PreviousClose {
			colorToDisplay = _red
		}

		c.currentPrice.Color = colorToDisplay
		c.changePrice.Color = colorToDisplay

		c.openPrice.Text = fmt.Sprintf("Open: $%.4f %s", goldPrice.PreviousClose, services.Currency)
		c.currentPrice.Text = fmt.Sprintf("Current: $%.4f %s", goldPrice.Price, services.Currency)
		c.changePrice.Text = fmt.Sprintf("Change: $%.4f %s", goldPrice.Change, services.Currency)
	}

	c.openPrice.Alignment = fyne.TextAlignLeading
	c.currentPrice.Alignment = fyne.TextAlignCenter
	c.changePrice.Alignment = fyne.TextAlignTrailing
}

func (c *Content) refreshImage() {
	if err := c.ImageDownloaderService.Download(); err != nil {
		//use bundle image
		c.imageGraph = canvas.NewImageFromResource(images.ResourceUnreachablePng)
	} else {
		c.imageGraph = canvas.NewImageFromFile(services.DownloadedFileName)
	}

	c.imageGraph.SetMinSize(fyne.Size{
		Width:  720,
		Height: 410,
	})

	c.imageGraph.FillMode = canvas.ImageFillOriginal
}

func (c *Content) refreshAll() {
	c.InfoLog.Println("Refreshing data...")
	c.refreshGoldPrices()
	c.refreshImage()

	c.imageContainer.Objects = []fyne.CanvasObject{c.imageGraph}
	c.imageContainer.Refresh()

	c.priceContainer.Objects = []fyne.CanvasObject{c.openPrice, c.currentPrice, c.changePrice}
	c.priceContainer.Refresh()
}
