package ui

import (
	"strings"
	"sync"

	"github.com/rivo/tview"
)

// NewLogBlock creates a new text view that runs the provided command and writes its output to the view.
func NewLogBlock(wg *sync.WaitGroup, command string, args []string, app *tview.Application) *tview.TextView {
	textView := tview.NewTextView()
	textView.SetDynamicColors(true)
	textView.SetBorder(true)
	textView.SetTitle(command + " " + strings.Join(args, " "))
	textView.SetTitleAlign(tview.AlignLeft)
	textView.SetBorderPadding(1, 1, 2, 2)
	textView.SetScrollable(true)
	textView.SetChangedFunc(func() {
		textView.ScrollToEnd()
		app.Draw()
	})
	w := tview.ANSIWriter(textView)
	go LaunchCommand(w, command, args, wg)
	// go textViewWrite(w, command+" "+strings.Join(args, " ")+"\n")
	return textView
}
