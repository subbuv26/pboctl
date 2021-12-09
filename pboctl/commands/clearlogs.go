package pboctl

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// pushCmd displays configuration
var clearLogsCmd = &cobra.Command{
	Use:   "clear-logs",
	Short: "clears logs",
	Long:  `clears logs`,
	Run:   clearLogsRun,
}

func clearLogsRun(cmd *cobra.Command, args []string) {
	fl := viper.GetString("log-file")
	if fl == "" {
		cmd.Help()
		return
	}
	if err := os.Truncate(fl, 0); err != nil {
		fmt.Println("Failed to truncate log file")
	}
}

func init() {
	pboctlCmd.AddCommand(clearLogsCmd)
}
