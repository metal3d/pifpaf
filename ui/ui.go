package ui

import (
	"fmt"
	"io"
	"log"
	"os/exec"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/rivo/tview"
)

var launchedCommands = []*exec.Cmd{}

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

// LaunchCommand starts a command and writes its stdout and stderr to the provided writer.
func LaunchCommand(w io.Writer, command string, args []string, wg *sync.WaitGroup) {
	cmd := exec.Command(command, args...)
	cmd.WaitDelay = 20 * time.Second
	launchedCommands = append(launchedCommands, cmd)
	stdout, stderr, err := getCommandPipes(cmd)
	defer func(cmd *exec.Cmd) {
		// kill the child process of cmd
		log.Println("Ending command", cmd)
		wg.Done()
		// kill child Process
	}(cmd)

	if err != nil {
		writeError(w, "Failed to get command pipes", err)
		return
	}

	if err := cmd.Start(); err != nil {
		writeError(w, "Failed to start command", err)
		return
	}

	go streamOutput(w, stdout)
	go streamOutput(w, stderr)

	if err := cmd.Wait(); err != nil {
		writeError(w, "Command execution failed", err)
	} else {
		fmt.Fprintf(w, "Command %s finished successfully\n", command)
	}
	cmd.Wait()
}

// getCommandPipes sets up stdout and stderr pipes for the command.
func getCommandPipes(cmd *exec.Cmd) (io.ReadCloser, io.ReadCloser, error) {
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, nil, fmt.Errorf("error getting stdout pipe: %w", err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, nil, fmt.Errorf("error getting stderr pipe: %w", err)
	}

	return stdout, stderr, nil
}

// streamOutput reads from the provided pipe and writes to the writer.
func streamOutput(w io.Writer, pipe io.ReadCloser) {
	if _, err := io.Copy(w, pipe); err != nil {
		writeError(w, "Error streaming output", err)
	}
}

// writeError writes an error message to the writer.
func writeError(w io.Writer, message string, err error) {
	fmt.Fprintf(w, "%s: %v\n", message, err)
}

// UI creates a new tview application and runs the provided commands in a grid layout.
// The maxCols parameter specifies the maximum number of columns in the grid.
func UI(cmds []string, maxCols int) {
	wg := sync.WaitGroup{}

	app := tview.NewApplication()
	app.EnableMouse(true)

	currentCol := 0
	mainbox := tview.NewFlex().SetDirection(tview.FlexRow)
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
		log.Println("Command finished", cmd)
	}
	wg.Wait()
}
