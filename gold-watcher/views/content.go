package views

import (
	"database/sql"
	"fmt"
	"image/color"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/Agustincou/fyne-demo/gold-watcher/repository"
	"github.com/Agustincou/fyne-demo/gold-watcher/services"
	"github.com/Agustincou/fyne-demo/gold-watcher/views/images"

	_ "github.com/glebarez/go-sqlite"
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
	GoldPriceService services.GoldPriceClient

	ToolBar *widget.Toolbar

	//Containers
	priceContainer                       *fyne.Container
	openPrice, currentPrice, changePrice *canvas.Text

	imageContainer *fyne.Container
	imageGraph     *canvas.Image

	holdingsContainer *fyne.Container
	holdingsTable     *widget.Table
	holdings          [][]interface{}

	AddHoldingsPurchaseAmountEntry *widget.Entry
	AddHoldingsPurchaseDateEntry   *widget.Entry
	AddHoldingsPurchasePriceEntry  *widget.Entry

	Repository repository.Repository
}

func NewContent() *Content {
	var c Content

	c.InfoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	c.ErrorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	//Service init
	//ToDo: PreferredCurrency should be modified by specific service. To improve
	services.PreferredCurrency = fyne.CurrentApp().Preferences().StringWithFallback("currency", services.PreferredCurrency)
	c.GoldPriceService = services.NewHTTPGoldPriceClient(services.GoldPriceOrgBaseURL, http.DefaultClient, services.PreferredCurrency)

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

	c.AddHoldingsPurchaseAmountEntry = widget.NewEntry()
	c.AddHoldingsPurchaseDateEntry = widget.NewEntry()
	c.AddHoldingsPurchasePriceEntry = widget.NewEntry()

	c.ToolBar = widget.NewToolbar(
		widget.NewToolbarSpacer(),
		widget.NewToolbarAction(theme.DocumentCreateIcon(), c.openHoldingsDialog),
		widget.NewToolbarAction(theme.ViewRefreshIcon(), c.refreshAll),
		widget.NewToolbarAction(theme.SettingsIcon(), c.showPreferences),
	)

	var dbPath = os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = fyne.CurrentApp().Storage().RootURI().Path() + "/sql.db"
	}
	c.InfoLog.Println("DB in:", dbPath)
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		c.ErrorLog.Println("error initiliazing database:", err.Error())
		log.Panic()
	}

	c.Repository = repository.NewSQLiteRespository(db)

	c.holdingsTable = widget.NewTable(
		func() (int, int) {
			return len(c.holdings), len(c.holdings[0])
		},
		func() fyne.CanvasObject {
			ctr := container.NewVBox(widget.NewLabel(""))
			return ctr
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {
			// if the last cell of a row
			if i.Col == (len(c.holdings[0])-1) && i.Row != 0 {
				o.(*fyne.Container).Objects = []fyne.CanvasObject{
					c.holdings[i.Row][i.Col].(*widget.Button),
				}
			} else {
				// we're just putting in textual information
				o.(*fyne.Container).Objects = []fyne.CanvasObject{
					widget.NewLabel(c.holdings[i.Row][i.Col].(string)),
				}
			}
		})

	return &c
}

func (c *Content) GetPriceContainer() *fyne.Container {
	c.refreshGoldPrices()

	c.priceContainer.Objects = []fyne.CanvasObject{c.openPrice, c.currentPrice, c.changePrice}
	c.priceContainer.Refresh()

	return c.priceContainer
}

func (c *Content) GetHoldingsTabContainer() *fyne.Container {
	c.holdings = c.getHoldingsTableItems()

	colWidths := []float32{50, 200, 200, 200, 110}
	for i := 0; i < len(colWidths); i++ {
		c.holdingsTable.SetColumnWidth(i, colWidths[i])
	}

	c.holdingsContainer = container.NewBorder(
		nil,
		nil,
		nil,
		nil,
		container.NewAdaptiveGrid(1, c.holdingsTable),
	)

	return c.holdingsContainer
}

func (c *Content) GetGraphTabContainer() *fyne.Container {
	c.refreshImage()
	height := float32(100)

	//ToDo: To improve. Algorithm to adjust the size of the image container so that it is not larger than the MAX size of the window
	if c.imageContainer != nil && len(c.imageContainer.Objects) != 0 {
		fmt.Println("Posicion container imagen: ", c.imageContainer.Position())
		height = fyne.CurrentApp().Driver().AllWindows()[0].Canvas().Size().Height - c.imageContainer.Position().Y - 100
		fmt.Println("Alto imagen:", c.imageGraph.MinSize().Height)
		if height > c.imageGraph.MinSize().Height {
			height = c.imageGraph.MinSize().Height
		}
	}

	//Previous refresh all don't render the initial image for unknown reasons. It only works with NewVBox
	scroll := container.NewScroll(c.imageGraph)
	scroll.SetMinSize(fyne.Size{
		Height: height,
		Width: fyne.CurrentApp().Driver().AllWindows()[0].Canvas().Size().Width - 100,
	})
	grid := container.NewAdaptiveGrid(1, scroll)
	fmt.Println("Size elemento:", height)
	fmt.Println("Max Y ventana:", fyne.CurrentApp().Driver().AllWindows()[0].Canvas().Size().Height)
	c.imageContainer.Objects = []fyne.CanvasObject{grid}
	c.imageContainer.Refresh()

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

		c.openPrice.Text = fmt.Sprintf("Open: $%.4f %s", goldPrice.PreviousClose, services.PreferredCurrency)
		c.currentPrice.Text = fmt.Sprintf("Current: $%.4f %s", goldPrice.Price, services.PreferredCurrency)
		c.changePrice.Text = fmt.Sprintf("Change: $%.4f %s", goldPrice.Change, services.PreferredCurrency)
	}

	c.openPrice.Alignment = fyne.TextAlignLeading
	c.currentPrice.Alignment = fyne.TextAlignCenter
	c.changePrice.Alignment = fyne.TextAlignTrailing
}

func (c *Content) getHoldingsTableItems() [][]interface{} {
	var holdingsTableItems [][]interface{}

	holdings, err := c.Repository.AllHoldings()
	if err != nil {
		c.InfoLog.Println(err)
		return nil
	}

	//first row
	holdingsTableItems = append(holdingsTableItems, []interface{}{"ID", "Amount", "Price", "Date", "Delete?"})

	for _, holding := range holdings {
		deleteButton := widget.NewButtonWithIcon("Delete", theme.DeleteIcon(), func() {
			dialog.ShowConfirm("Delete?", "", func(deleted bool) {
				if deleted {
					if err := c.Repository.DeleteHolding(holding.ID); err != nil {
						c.ErrorLog.Println(err)
					}
				}
				// refresh the holdings table
				c.refreshHoldingsTable()
			}, fyne.CurrentApp().Driver().AllWindows()[0])
		})
		deleteButton.Importance = widget.HighImportance

		var currentRow []interface{}

		currentRow = append(currentRow, strconv.FormatInt(holding.ID, 10))
		currentRow = append(currentRow, fmt.Sprintf("%d toz", holding.Amount))
		currentRow = append(currentRow, fmt.Sprintf("$%.2f", float32(holding.PurchasePrice)/100))
		currentRow = append(currentRow, holding.PurchaseDate.Format("2006-01-02"))
		currentRow = append(currentRow, deleteButton)

		holdingsTableItems = append(holdingsTableItems, currentRow)
	}

	return holdingsTableItems
}

func (c *Content) refreshImage() {
	fyneURI, _ := storage.ParseURI(fmt.Sprintf("https://goldprice.org/charts/gold_3d_b_o_%s_x.png", services.PreferredCurrency))
	c.imageGraph = canvas.NewImageFromURI(fyneURI)

	if c.imageGraph == nil || c.imageGraph.Resource == nil {
		c.imageGraph = canvas.NewImageFromResource(images.ResourceUnreachablePng)
	}

	c.imageGraph.SetMinSize(fyne.Size{
		Width:  720,
		Height: 410,
	})

	c.imageGraph.FillMode = canvas.ImageFillContain
}

func (c *Content) refreshAll() {
	c.GetGraphTabContainer()
	c.GetPriceContainer()
}

func (c *Content) refreshHoldingsTable() {
	c.holdings = c.getHoldingsTableItems()
	c.holdingsTable.Refresh()
}

func (c *Content) openHoldingsDialog() {
	c.AddHoldingsPurchaseDateEntry.Validator = func(s string) error {
		if _, err := time.Parse("2006-01-02", s); err != nil {
			return err
		}
		return nil
	}
	c.AddHoldingsPurchaseDateEntry.PlaceHolder = "YYYY-MM-DD"

	c.AddHoldingsPurchaseAmountEntry.Validator = func(s string) error {
		if _, err := strconv.Atoi(s); err != nil {
			return err
		}
		return nil
	}

	c.AddHoldingsPurchasePriceEntry.Validator = func(s string) error {
		if _, err := strconv.ParseFloat(s, 32); err != nil {
			return err
		}
		return nil
	}

	addForm := dialog.NewForm(
		"Add Gold",
		"Add",
		"Cancel",
		[]*widget.FormItem{
			{Text: "Amount in toz", Widget: c.AddHoldingsPurchaseAmountEntry},
			{Text: "Purchase price", Widget: c.AddHoldingsPurchasePriceEntry},
			{Text: "Purchase date", Widget: c.AddHoldingsPurchaseDateEntry},
		},
		func(valid bool) {
			if valid {
				amount, _ := strconv.Atoi(c.AddHoldingsPurchaseAmountEntry.Text)
				price, _ := strconv.ParseFloat(c.AddHoldingsPurchasePriceEntry.Text, 32)
				date, _ := time.Parse("2006-01-02", c.AddHoldingsPurchaseDateEntry.Text)

				if _, err := c.Repository.InsertHolding(repository.Holding{
					Amount:        amount,
					PurchaseDate:  date,
					PurchasePrice: int(price * 100),
				}); err != nil {
					c.ErrorLog.Println(err)
				}
				c.refreshHoldingsTable()
			}
		},
		fyne.CurrentApp().Driver().AllWindows()[0],
	)

	addForm.Resize(fyne.Size{Width: 400})
	addForm.Show()
}

func (c *Content) showPreferences() {
	win := fyne.CurrentApp().NewWindow("Preferences")

	label := widget.NewLabel("Preferred currency")
	cur := widget.NewSelect([]string{"USD", "ARS", "CAD"}, func(selected string) {
		//ToDo: PreferredCurrency should be modified by specific service. To improve
		services.PreferredCurrency = selected
		fyne.CurrentApp().Preferences().SetString("currency", selected)
	})
	cur.Selected = services.PreferredCurrency

	btn := widget.NewButton("Save", func() {
		win.Close()
		c.refreshAll()
	})
	btn.Importance = widget.HighImportance

	win.SetContent(container.NewVBox(label, cur, btn))

	win.Resize(fyne.NewSize(300, 200))
	win.Show()
}
