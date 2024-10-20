package main

import (
	"bufio"
	"fmt"
	"os"
	"runtime/debug"

	"github.com/metal3d/pifpaf/internal/ui"
	"github.com/spf13/cobra"
)

const examples = `
Examples:

	# commands in arguments
	pifpaf "ping google.com" "podman run --rm -it metal3d/xmrig"

	# commands from stdin
	echo -e "ping google.com\npodman run --rm -it metal3d/xmrig" | pifpaf

	pifpaf <<EOF
	ping google.comm
	podman run --rm -it metal3d/xmrig
	EOF

	pifpaf < file_with_commands.txt
`

const (
	shortDescription = "pifpaf is a tool to run multiple commands and display their output in a grid layout."
	longDescription  = shortDescription + "\n\n" +
		"Each command should be one string, that means that you certainly need to quote the command if it has spaces.\n" +
		"You can also pass the commands separated by newlines to stdin." +
		examples
)

// Version is the version of the application, can set at build time
var Version = "dev"

func main() {
	// maximum number of columns in the grid layout
	maxCols := 3

	// root command
	cmd := &cobra.Command{
		Use:     "pifpaf [options] command1 [command2] ...",
		Short:   shortDescription,
		Long:    longDescription,
		Version: buildVersion(),
		Args:    cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if maxCols < 1 {
				return cmd.Help()
			}

			// accept stdin command list separated by newlines
			if stdinCommands := getCommandsFromStdin(); len(stdinCommands) > 0 {
				args = append(args, stdinCommands...)
			}

			if len(args) == 0 {
				return fmt.Errorf("no command to run")
			}

			ui.UI(args, maxCols)
			return nil
		},
	}
	cmd.Flags().IntVarP(
		&maxCols, "max-cols", "c", maxCols,
		fmt.Sprintf("Maximum number of columns in the grid layout, default is %d, must be greater than 0", maxCols),
	)

	// add version subcommand
	versionSubCmd := &cobra.Command{
		Use:   "version",
		Short: "Print the version number of multilogs",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Printf("multilogs version %s\n", buildVersion())
		},
	}
	cmd.AddCommand(versionSubCmd)

	// let's go
	cmd.Execute()
}

// buildVersion returns the version of the application. If the version is "dev", it tries to get the version from the go build info.
// If it fails, it returns the default version.
// This helps to get the version of the application when it is built with go build, go install, or from the release page.
func buildVersion() string {
	if Version == "dev" {
		// detetct the version from go build release tag
		if build, ok := debug.ReadBuildInfo(); !ok {
			return Version
		} else if build.Main.Version != "" {
			return build.Main.Version
		}
	}
	return Version
}

// getCommandsFromStdin reads the commands from stdin and returns them as a slice of strings.
func getCommandsFromStdin() []string {
	stdinStat, err := os.Stdin.Stat()
	if err != nil {
		return nil
	}

	if (stdinStat.Mode() & os.ModeCharDevice) == 0 {
		var commands []string
		s := bufio.NewScanner(os.Stdin)
		for s.Scan() {
			commands = append(commands, s.Text())
		}
		return commands
	}
	return nil
}
