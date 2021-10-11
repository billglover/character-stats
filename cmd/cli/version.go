package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of the CLI",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Skritter CLI v?.?.? (dev)")
	},
}
