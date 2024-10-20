package main

import (
	"github.com/metal3d/pifpaf/internal/ui"

	"github.com/spf13/cobra"
)

const (
	shortDescription = "multilogs is a tool to run multiple commands and display their output in a grid layout."
	longDescription  = shortDescription + "\n" +
		"Each command should be one string, that means that you certainly need to quote the command if it has spaces."
)

func main() {
	/// cmds := []string{}
	maxCols := 3

	cmd := &cobra.Command{
		Use:   "multilogs [options] command1 [command2] ...",
		Short: shortDescription,
		Long:  longDescription,
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			ui.UI(args, maxCols)
		},
	}
	cmd.Flags().IntVarP(&maxCols, "max-cols", "c", 3, "Maximum number of columns in the grid layout")
	// commands are the list of strings arguments passed to the commands
	cmd.Execute()
}
