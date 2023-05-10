package main

import (
	"testing"

	"fyne.io/fyne/v2/test"
	"github.com/stretchr/testify/assert"
)

func Test_makeUI(t *testing.T) {
	testContent := NewContent()

	//Simulate an action over a widget
	test.Type(testContent.EditWidget, "Hello")

	assert.Equal(t, "Hello", testContent.PreviewWidget.String(), "Failed -- did not find expected value in preview")
}

func Test_runApp(t *testing.T) {
	_ = test.NewApp()
	testWin := NewWindow("Test MarkDown")

	//fyne.CurrentApp().Run()

	test.Type(testWin.content.EditWidget, "Some Text")

	assert.Equal(t, "Some Text", testWin.content.PreviewWidget.String())
}
