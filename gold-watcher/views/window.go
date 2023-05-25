package views

import (
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
)

type LocalWindow struct {
	fyneWin fyne.Window
	content *Content
	tabs    *container.AppTabs
}

func NewWindow(title string) *LocalWindow {
	var localWin LocalWindow

	localWin.fyneWin = fyne.CurrentApp().NewWindow(title)
	localWin.content = NewContent()

	localWin.content.Repository.Migrate()

	localWin.tabs = container.NewAppTabs(
		container.NewTabItemWithIcon("Prices", theme.HomeIcon(), localWin.content.GetGraphTabContainer()),
		container.NewTabItemWithIcon("Holdings", theme.InfoIcon(), localWin.content.GetHoldingsTabContainer()),
	)

	return &localWin
}

func (l *LocalWindow) ShowAndRun() {
	//setup window content
	l.tabs.SetTabLocation(container.TabLocationTop) //default TabLocationTop

	l.fyneWin.SetContent(container.NewVBox(
		l.content.GetPriceContainer(),
		l.content.ToolBar,
		l.tabs,
	),
	)

	//setup window view
	l.fyneWin.Resize(fyne.Size{Width: 720, Height: 480})
	//l.fyneWin.SetFixedSize(true)
	l.fyneWin.SetMaster()
	l.fyneWin.CenterOnScreen()

	go func() {
		for range time.Tick(time.Second * 5) {
			l.content.refreshAll()
		}
	}()

	l.fyneWin.ShowAndRun()
}
