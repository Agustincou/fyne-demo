package main

import (
	"fyne.io/fyne/v2/app"
	"github.com/Agustincou/fyne-demo/gold-watcher/views"
)

const (
	_appTitle = "GoldWatcher"
	_appId = "com.github.agustincou.fyne-demo.gold-watcher"
)

func main() {
	//create fyne app
	_ = app.NewWithID(_appId)

	localWindow := views.NewWindow(_appTitle)

	localWindow.ShowAndRun()
}
