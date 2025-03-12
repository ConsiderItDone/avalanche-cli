package fireblocks

import (
	"io"
	"os"
	"time"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/utils/crypto/keychain"
	"github.com/ava-labs/avalanchego/utils/set"
)

var (
	_ keychain.Keychain = &FireblocksKeychain{}
	_ keychain.Signer   = &FireblocksSigner{}
)

type FireblocksKeychain struct {
	pk    string
	ak    string
	vault string
}

type FireblocksSigner struct {
	sdk   *SDK
	vault string
}

func NewFireblocksKeychain(pk, ak, vault string) (*FireblocksKeychain, error) {
	return &FireblocksKeychain{pk, ak, vault}, nil
}

// The returned Signer can provide a signature for [addr]
func (fk *FireblocksKeychain) Get(addr ids.ShortID) (keychain.Signer, bool) {
	signer, err := NewFireblocksSigner(fk.pk, fk.ak, fk.vault)
	if err != nil {
		return nil, false
	}
	return signer, true
}

// Returns the set of addresses for which the accessor keeps an associated
// signer
func (*FireblocksKeychain) Addresses() set.Set[ids.ShortID] {
	s := set.NewSet[ids.ShortID](1)
	s.Add(ids.ShortEmpty)
	return s
}

func NewFireblocksSigner(pk, ak, vault string) (*FireblocksSigner, error) {
	f, err := os.Open(pk)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	pkBytes, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	return &FireblocksSigner{
		NewInstance(pkBytes, ak, "https://sandbox-api.fireblocks.io", time.Hour),
		vault,
	}, nil
}

func (*FireblocksSigner) SignHash([]byte) ([]byte, error) {
	panic("impelement me")
}

func (*FireblocksSigner) Sign([]byte) ([]byte, error) {
	panic("impelement me")
}

func (*FireblocksSigner) Address() ids.ShortID {
	panic("doesn't support")
}
