package fireblocks

import (
	"encoding/hex"
	"io"
	"os"
	"sync"
	"time"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/utils/crypto/keychain"
	"github.com/ava-labs/avalanchego/utils/crypto/secp256k1"
	"github.com/ava-labs/avalanchego/utils/hashing"
	"github.com/ava-labs/avalanchego/utils/set"
)

var (
	_ keychain.Keychain = &FireblocksKeychain{}
	_ keychain.Signer   = &FireblocksSigner{}
)

type FireblocksKeychain struct {
	signer *FireblocksSigner
}

type FireblocksSigner struct {
	sdk     *SDK
	vaultid string
	assetid string

	addr ids.ShortID
	mu   sync.Mutex
}

func NewFireblocksKeychain(pk, ak, vaultid, assetid string) (*FireblocksKeychain, error) {
	signer, err := NewFireblocksSigner(pk, ak, vaultid, assetid)
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
	s := set.NewSet[ids.ShortID](1)
	s.Add(fk.signer.Address())
	return s
}

func NewFireblocksSigner(pk, ak, vaultid, assetid string) (*FireblocksSigner, error) {
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
		sdk:     NewInstance(pkBytes, ak, "https://sandbox-api.fireblocks.io", time.Hour),
		vaultid: vaultid,
		assetid: assetid,

		addr: ids.ShortEmpty,
		mu:   sync.Mutex{},
	}, nil
}

func (fs *FireblocksSigner) SignHash(hash []byte) ([]byte, error) {
	sig, _, err := fs.sdk.SignData(fs.vaultid, fs.assetid, hash)
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

		_, rawpb, err := fs.sdk.SignData(fs.vaultid, fs.assetid, msg)
		if err != nil {
			panic(err)
		}

		pb, err := secp256k1.ToPublicKey(rawpb)
		if err != nil {
			panic(err)
		}

		fs.addr = pb.Address()
	}

	return fs.addr
}
