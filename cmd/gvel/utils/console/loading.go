package console

import (
	"fmt"
	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	"io"
	"time"
)

var loadSpinner *spinner.Spinner
var defaultCharset = []string{
	"💞",
	"💗",
	"💖",
}
var DefaultLoadWriter io.Writer = color.Output

func StartLoading(format string, args ...interface{}) {
	if loadSpinner != nil {
		loadSpinner.Stop()
	}

	message := fmt.Sprintf(format, args...)

	loadSpinner = spinner.New(defaultCharset, 100*time.Millisecond)
	loadSpinner.FinalMSG = message + "\n"
	loadSpinner.Suffix = " " + message
	loadSpinner.HideCursor = true
	loadSpinner.Writer = DefaultLoadWriter
	loadSpinner.Start()
}

func StopLoading() {
	if loadSpinner != nil {
		loadSpinner.Stop()
	}
	loadSpinner = nil
}
