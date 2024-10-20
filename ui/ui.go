package ui

import (
	"log"
	"strings"
	"sync"
	"syscall"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// UI creates a new tview application and runs the provided commands in a grid layout.
// The maxCols parameter specifies the maximum number of columns in the grid.
func UI(cmds []string, maxCols int) {
	wg := sync.WaitGroup{}

	app := tview.NewApplication()
	app.EnableMouse(true)

	currentCol := 0
	mainbox := tview.NewFlex().SetDirection(tview.FlexRow)

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEsc {
			// ask the user if they want to quit
			modal := tview.NewModal()
			modal.SetText("Do you really want to quit?\nThis will stop all running commands.")
			modal.AddButtons([]string{"Yes", "No"})
			modal.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
				if buttonIndex == 0 {
					app.Stop()
				} else {
					app.SetRoot(mainbox, true)
				}
			})
			app.SetRoot(modal, true)
		}
		return event
	})

	currentRow := tview.NewFlex().SetDirection(tview.FlexColumn)

	for _, cmd := range cmds {
		wg.Add(1)
		parts := strings.Fields(cmd)
		command := parts[0]
		args := parts[1:]
		block := NewLogBlock(&wg, command, args, app)
		currentRow.AddItem(block, 0, 1, false)
		currentCol++
		if currentCol == maxCols {
			currentCol = 0
			mainbox.AddItem(currentRow, 0, 1, false)
			currentRow = tview.NewFlex().SetDirection(tview.FlexColumn)
		}
	}

	if currentCol != 0 {
		mainbox.AddItem(currentRow, 0, 1, false)
	}

	if err := app.SetRoot(mainbox, true).Run(); err != nil {
		panic(err)
	}
	log.Println("Waiting for commands to finish")
	for _, cmd := range launchedCommands {
		log.Println("Stopping command", cmd)
		cmd.Process.Signal(syscall.SIGTERM)
	}
	wg.Wait()
}
