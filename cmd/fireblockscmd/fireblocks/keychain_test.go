package fireblocks

import "testing"

func TestSigner(t *testing.T) {
	signer, err := NewFireblocksSigner("https://sandbox-api.fireblocks.io", "/Users/n0cte/Downloads/fireblocks_secret_editor_e4fafe6f-742f-423c-b5fa-2af197e932d8.key", "e4fafe6f-742f-423c-b5fa-2af197e932d8", "219", "AVAXTEST")
	if err != nil {
		t.Fatal(err)
	}
	address := signer.Address()
	straddr := address.String()
	t.Logf("Signer: %s %s", address, straddr)
}
