package fireblockscmd

import (
	"github.com/ava-labs/avalanche-cli/pkg/application"
	"github.com/ava-labs/avalanche-cli/pkg/cobrautils"
	"github.com/spf13/cobra"
)

var app *application.Avalanche

func NewCmd(injectedApp *application.Avalanche) *cobra.Command {
	app = injectedApp

	cmd := &cobra.Command{
		Use:   "fireblocks",
		Short: "Fireblocks helper functions",
		RunE:  cobrautils.CommandSuiteUsage,
	}

	cmd.AddCommand(newAddressCmd())
	return cmd
}
