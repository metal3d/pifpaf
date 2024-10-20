package main

import (
	"bufio"
	"fmt"
	"os"
	"runtime/debug"
	"strings"
	"text/template"

	"github.com/metal3d/pifpaf/internal/ui"
	"github.com/spf13/cobra"
)

const examples = `
  # commands in arguments
  {{.Appname}} launch "ping google.com" "podman run --rm -it metal3d/xmrig"

  # commands from stdin
  echo -e "ping google.com\npodman run --rm -it metal3d/xmrig" | %[1]s launch

  {{.Appname}} launch <<EOF
  ping google.comm
  podman run --rm -it metal3d/xmrig
  EOF

  {{.Appname}} launch < file_with_commands.txt
`

var tplConfig = map[string]any{
	"Appname": "pifpaf",
}

// Version is the version of the application, can set at build time
var Version = "dev"

func main() {
	// maximum number of columns in the grid layout
	maxCols := 3

	// launch subcommand
	launch := &cobra.Command{
		Use:     "launch [flags] command1 [command2] ...",
		Short:   "Run multiple commands and display their output in a grid layout",
		Long:    "Run multiple commands and display their output in a grid layout. The commands can be passed as arguments or from stdin.",
		Example: tplExec(examples, tplConfig),
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
	launch.Flags().IntVarP(
		&maxCols, "max-cols", "c", maxCols,
		fmt.Sprintf("Maximum number of columns in the grid layout, default is %d, must be greater than 0", maxCols),
	)

	// add version subcommand
	versionSubCmd := &cobra.Command{
		Use:   "version",
		Short: "Print the version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(buildVersion())
		},
	}

	// root command
	rootCmd := &cobra.Command{
		Use:     tplConfig["Appname"].(string),
		Version: buildVersion(),
		Example: tplExec(examples, tplConfig),
	}
	rootCmd.AddCommand(launch)
	rootCmd.AddCommand(versionSubCmd)

	// let's go
	rootCmd.Execute()
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

// convience function to execute a template string with arguments
func tplExec(s string, args any) string {
	buff := &strings.Builder{}
	t := template.Must(template.New("").Parse(s))
	t.Execute(buff, args)
	return buff.String()
}
