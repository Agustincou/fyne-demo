package principal

import (
	"log"
	"os"
)

type Content struct {
	//App fyne.App -> if needed
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}

func NewContent() *Content {
	var cont Content

	cont.InfoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	cont.ErrorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	return &cont
}
