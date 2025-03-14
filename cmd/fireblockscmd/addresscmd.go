package fireblockscmd

import (
	"fmt"

	"github.com/ava-labs/avalanche-cli/cmd/fireblockscmd/fireblocks"
	"github.com/ava-labs/avalanche-cli/pkg/cobrautils"
	"github.com/ava-labs/avalanchego/utils/formatting/address"
	"github.com/spf13/cobra"
)

var (
	apiAddr    string
	priKey     string
	apiKey     string
	vaultId    string
	assetId    string
	chainAlias string
	chainHrp   string
)

func newAddressCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:    "address",
		Short:  "show fireblocks avalanche short id",
		RunE:   addresscmd,
		Args:   cobrautils.ExactArgs(0),
		Hidden: false,
	}
	cmd.Flags().StringVar(&apiAddr, "api-addr", "https://sandbox-api.fireblocks.io", "fireblocks api address")
	cmd.Flags().StringVar(&priKey, "private-key", "/Users/n0cte/Downloads/fireblocks_secret_editor_e4fafe6f-742f-423c-b5fa-2af197e932d8.key", "absolute path to fireblocks api private key")
	cmd.Flags().StringVar(&apiKey, "api-key", "e4fafe6f-742f-423c-b5fa-2af197e932d8", "fireblocks api key")
	cmd.Flags().StringVar(&vaultId, "vault-id", "219", "fireblocks vault id")
	cmd.Flags().StringVar(&assetId, "asset-id", "AVAXTEST", "fireblocks asset id")
	cmd.Flags().StringVar(&chainAlias, "avalanche-chain-alias", "P", "avalanche network alias asset id")
	cmd.Flags().StringVar(&chainHrp, "avalanche-chain-hrp", "fuji", "avalanche network hrp")
	return cmd
}

func addresscmd(_ *cobra.Command, _ []string) error {
	signer, err := fireblocks.NewFireblocksSigner(apiAddr, priKey, apiKey, vaultId, assetId)
	if err != nil {
		return err
	}
	addr, err := address.Format(chainAlias, chainHrp, signer.Address().Bytes())
	if err != nil {
		return err
	}
	fmt.Printf("Fireblocks address: %s\n", addr)
	return nil
}
