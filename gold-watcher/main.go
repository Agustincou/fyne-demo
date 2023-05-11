package main

import (
	"fyne.io/fyne/v2/app"
	"github.com/Agustincou/fyne-demo/gold-watcher/views/principal"
)

const (
	appTitle = "GoldWatcher"
)

func main() {
	//create fyne app
	_ = app.NewWithID("com.github.agustincou.fyne-demo.gold-watcher")

	localWindow := principal.NewWindow(appTitle)

	localWindow.FyneWin.ShowAndRun()
}
