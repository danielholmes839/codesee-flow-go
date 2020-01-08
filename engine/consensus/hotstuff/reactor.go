package hotstuff

import (
	"github.com/dapperlabs/flow-go/engine/consensus/hotstuff/types"
)

type Reactor interface {
	GetQCForNextBlock(view uint64) *types.QuorumCertificate
	BlocksForView(view uint64) []*types.BlockProposal
	FindBlockProposalByViewAndBlockMRH(view uint64, blockMRH types.MRH) *types.BlockProposal
	FinalizedView() uint64
	IsSafeNode(block *types.BlockProposal) bool

	AddBlock(*types.BlockProposal) // ToDo rename to AddBlock
	ProcessQcFromVotes(*types.QuorumCertificate)

	// IsKnownBlock returns true if the consensus reactor knows the specified block
	IsKnownBlock([]byte, uint64) bool

	// IsProcessingNeeded returns true if consensus reactor should process the specified block
	IsProcessingNeeded([]byte, uint64) bool




	// MakeForkChoice prompts the ForkChoice to generate a fork choice.
	// The fork choice is a qc that should be used for building the primaries block
	MakeForkChoice(viewNumber uint64) *types.QuorumCertificate
}
