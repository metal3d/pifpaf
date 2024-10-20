package main

import (
	"fmt"
	"runtime/debug"

	"github.com/metal3d/pifpaf/internal/ui"

	"github.com/spf13/cobra"
)

const (
	shortDescription = "multilogs is a tool to run multiple commands and display their output in a grid layout."
	longDescription  = shortDescription + "\n" +
		"Each command should be one string, that means that you certainly need to quote the command if it has spaces."
)

// Version is the version of the application, can set at build time
var Version = "dev"

func main() {
	/// cmds := []string{}
	maxCols := 3

	// root command
	cmd := &cobra.Command{
		Use:     "multilogs [options] command1 [command2] ...",
		Short:   shortDescription,
		Long:    longDescription,
		Version: buildVersion(),
		Args:    cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if maxCols < 1 {
				return cmd.Help()
			}
			ui.UI(args, maxCols)
			return nil
		},
	}
	cmd.Flags().IntVarP(
		&maxCols, "max-cols", "c", 3,
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
