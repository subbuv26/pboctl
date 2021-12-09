package pboctl

import (
	"fmt"
	"io/ioutil"

	resty "github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	RESTAPIPATH = ""
)

// pushCmd displays configuration
var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "Pushes configuration to BIGIP Orchestrator",
	Long:  `Pushes configuration to BIGIP Orchestrator`,
	Run:   pushRun,
}

func pushRun(cmd *cobra.Command, args []string) {
	fl := viper.GetString("file")
	if fl == "" {
		cmd.Help()
		return
	}

	host := viper.GetString("host")
	port := viper.GetString("port")

	fileBytes, err := ioutil.ReadFile(fl)
	if err != nil {
		fmt.Println("Failed to read file:", fl)
		return
	}

	logrus.Info("Extracted Configuration from:", fl)
	logrus.Info("Pushing Configuration to BIGIP:", host)

	client := resty.New()

	api := "http://" + host + RESTAPIPATH + ":" + port
	resp, err := client.R().
		SetBody(fileBytes).
		SetHeader("Content-Type", "application/json").
		Post(api)
	if err != nil {
		fmt.Println("Failed to push configuration. Error:", err)
		return
	}

	if resp.IsError() {
		fmt.Println("Error ocurred")
		return
	}

	fmt.Println("Succesfully Pushed Configuration")
}

func init() {
	pboctlCmd.AddCommand(pushCmd)

	pushCmd.PersistentFlags().StringP("file", "f", "", "BIGIP Resource configuration file")
	viper.BindPFlag("file", pushCmd.PersistentFlags().Lookup("file"))
}
