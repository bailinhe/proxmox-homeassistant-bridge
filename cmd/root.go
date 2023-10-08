// Package cmd is our cobra/viper cli implementation
package cmd

import (
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

const appName = "proxmox-ha-bridge"

var (
	cfgFile string
	logger  *zap.SugaredLogger
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "proxmox-ha-bridge",
	Short: "Start proxmox vms with home assistant",
	Long:  "Start proxmox vms with home assistant",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file")

	rootCmd.PersistentFlags().Bool("debug", false, "enable debug logging")
	viperBindFlag("logging.debug", rootCmd.PersistentFlags().Lookup("debug"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	}

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	viper.SetEnvPrefix("proxmox_ha_bridge")
	viper.AutomaticEnv() // read in environment variables that match

	err := viper.ReadInConfig()

	setupLogging()

	if err == nil {
		logger.Infow("using config file", "file", viper.ConfigFileUsed())
	}
}

func setupLogging() {
	cfg := zap.NewProductionConfig()

	if viper.GetBool("logging.debug") {
		cfg = zap.NewDevelopmentConfig()
		cfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	} else {
		cfg.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	}

	l, err := cfg.Build()
	if err != nil {
		panic(err)
	}

	logger = l.Sugar().With("app", appName)
	defer logger.Sync() //nolint:errcheck
}

// viperBindFlag provides a wrapper around the viper bindings that handles error checks
func viperBindFlag(name string, flag *pflag.Flag) {
	err := viper.BindPFlag(name, flag)
	if err != nil {
		panic(err)
	}
}
