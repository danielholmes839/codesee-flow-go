package flow

import (
	"fmt"
	"time"
)

// GenesisTime defines the timestamp of the genesis block.
var GenesisTime = time.Date(2018, time.December, 19, 22, 32, 30, 42, time.UTC)

// DefaultProtocolVersion is the default protocol version, indicating none was
// explicitly set during bootstrapping.
const DefaultProtocolVersion = 0

// DefaultTransactionExpiry is the default expiry for transactions, measured in blocks.
// The default value is equivalent to 10 minutes for a 1-second block time.
//
// Let E by the transaction expiry. If a transaction T specifies a reference
// block R with height H, then T may be included in any block B where:
// * R<-*B - meaning B has R as an ancestor, and
// * R.height < B.height <= R.height+E
const DefaultTransactionExpiry = 10 * 60

// DefaultTransactionExpiryBuffer is the default buffer time between a transaction being ingested by a
// collection node and being included in a collection and block.
const DefaultTransactionExpiryBuffer = 30

// DefaultMaxTransactionGasLimit is the default maximum value for the transaction gas limit.
const DefaultMaxTransactionGasLimit = 9999

// DefaultMaxTransactionByteSize is the default maximum transaction byte size. (~1.5MB)
const DefaultMaxTransactionByteSize = 1_500_000

// DefaultMaxCollectionByteSize is the default maximum value for a collection byte size.
const DefaultMaxCollectionByteSize = 3_000_000 // ~3MB. This is should always be higher than the limit on single tx size.

// DefaultMaxCollectionTotalGas is the default maximum value for total gas allowed to be included in a collection.
const DefaultMaxCollectionTotalGas = 10_000_000 // 10M

// DefaultMaxCollectionSize is the default maximum number of transactions allowed inside a collection.
const DefaultMaxCollectionSize = 100

// DefaultMaxAddressIndex is the default for the maximum address index allowed to be acceptable by collection and acccess nodes.
const DefaultMaxAddressIndex = 20_000_000

// DefaultValueLogGCFrequency is the default frequency in blocks that we call the
// badger value log GC. Equivalent to 10 mins for a 1 second block time
const DefaultValueLogGCFrequency = 10 * 60

// DomainTagLength is set to 32 bytes.
//
// Signatures on Flow that needs to be scoped to a certain domain need to
// have the same length in order to avoid tag collision issues, when prefixing the
// message to sign.
const DomainTagLength = 32

const TransactionTagString = "FLOW-V0.0-transaction"

// TransactionDomainTag is the prefix of all signed transaction payloads.
//
// The tag is the string `TransactionTagString` encoded as UTF-8 bytes,
// right padded to a total length of 32 bytes.
var TransactionDomainTag = paddedDomainTag(TransactionTagString)

// paddedDomainTag padds string tags to form the actuatl domain separation tag used for signing
// and verifiying.
//
// A domain tag is encoded as UTF-8 bytes, right padded to a total length of 32 bytes.
func paddedDomainTag(s string) [DomainTagLength]byte {
	var tag [DomainTagLength]byte

	if len(s) > DomainTagLength {
		panic(fmt.Sprintf("domain tag %s cannot be longer than %d characters", s, DomainTagLength))
	}

	copy(tag[:], s)

	return tag
}
