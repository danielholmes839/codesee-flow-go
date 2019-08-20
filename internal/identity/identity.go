// Package identity provides an interface that bundles all relevant information about staked nodes.
// A simple in-memory implementation is included.
package identity

import (
	"fmt"
	"math/big"
	"sort"
)

type NodeRole string

const (
	CollectorRole NodeRole = "collector"
	ConsensusRole NodeRole = "consensus"
	ExecutorRole  NodeRole = "executor"
	VerifierRole  NodeRole = "verifier"
	ObserverRole  NodeRole = "observer"
)

type NodeIdentity interface {
	ID() uint
	Address() string
	Role() NodeRole
	Stake() *big.Int
	Index() uint
}

type Table interface {
	Count() uint
	Nodes() []NodeIdentity

	GetByIndex(uint) (NodeIdentity, error)
	GetByID(uint) (NodeIdentity, error)
	GetByAddress(string) (NodeIdentity, error)

	TotalStake() *big.Int

	FilterByID([]uint) (Table, error)
	FilterByAddress([]string) (Table, error)
	FilterByRole(NodeRole) (Table, error)
	FilterByIndex([]uint) (Table, error)
}

// NodeRecord provides information about one Node that (independent of potential other nodes)
type NodeRecord struct {
	ID      uint
	Address string
	Role    NodeRole
	Stake   *big.Int
}

// NodeRecords is a slice of *NodeRecord which implements sort.Interface
// Sorting is based solely on NodeRecord.ID
type NodeRecords []*NodeRecord

func (ns NodeRecords) Len() int           { return len(ns) }
func (ns NodeRecords) Less(i, j int) bool { return ns[i].ID < ns[j].ID }
func (ns NodeRecords) Swap(i, j int)      { ns[i], ns[j] = ns[j], ns[i] }



// Implementation of NodeIdentity interface
type nodeIdentity struct {
	coreID *NodeRecord
	index  uint
}

func (i nodeIdentity) ID() uint {
	return i.coreID.ID
}

func (i nodeIdentity) Address() string {
	return i.coreID.Address
}

func (i nodeIdentity) Role() NodeRole {
	return i.coreID.Role
}

func (i nodeIdentity) Stake() *big.Int {
	return i.coreID.Stake
}

func (i nodeIdentity) Index() uint {
	return i.index
}

// nodeIdentities is a slice of *nodeIdentity which implements sort.Interface
// Sorting is based solely on nodeIdentity.ID
type nodeIdentities []*nodeIdentity

func (ns nodeIdentities) Len() int           { return len(ns) }
func (ns nodeIdentities) Less(i, j int) bool { return ns[i].coreID.ID < ns[j].coreID.ID }
func (ns nodeIdentities) Swap(i, j int)      { ns[i], ns[j] = ns[j], ns[i] }

// InMemoryIdentityTable is a In-memory implementation of the interface identity.Table
type InMemoryIdentityTable struct {
	nodes      []*nodeIdentity
	addressMap map[string]*nodeIdentity
	idMap      map[uint]*nodeIdentity
}

func (t InMemoryIdentityTable) Count() uint {
	return uint(len(t.nodes))
}

func (t InMemoryIdentityTable) Nodes() []NodeIdentity {
	identities := make([]NodeIdentity, len(t.nodes))
	for i, n := range t.nodes { // converting explicitly from nodeIdentity ti interface NodeIdentity
		identities[i] = n
	}
	return identities
}

func (t InMemoryIdentityTable) GetByIndex(idx uint) (NodeIdentity, error) {
	if int(idx) > len(t.nodes) {
		return nil, &NodeNotFoundError{fmt.Sprint(idx)}
	}
	return t.nodes[idx], nil
}

func (t InMemoryIdentityTable) GetByID(id uint) (NodeIdentity, error) {
	value, found := t.idMap[id]
	if !found {
		return nil, &NodeNotFoundError{fmt.Sprint(id)}
	}
	return value, nil
}

func (t InMemoryIdentityTable) GetByAddress(address string) (NodeIdentity, error) {
	value, found := t.addressMap[address]
	if !found {
		return nil, &NodeNotFoundError{address}
	}
	return value, nil
}

func (t InMemoryIdentityTable) TotalStake() *big.Int {
	s := big.NewInt(0)
	for _, n := range t.nodes {
		s.Add(s, n.Stake())
	}
	return s
}

func (t InMemoryIdentityTable) FilterByID(ids []uint) (Table, error) {
	nodes := make([]*NodeRecord, len(ids))
	missing := []uint{}
	var n *nodeIdentity
	var found bool
	var idx int = 0
	for _, id := range ids {
		n, found = t.idMap[id]
		if found {
			nodes[idx] = n.coreID
			idx++
		} else {
			missing = append(missing, id)
		}
	}
	if len(missing) > 0 {
		return NewInMemoryIdentityTable(nodes[0:idx]), &NodeNotFoundError{fmt.Sprint(missing)}
	}
	return NewInMemoryIdentityTable(nodes[0:idx]), nil
}

func (t InMemoryIdentityTable) FilterByAddress(addresses []string) (Table, error) {
	nodes := make([]*NodeRecord, len(addresses))
	missing := []string{}
	var n *nodeIdentity
	var found bool
	var idx int = 0
	for _, addr := range addresses {
		n, found = t.addressMap[addr]
		if found {
			nodes[idx] = n.coreID
			idx++
		} else {
			missing = append(missing, addr)
		}
	}
	if len(missing) > 0 {
		return NewInMemoryIdentityTable(nodes[0:idx]), &NodeNotFoundError{fmt.Sprint(missing)}
	}
	return NewInMemoryIdentityTable(nodes[0:idx]), nil
}

func (t InMemoryIdentityTable) FilterByRole(role NodeRole) (Table, error) {
	nodes := make([]*NodeRecord, t.Count())
	var idx int = 0
	for _, n := range t.nodes {
		if n.Role() == role {
			nodes[idx] = n.coreID
			idx++
		}
	}
	if idx == 0 {
		return nil, &NodeNotFoundError{fmt.Sprint(role)}
	}
	return NewInMemoryIdentityTable(nodes[0:idx]), nil
}

func (t InMemoryIdentityTable) FilterByIndex(indices []uint) (Table, error) {
	nodes := make([]*NodeRecord, len(indices))
	missing := []uint{}
	var idx int = 0
	_nodeCount := uint(t.Count())
	_allNodes := t.nodes
	for _, i := range indices {
		if i < _nodeCount {
			nodes[idx] = _allNodes[i].coreID
			idx++
		} else {
			missing = append(missing, i)
		}
	}
	if len(missing) > 0 {
		return NewInMemoryIdentityTable(nodes[0:idx]), &NodeNotFoundError{fmt.Sprint(missing)}
	}
	return NewInMemoryIdentityTable(nodes[0:idx]), nil
}

type NodeNotFoundError struct {
	key string
}

func (e *NodeNotFoundError) Error() string {
	return fmt.Sprintf("node with '%s' not found", e.key)
}

func NewInMemoryIdentityTable(nodes []*NodeRecord) *InMemoryIdentityTable {
	nidentities := newSortedNodeIdentities(nodes)

	addressMap := make(map[string]*nodeIdentity)
	idMap := make(map[uint]*nodeIdentity)
	var _last *NodeRecord = nil // reference to previous nodeIdentity to detect duplicates
	for i, n := range nidentities {
		if _last == n.coreID {
			panic("Duplicate NodeRecord not supported")
		}
		_last = n.coreID
		n.index = uint(i)
		addressMap[n.coreID.Address] = n
		idMap[n.coreID.ID] = n
	}

	return &InMemoryIdentityTable{nidentities, addressMap, idMap}
}

// Wraps each NodeRecord as into a nodeIdentity type with default `index=0` and sorts the elements
// Checks for nil elements in
func newSortedNodeIdentities(nodes []*NodeRecord) []*nodeIdentity {
	// While the slice `nodes` is copied, the data in the slice is not sufficient to sort the slice without mutating the underlying array
	// For more details, see https://blog.golang.org/go-slices-usage-and-internals
	nidentities := make([]*nodeIdentity, len(nodes))
	for i, n := range nodes {
		if n == nil {
			panic("NodeRecord cannot be nil")
		}
		nidentities[i] = &nodeIdentity{coreID: n}
	}
	sort.Sort(nodeIdentities(nidentities))
	return nidentities
}
