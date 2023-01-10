package postgres

import (
	blockstore "github.com/ipfs/go-ipfs-blockstore"
	dshelp "github.com/ipfs/go-ipfs-ds-help"
	_ "github.com/lib/pq" //postgres driver
	"github.com/multiformats/go-multihash"
)

// MultihashKeyFromKeccak256 converts keccak256 hash bytes into a blockstore-prefixed multihash db key string
func MultihashKeyFromKeccak256(h []byte) (string, error) {
	mh, err := multihash.Encode(h, multihash.KECCAK_256)
	if err != nil {
		return "", err
	}
	dbKey := dshelp.MultihashToDsKey(mh)
	return blockstore.BlockPrefix.String() + dbKey.String(), nil
}
