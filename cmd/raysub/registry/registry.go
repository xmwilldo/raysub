package registry

import "github.com/spf13/cobra"

type CliCommad struct {
	Command *cobra.Command
	Parent  *cobra.Command
}

var (
	Commands []CliCommad
)
