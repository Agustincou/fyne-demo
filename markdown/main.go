package main

import (
	"io/ioutil"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
)

const(
	appTitle = "Markdown"
	fileMenuName = "File"
	fileMenuOpenItemName = "Open..."
	fileMenuSaveItemName = "Save"
	fileMenuSaveAsItemName = "Save as..."
)

var filter = storage.NewExtensionFileFilter([]string{".md", ".MD"})

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
	window := NewWindow(appTitle)

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
	openMenuItem := fyne.NewMenuItem(fileMenuOpenItemName, win.openFunc())
	saveAsMenuItem := fyne.NewMenuItem(fileMenuSaveAsItemName, win.saveAsFunc())
	//SaveMenuItem saved to reference him later and alter his behavior
	win.content.SaveMenuItem = fyne.NewMenuItem(fileMenuSaveItemName, win.saveFunc())
	win.content.SaveMenuItem.Disabled = true

	fileMenuList := fyne.NewMenu(fileMenuName, openMenuItem, win.content.SaveMenuItem, saveAsMenuItem)

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
			
			if !strings.HasSuffix(strings.ToLower(write.URI().String()), ".md") {
				dialog.ShowInformation("Error", "Please name your file with a .md extension!", win.fyneWin)
			}

			write.Write([]byte(win.content.EditWidget.Text))
			win.content.CurrentFile = write.URI()

			defer write.Close()

			win.fyneWin.SetTitle(appTitle + " - " + write.URI().Name())
			win.content.SaveMenuItem.Disabled = false

		}, win.fyneWin)
		
		saveDialog.SetFileName("untitled.md")
		saveDialog.SetFilter(filter)
		saveDialog.Show()
	}
}

func (win *LocalWindow) openFunc() func() {
	return func() {
		openDialog := dialog.NewFileOpen(func(read fyne.URIReadCloser, err error) {
			if err != nil {
				dialog.ShowError(err, win.fyneWin)
			}

			if read == nil {
				return
			}

			defer read.Close()

			data, err := ioutil.ReadAll(read)
			if err != nil {
				dialog.ShowError(err, win.fyneWin)
			}

			win.content.EditWidget.SetText(string(data))

			win.content.CurrentFile = read.URI()
			win.fyneWin.SetTitle(appTitle + " - " + read.URI().Name())
			win.content.SaveMenuItem.Disabled = false
		}, win.fyneWin)

		openDialog.SetFilter(filter)
		openDialog.Show()
	}
}

func (win *LocalWindow) saveFunc() func () {
	return func() {
		if win.content.CurrentFile != nil {
			write, err := storage.Writer(win.content.CurrentFile)
			if err != nil {
				dialog.ShowError(err, win.fyneWin)
				return
			}

			write.Write([]byte(win.content.EditWidget.Text))
			defer write.Close()
		}
	}
}
