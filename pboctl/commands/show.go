package pboctl

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// showCmd displays configuration
var showCmd = &cobra.Command{
	Use:   "show-cliconfig",
	Short: "A brief description of pboctl configuration",
	Long:  `A brief description of pboctl configuration`,
	Run:   showRun,
}

func showRun(cmd *cobra.Command, args []string) {
	for k := range cfgCmdArgs {
		fmt.Println(k, ":\t", viper.Get(k))
	}
}

func init() {
	pboctlCmd.AddCommand(showCmd)
}
