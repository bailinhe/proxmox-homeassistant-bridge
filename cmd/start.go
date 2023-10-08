package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gitlab.com/bighotel/proxmox-ha-bridge/pkg/proxmox"
)

// startCmd starts this app
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "starts the governor api server",
	RunE: func(cmd *cobra.Command, args []string) error {
		return startServer()
	},
}

func init() {
	rootCmd.AddCommand(startCmd)

	startCmd.Flags().String("token-id", "", "Proxmox API token ID")
	viperBindFlag("token-id", startCmd.Flags().Lookup("token-id"))

	startCmd.Flags().String("token-secret", "", "Proxmox API token secret")
	viperBindFlag("token-secret", startCmd.Flags().Lookup("token-secret"))

	startCmd.Flags().String("server-url", "https://proxmox.int.bailinhe.com", "Proxmox API server URL")
	viperBindFlag("server-url", startCmd.Flags().Lookup("server-url"))

	startCmd.Flags().Bool("insecure", false, "Skip TLS verify for proxmox server")
	viperBindFlag("insecure", startCmd.Flags().Lookup("insecure"))
}

func startServer() error {
	tokenID := viper.GetString("token-id")
	if tokenID == "" {
		logger.Panic("proxmox api token not provided")
	}

	apisecret := viper.GetString("token-secret")
	if apisecret == "" {
		logger.Panic("proxmox api secret not provided")
	}

	serverURL := viper.GetString("server-url")
	if serverURL == "" {
		logger.Panic("proxmox server URL not provided")
	}

	opts := []proxmox.Opt{
		proxmox.WithAPIToken(tokenID, apisecret),
		proxmox.WithLogger(logger.Desugar()),
	}

	if viper.GetBool("insecure") {
		logger.Debugf("insecure enabled")
		opts = append(opts, proxmox.WithInsecure())
	}

	proxmox.NewServer(serverURL, opts...).Start()

	return nil
}
