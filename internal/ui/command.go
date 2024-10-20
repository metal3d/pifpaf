package ui

import (
	"fmt"
	"io"
	"log"
	"os/exec"
	"sync"
	"time"
)

var launchedCommands = []*exec.Cmd{}

// LaunchCommand starts a command and writes its stdout and stderr to the provided writer.
func LaunchCommand(w io.Writer, command string, args []string, wg *sync.WaitGroup) {
	cmd := exec.Command(command, args...)
	cmd.WaitDelay = 20 * time.Second
	launchedCommands = append(launchedCommands, cmd)
	stdout, stderr, err := getCommandPipes(cmd)
	defer func(cmd *exec.Cmd) {
		log.Println("Ending command", cmd)
		wg.Done()
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
