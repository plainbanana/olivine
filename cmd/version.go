package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:              "version",
	Short:            "return versions.",
	Long:             `All software has versions`,
	PersistentPreRun: setLogMinLevel,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("olivin version:", version)
	},
}
