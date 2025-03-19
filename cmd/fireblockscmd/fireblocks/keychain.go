package fireblocks

import (
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/utils/crypto/keychain"
	"github.com/ava-labs/avalanchego/utils/crypto/secp256k1"
	"github.com/ava-labs/avalanchego/utils/hashing"
	"github.com/ava-labs/avalanchego/utils/set"

	"github.com/ava-labs/avalanche-cli/pkg/prompts"
)

var (
	_ keychain.Keychain = &FireblocksKeychain{}
	_ keychain.Signer   = &FireblocksSigner{}
)

type FireblocksKeychain struct {
	signer *FireblocksSigner
}

type FireblocksSigner struct {
	sdk *SDK

	account      int
	addressIndex int

	addr ids.ShortID
	mu   sync.Mutex
}

func NewFireblocksKeychain(params *prompts.FireblocksParams) (keychain.Keychain, error) {
	var apiEndpoint string
	if params.UseSandbox {
		apiEndpoint = "https://sandbox-api.fireblocks.io"
	} else {
		apiEndpoint = "https://api.fireblocks.io"
	}

	signer, err := NewFireblocksSigner(apiEndpoint, params.PrivateKeyPath, params.APIKey, params.Account, params.AddressIndex)
	if err != nil {
		return nil, err
	}

	return &FireblocksKeychain{
		signer: signer,
	}, nil
}

// The returned Signer can provide a signature for [addr]
func (fk *FireblocksKeychain) Get(addr ids.ShortID) (keychain.Signer, bool) {
	if fk.signer.Address().Compare(addr) != 0 {
		return nil, false
	}
	return fk.signer, true
}

// Returns the set of addresses for which the accessor keeps an associated
// signer
func (fk *FireblocksKeychain) Addresses() set.Set[ids.ShortID] {
	return set.Of(fk.signer.Address())
}

func NewFireblocksSigner(apiAddr, privateKeyPath, apiKey string, account, addressIndex int) (*FireblocksSigner, error) {
	f, err := os.Open(privateKeyPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	pkBytes, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	return &FireblocksSigner{
		sdk:          NewInstance(pkBytes, apiKey, apiAddr, time.Hour),
		account:      account,
		addressIndex: addressIndex,

		addr: ids.ShortEmpty,
		mu:   sync.Mutex{},
	}, nil
}

func (fs *FireblocksSigner) SignHash(hash []byte) ([]byte, error) {
	sig, _, err := fs.sdk.SignData(fs.account, fs.addressIndex, hash)
	return sig, err
}

func (fs *FireblocksSigner) Sign(data []byte) ([]byte, error) {
	return fs.SignHash(hashing.ComputeHash256(data))
}

func (fs *FireblocksSigner) Address() ids.ShortID {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	if fs.addr.Compare(ids.ShortEmpty) == 0 {
		msg, err := hex.DecodeString("802a5a961895b3f8c6556e31d0960a5778d7135be7d04bbbadd5e406c4bac381")
		if err != nil {
			panic(err)
		}

		rawSignature, rawPublicKey, err := fs.sdk.SignData(fs.account, fs.addressIndex, msg)
		if err != nil {
			panic(err)
		}

		pb, err := secp256k1.ToPublicKey(rawPublicKey)
		if err != nil {
			panic(err)
		}

		fs.addr = pb.Address()
		fmt.Printf("PB1 ShortID %s\n", fs.addr)

		pb2, err := secp256k1.RecoverPublicKeyFromHash(msg, rawSignature)
		if err != nil {
			panic(err)
		}
		pb2b := pb2.Bytes()
		fmt.Printf("PB2 public key %s\n", hex.EncodeToString(pb2b))

		pb2Addr := pb2.Address()
		fmt.Printf("PB2 ShortID %s\n", pb2Addr)
	}

	return fs.addr
}
