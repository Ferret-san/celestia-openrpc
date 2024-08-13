package share

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"hash"

	"github.com/celestiaorg/celestia-openrpc/types/appconsts"
	"github.com/celestiaorg/nmt"
)

type NamespacedRow struct {
	Shares []Share    `json:"shares"`
	Proof  *nmt.Proof `json:"proof"`
}

// NamespacedShares represents all shares with proofs within a specific namespace of an EDS.
type NamespacedShares []NamespacedRow

// DefaultRSMT2DCodec sets the default rsmt2d.Codec for shares.
var DefaultRSMT2DCodec = appconsts.DefaultCodec

const (
	// Size is a system-wide size of a share, including both data and namespace GetNamespace
	Size = appconsts.ShareSize
)

// MaxSquareSize is currently the maximum size supported for unerasured data in
// rsmt2d.ExtendedDataSquare.
var MaxSquareSize = appconsts.SquareSizeUpperBound(appconsts.LatestVersion)

// Share contains the raw share data without the corresponding namespace.
// NOTE: Alias for the byte is chosen to keep maximal compatibility, especially with rsmt2d.
// Ideally, we should define reusable type elsewhere and make everyone(Core, rsmt2d, ipld) to rely
// on it.
type Share = []byte

// GetNamespace slices Namespace out of the Share.
func GetNamespace(s Share) Namespace {
	return s[:appconsts.NamespaceSize]
}

// GetData slices out data of the Share.
func GetData(s Share) []byte {
	return s[appconsts.NamespaceSize:]
}

// DataHash is a representation of the Root hash.
type DataHash []byte

func (dh DataHash) Validate() error {
	if len(dh) != 32 {
		return fmt.Errorf("invalid hash size, expected 32, got %d", len(dh))
	}
	return nil
}

func (dh DataHash) String() string {
	return fmt.Sprintf("%X", []byte(dh))
}

// MustDataHashFromString converts a hex string to a valid datahash.
func MustDataHashFromString(datahash string) DataHash {
	dh, err := hex.DecodeString(datahash)
	if err != nil {
		panic(fmt.Sprintf("datahash conversion: passed string was not valid hex: %s", datahash))
	}
	err = DataHash(dh).Validate()
	if err != nil {
		panic(fmt.Sprintf("datahash validation: passed hex string failed: %s", err))
	}
	return dh
}

// NewSHA256Hasher returns a new instance of a SHA-256 hasher.
func NewSHA256Hasher() hash.Hash {
	return sha256.New()
}
