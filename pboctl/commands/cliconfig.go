package pboctl

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// showCmd displays configuration
var cfgCmd = &cobra.Command{
	Use:   "cli-config",
	Short: "Updates configuration",
	Long:  `Updates configuration`,
	Run:   cfgRun,
}

var cfgCmdArgs = map[string]bool{
	"host":      true,
	"log-type":  true,
	"log-file":  true,
	"log-level": true,
	"port":      true,
}

func cfgRun(cmd *cobra.Command, args []string) {
	for _, arg := range args {
		argItems := strings.Split(arg, "=")
		if len(argItems) != 2 {
			fmt.Println("Skipping Invalid Argument:", arg)
			continue
		}
		k, v := argItems[0], argItems[1]

		fmt.Println(k, v)
		if _, ok := cfgCmdArgs[arg]; ok {
			viper.Set(k, v)
		} else {
			fmt.Println("Skipping Unsupported Argument:", k)
		}
	}
	viper.WriteConfig()
}

var usageFunc func(*cobra.Command) error

func cfgUsage(cmd *cobra.Command) error {
	usageFunc(cmd)
	fmt.Println("Specific usage: Global Flags <flag>=<value> format to update config")
	fmt.Println("Note: except --cli-config")
	fmt.Println("Eg: pboctl cli-config --log-type=file")
	fmt.Println()
	return nil
}

func init() {
	pboctlCmd.AddCommand(cfgCmd)
	usageFunc = cfgCmd.UsageFunc()
	cfgCmd.SetUsageFunc(cfgUsage)
}
