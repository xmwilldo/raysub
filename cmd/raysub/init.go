package main

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/xmwilldo/v2ray-sub/cmd/raysub/registry"
)

var (
	initCommand = &cobra.Command{
		Use:   "init",
		Short: "",
		Long:  "",
		RunE:  initRun,
		Args:  cobra.MinimumNArgs(1),
	}
)

func init() {
	registry.Commands = append(registry.Commands, registry.CliCommad{
		Command: initCommand,
	})
	initFlags(initCommand.Flags())
}

func initFlags(flags *pflag.FlagSet) {

}

func initRun(cmd *cobra.Command, args []string) error {

	subscriptionURL := args[0]
	viper.Set("subscriptionUrl", subscriptionURL)
	viper.Set("v2rayConfigPath", "/etc/v2ray/config.json")
	if err := viper.WriteConfig(); err != nil {
		return err
	}

	return nil
}
