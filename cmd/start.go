package cmd

import (
	"time"

	"github.com/nats-io/nats.go"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gitlab.com/bighotel/proxmox-ha-bridge/internal/api"
	natsclient "gitlab.com/bighotel/proxmox-ha-bridge/pkg/events/nats"
	"gitlab.com/bighotel/proxmox-ha-bridge/pkg/proxmox"
)

// startCmd starts this app
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "starts the proxmox ha bridge server",
	RunE: func(cmd *cobra.Command, args []string) error {
		return startServer()
	},
}

func init() {
	rootCmd.AddCommand(startCmd)

	startCmd.Flags().String("token-id", "", "Proxmox API token ID")
	viperBindFlag("proxmox.token-id", startCmd.Flags().Lookup("token-id"))

	startCmd.Flags().String("token-secret", "", "Proxmox API token secret")
	viperBindFlag("proxmox.token-secret", startCmd.Flags().Lookup("token-secret"))

	startCmd.Flags().String("server-url", "https://proxmox.int.bailinhe.com", "Proxmox API server URL")
	viperBindFlag("proxmox.server-url", startCmd.Flags().Lookup("server-url"))

	startCmd.Flags().Bool("insecure", false, "Skip TLS verify for proxmox server")
	viperBindFlag("proxmox.insecure", startCmd.Flags().Lookup("insecure"))

	startCmd.Flags().IntSlice("vm-ids", []int{}, "VM IDs to monitor")
	viperBindFlag("proxmox.vm-ids", startCmd.Flags().Lookup("vm-ids"))

	startCmd.Flags().String("probe-interval", "", "Probe interval in seconds")
	viperBindFlag("proxmox.probe-interval", startCmd.Flags().Lookup("probe-interval"))

	startCmd.Flags().String("mqtt-server", "", "MQTT server URL")
	viperBindFlag("mqtt.server", startCmd.Flags().Lookup("mqtt-server"))

	startCmd.Flags().String("mqtt-username", "", "MQTT username")
	viperBindFlag("mqtt.username", startCmd.Flags().Lookup("mqtt-username"))

	startCmd.Flags().String("mqtt-password", "", "MQTT password")
	viperBindFlag("mqtt.password", startCmd.Flags().Lookup("mqtt-password"))

	startCmd.Flags().String("mqtt-availability-publish-interval", "", "MQTT availability publish interval in seconds")
	viperBindFlag("mqtt.availability-publish-interval", startCmd.Flags().Lookup("mqtt-availability-publish-interval"))

	startCmd.Flags().String("nats-server", "", "NATS server URL")
	viperBindFlag("nats.server", startCmd.Flags().Lookup("nats-server"))

	startCmd.Flags().String("nats-username", "", "NATS username")
	viperBindFlag("nats.username", startCmd.Flags().Lookup("nats-username"))

	startCmd.Flags().String("nats-password", "", "NATS password")
	viperBindFlag("nats.password", startCmd.Flags().Lookup("nats-password"))
}

func startServer() error {
	// config proxmox client from cli
	tokenID := viper.GetString("proxmox.token-id")
	if tokenID == "" {
		logger.Panic("proxmox api token not provided")
	}

	apisecret := viper.GetString("proxmox.token-secret")
	if apisecret == "" {
		logger.Panic("proxmox api secret not provided")
	}

	serverURL := viper.GetString("proxmox.server-url")
	if serverURL == "" {
		logger.Panic("proxmox server URL not provided")
	}

	opts := []proxmox.Opt{
		proxmox.WithAPIToken(tokenID, apisecret),
		proxmox.WithLogger(logger.Desugar()),
	}

	if viper.GetBool("proxmox.insecure") {
		logger.Debugf("insecure enabled")

		opts = append(opts, proxmox.WithInsecure())
	}

	vmids := viper.GetIntSlice("proxmox.vm-ids")
	logger.Debug("vmids: %v", vmids)

	c := proxmox.NewClient(serverURL, opts...)

	// config nats client from cli
	url := viper.GetString("nats.server")
	if url == "" {
		logger.Panic("nats server URL not provided")
	}

	username := viper.GetString("nats.username")
	if username == "" {
		logger.Panic("nats username not provided")
	}

	password := viper.GetString("nats.password")
	if password == "" {
		logger.Panic("nats password not provided")
	}

	apiserveropts := []api.Opt{
		api.WithProxmoxClient(c),
		api.WithLogger(logger.Desugar()),
		api.WithEventsClient(natsclient.NewClient(
			natsclient.WithNATSURL(url),
			natsclient.WithLogger(logger.Desugar()),
			natsclient.WithNATSOpts(
				nats.UserInfo(username, password),
				nats.PingInterval(1*time.Second),
				nats.Timeout(1*time.Second),
			),
		)),
	}

	probeIntervalStr := viper.GetString("proxmox.probe-interval")
	if probeIntervalStr != "" {
		probeInterval, err := time.ParseDuration(probeIntervalStr)
		if err != nil {
			logger.Panicf("failed to parse probe interval: %v", err)
		}

		logger.Debug("probe interval: %v", probeInterval)

		apiserveropts = append(apiserveropts, api.WithProbeInterval(probeInterval))
	}

	availabilityPublishIntervalStr := viper.GetString("mqtt.availability-publish-interval")
	if availabilityPublishIntervalStr != "" {
		availabilityPublishInterval, err := time.ParseDuration(availabilityPublishIntervalStr)
		if err != nil {
			logger.Panicf("failed to parse availability publish interval: %v", err)
		}

		logger.Debug("availability publish interval: %v", availabilityPublishInterval)

		apiserveropts = append(apiserveropts, api.WithAvailabilityPublishInterval(availabilityPublishInterval))
	}

	s := api.NewServer(
		vmids,
		apiserveropts...,
	)

	if err := s.Start(); err != nil {
		logger.Panic(err.Error())
	}

	return nil
}
