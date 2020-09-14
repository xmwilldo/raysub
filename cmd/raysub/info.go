package main

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/xmwilldo/v2ray-sub/cmd/raysub/registry"
)

var (
	infoCommand = &cobra.Command{
		Use:   "info",
		Short: "",
		Long:  "",
		RunE:  info,
	}
)

func init() {
	registry.Commands = append(registry.Commands, registry.CliCommad{
		Command: infoCommand,
	})
	infoFlags(infoCommand.Flags())
}

func infoFlags(flags *pflag.FlagSet) {

}

func info(cmd *cobra.Command, args []string) error {
	// get some info
	log.Println(args)
	return nil
}
