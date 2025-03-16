package fireblocks

import (
	"strings"

	"github.com/ava-labs/avalanche-cli/pkg/prompts"
)

func PromptFireblocks(prompt prompts.Prompter) (*FireblocksKeychain, error) {
	useSandbox, err := chooseSandboxOrProd(prompt)
	if err != nil {
		return nil, err
	}

	privateKeyPath, err := prompt.CaptureString("Fireblocks private key path")
	if err != nil {
		return nil, err
	}
	privateKeyPath = strings.TrimSpace(privateKeyPath)

	apiKey, err := prompt.CaptureString("Fireblocks api key")
	if err != nil {
		return nil, err
	}

	account, err := prompt.CaptureInt("Fireblocks bip44 account", func(n int) error {
		return nil
	})

	addressIndex, err := prompt.CaptureInt("Fireblocks bip44 address index", func(n int) error {
		return nil
	})

	var apiEndpoint string
	if useSandbox {
		apiEndpoint = "https://sandbox-api.fireblocks.io"
	} else {
		apiEndpoint = "https://api.fireblocks.io"
	}

	fireblocksKc, err := NewFireblocksKeychain(apiEndpoint, privateKeyPath, apiKey, account, addressIndex)
	if err != nil {
		return nil, err
	}

	return fireblocksKc, nil
}

// chooseSandboxOrProd returns true if Sandbox environment is selected
func chooseSandboxOrProd(prompt prompts.Prompter) (bool, error) {
	const (
		sandboxOption = "Sandbox"
		prodOption    = "Production"
	)
	option, err := prompt.CaptureList("What Fireblocks environment should be used?", []string{sandboxOption, prodOption})
	if err != nil {
		return false, err
	}
	return option == sandboxOption, nil
}
