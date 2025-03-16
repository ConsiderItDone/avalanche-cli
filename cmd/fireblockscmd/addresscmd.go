package fireblockscmd

import (
	"fmt"

	"github.com/ava-labs/avalanchego/utils/formatting/address"
	"github.com/spf13/cobra"

	"github.com/ava-labs/avalanche-cli/cmd/fireblockscmd/fireblocks"
	"github.com/ava-labs/avalanche-cli/pkg/cobrautils"
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
	cmd.Flags().StringVar(&priKey, "private-key", "/srv/landslide/fireblocks_secret_editor_e4fafe6f-742f-423c-b5fa-2af197e932d8.key", "absolute path to fireblocks api private key")
	cmd.Flags().StringVar(&apiKey, "api-key", "e4fafe6f-742f-423c-b5fa-2af197e932d8", "fireblocks api key")
	cmd.Flags().StringVar(&vaultId, "vault-id", "220", "fireblocks vault id")
	cmd.Flags().StringVar(&assetId, "asset-id", "BTC_TEST", "fireblocks asset id")
	cmd.Flags().StringVar(&chainAlias, "avalanche-chain-alias", "P", "avalanche network alias asset id")
	cmd.Flags().StringVar(&chainHrp, "avalanche-chain-hrp", "fuji", "avalanche network hrp")
	return cmd
}

func addresscmd(_ *cobra.Command, _ []string) error {
	keychain, err := fireblocks.PromptFireblocks(app.Prompt)
	if err != nil {
		return err
	}
	addresses := keychain.Addresses()
	addr, exists := addresses.Peek()
	if !exists {
		return fmt.Errorf("no address")
	}
	pChainAddress, err := address.Format(chainAlias, chainHrp, addr.Bytes())
	if err != nil {
		return err
	}
	fmt.Printf("Fireblocks address: %s\n", pChainAddress)
	return nil
}
