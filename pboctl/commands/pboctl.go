package pboctl

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

var (
	cfgFile string
	// pboctlCmd is the root command to configure Programmable BIG-IP
	pboctlCmd = NewPBOctl()
)

func Execute() {
	if err := pboctlCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run(ccmd *cobra.Command, args []string) {
	ccmd.HelpFunc()(ccmd, args)
}

func persistentPreRun(ccmd *cobra.Command, args []string) {

	// if --config is passed, attempt to parse the config file
	if cfgFile != "" {

		// get the filepath
		abs, err := filepath.Abs(cfgFile)
		if err != nil {
			fmt.Println("Error reading filepath: ", err.Error())
		}

		base := filepath.Base(abs)
		path := filepath.Dir(abs)

		// Need to give file name with out extension
		viper.SetConfigName(strings.Split(base, ".")[0])
		viper.AddConfigPath(path)

		// Find and read the config file; Handle errors reading the config file
		if err := viper.ReadInConfig(); err != nil {
			fmt.Println("Failed to read config file: ", err.Error())
			os.Exit(1)
		}
	}
}

func NewPBOctl() *cobra.Command {
	pboctlCmd := &cobra.Command{
		Use:   "pboctl <subcommand> <args>",
		Short: "pboctl is the command line interface to connect to Orchestrator that configures Programmable BIG-IP",
		Long:  `pboctl is the command line interface to connect to Orchestrator that configures Programmable BIG-IP`,
	}

	pboctlCmd.Run = run

	pboctlCmd.PersistentPreRun = persistentPreRun

	return pboctlCmd
}

//setUpLogs set the log output ans the log level
func setUpLogs(logType, logFile, level string) error {
	if logType == "stdout" {
		logrus.SetOutput(os.Stdout)
	} else if logFile != "" {

		f, err := os.OpenFile(logFile, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
		if err != nil {
			fmt.Println(err)
			logrus.SetOutput(os.Stdout)
		} else {
			logrus.SetOutput(f)
		}
	}

	lvl, err := logrus.ParseLevel(level)
	if err != nil {
		return err
	}
	logrus.SetFormatter(&prefixed.TextFormatter{TimestampFormat: "2006-01-02 15:04:05", FullTimestamp: true, ForceFormatting: true})
	logrus.SetLevel(lvl)
	return nil
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".pbo.yaml"
		viper.AddConfigPath(home)
		viper.SetConfigName(".pbo")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println(err)
	}

	setUpLogs(
		viper.GetString("log-type"),
		viper.GetString("log-file"),
		viper.GetString("log-level"),
	)
	logrus.Info("pboctl client initiated")
}

func init() {

	cobra.OnInitialize(initConfig)

	// the argument value gets saved in the cfgFile Variable
	pboctlCmd.PersistentFlags().StringVar(&cfgFile, "cli-config", "", "config file (default is $HOME/.pbo.yaml)")

	// These String/StringP are to be bind with viper
	pboctlCmd.PersistentFlags().StringP("host", "H", "127.0.0.1", "BIGIP Orchestrator hostname/IP")
	pboctlCmd.PersistentFlags().String("log-type", "stdout", "The type of logging (stdout, file)")
	pboctlCmd.PersistentFlags().String("log-file", "/tmp/pbo.log", "If log-type=file, the /path/to/logfile; ignored otherwise")
	pboctlCmd.PersistentFlags().String("log-level", logrus.InfoLevel.String(), "Output level of logs (trace, debug, info, warn, error, debug)")
	pboctlCmd.PersistentFlags().StringP("port", "p", "8000", "BIGIP Orchestrator port")

	// viper saves all these in viper context, by using viper functions these values can be extracted
	viper.BindPFlag("host", pboctlCmd.PersistentFlags().Lookup("host"))
	viper.BindPFlag("log-type", pboctlCmd.PersistentFlags().Lookup("log-type"))
	viper.BindPFlag("log-file", pboctlCmd.PersistentFlags().Lookup("log-file"))
	viper.BindPFlag("log-level", pboctlCmd.PersistentFlags().Lookup("log-level"))
	viper.BindPFlag("port", pboctlCmd.PersistentFlags().Lookup("port"))
}
