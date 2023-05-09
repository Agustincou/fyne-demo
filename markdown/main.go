package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

type Content struct {
	EditWidget    *widget.Entry
	PreviewWidget *widget.RichText
	CurrentFile   fyne.URI
	SaveMenuItem  *fyne.MenuItem
}

type LocalWindow struct {
	fyneWin fyne.Window
	content *Content
}

func main() {
	//create fyne app
	_ = app.New()

	//create window & content
	window := NewWindow("Markdown")

	//show windows and run app
	window.fyneWin.ShowAndRun()
}

func NewContent() *Content{
	var content Content

	content.EditWidget = widget.NewMultiLineEntry()
	content.PreviewWidget = widget.NewRichTextFromMarkdown("")

	content.EditWidget.OnChanged = func(_ string) {
		content.PreviewWidget.ParseMarkdown(content.EditWidget.Text)
	}

	//or simplified...
	//app.EditWidget.OnChanged = app.PreviewWidget.ParseMarkdown

	return &content
}

func NewWindow(title string) *LocalWindow {
	var localWin LocalWindow

	localWin.fyneWin = fyne.CurrentApp().NewWindow(title)
	localWin.content = NewContent()

	localWin.addSaveMenuList()

	//set window content
	localWin.fyneWin.SetContent(container.NewHSplit(localWin.content.EditWidget, localWin.content.PreviewWidget))

	//setup window view
	localWin.fyneWin.Resize(fyne.Size{Width: 800, Height: 500})
	localWin.fyneWin.CenterOnScreen()

	return &localWin
}

func (win *LocalWindow) addSaveMenuList() {
	openMenuItem := fyne.NewMenuItem("Open...", func() {})
	saveAsMenuItem := fyne.NewMenuItem("Save as...", win.saveAsFunc())
	//SaveMenuItem saved to reference him later and alter his behavior
	win.content.SaveMenuItem = fyne.NewMenuItem("Save", func() {})
	win.content.SaveMenuItem.Disabled = true

	fileMenuList := fyne.NewMenu("File", openMenuItem, win.content.SaveMenuItem, saveAsMenuItem)

	menu := fyne.NewMainMenu(fileMenuList)

	win.fyneWin.SetMainMenu(menu)
}

func (win *LocalWindow) saveAsFunc() func() {
	return func() {
		saveDialog := dialog.NewFileSave(func(write fyne.URIWriteCloser, err error){
			if err != nil {
				dialog.ShowError(err, win.fyneWin)

				return
			}

			if write == nil {
				//user cancelled
				return
			}

			write.Write([]byte(win.content.EditWidget.Text))
			win.content.CurrentFile = write.URI()

			defer write.Close()

			win.fyneWin.SetTitle(win.fyneWin.Title() + " - " + write.URI().Name())
			win.content.SaveMenuItem.Disabled = false

		}, win.fyneWin)
		saveDialog.Show()
	}
}
