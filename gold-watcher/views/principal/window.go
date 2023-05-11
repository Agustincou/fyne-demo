package principal

import "fyne.io/fyne/v2"

type LocalWindow struct {
	FyneWin fyne.Window
	content *Content
}

func NewWindow(title string) *LocalWindow {
	var localWin LocalWindow

	localWin.FyneWin = fyne.CurrentApp().NewWindow(title)
	localWin.content = NewContent()

	//set window content
	//localWin.fyneWin.SetContent()

	//setup window view
	localWin.FyneWin.Resize(fyne.Size{Width: 300, Height: 200})
	localWin.FyneWin.SetFixedSize(true)
	localWin.FyneWin.SetMaster()
	localWin.FyneWin.CenterOnScreen()

	return &localWin
}
