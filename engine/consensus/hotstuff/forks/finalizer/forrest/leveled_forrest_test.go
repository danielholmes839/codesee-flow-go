package forrest

import (
	"testing"

	"github.com/dapperlabs/flow-go/engine/consensus/hotstuff/forks/finalizer/forrest/mock"
	"github.com/dapperlabs/flow-go/model/flow"
	"github.com/stretchr/testify/assert"
)

// ~~~~~~~~~~~~~~~~~~~~~ Mock implementation for Vertex ~~~~~~~~~~~~~~~~~~~~~ //
type VertexMock struct {
	id    flow.Identifier
	level uint64

	parentId    flow.Identifier
	parentLevel uint64
}

func (v *VertexMock) VertexID() flow.Identifier         { return v.id }
func (v *VertexMock) Level() uint64            { return v.level }
func (v *VertexMock) Parent() (flow.Identifier, uint64) { return v.parentId, v.parentLevel }

func NewVertexMock(vertexId string, vertexLevel uint64, parentId string, parentLevel uint64) *mock.Vertex {
	v := &mock.Vertex{}
	v.On("VertexID").Return(string2Identifyer(vertexId))
	v.On("Level").Return(vertexLevel)
	v.On("Parent").Return(string2Identifyer(parentId), parentLevel)
	return v
}

var TestVertices = map[string]*mock.Vertex{
	"A": NewVertexMock("A", 3, "Genesis", 0),
	"B": NewVertexMock("B", 1, "Genesis", 0),
	"C": NewVertexMock("C", 2, "B", 1),
	"D": NewVertexMock("D", 3, "C", 2),
	"W": NewVertexMock("W", 4, "Missing1", 3),
	"X": NewVertexMock("X", 5, "Missing2", 4),
	"Y": NewVertexMock("Y", 6, "X", 5),
	"Z": NewVertexMock("Z", 6, "X", 5),
}

// ~~~~~~~~~~~~~~~~~~~~~~~~ Tests for VertexIterator ~~~~~~~~~~~~~~~~~~~~~~~~ //

// TestVertexIterator tests VertexIterator on non-empty list
func TestVertexIterator(t *testing.T) {
	var vl VertexList
	for i := 1; i <= 10; i++ {
		b := VertexMock{level: uint64(i)}
		c := vertexContainer{vertex: &b}
		vl = append(vl, &c)
	}

	i := 0
	for I := newVertexIterator(vl); I.HasNext(); {
		assert.Equal(t, vl[i].vertex, I.NextVertex())
		i++
	}
}

// TestVertexIteratorOnEmpty tests VertexIterator on non-empty and nil lists
func TestVertexIteratorOnEmpty(t *testing.T) {
	I := newVertexIterator(nil)
	assert.False(t, I.HasNext())
	assert.True(t, I.NextVertex() == nil)

	I = newVertexIterator(VertexList{})
	assert.False(t, I.HasNext())
	assert.True(t, I.NextVertex() == nil)
}

// ~~~~~~~~~~~~~~~~~~~~~~~~ Tests for LeveledForrest ~~~~~~~~~~~~~~~~~~~~~~~~ //

// TestLeveledForrest_AddVertex tests that Vertex can be added twice without problems
func TestLeveledForrest_AddVertex(t *testing.T) {
	F := NewLeveledForrest()
	v := NewVertexMock("A", 3, "Genesis", 0)
	if err := F.VerifyVertex(v); err != nil {
		assert.Fail(t, err.Error())
	}
	F.AddVertex(v)
	assert.True(t, F.HasVertex(string2Identifyer("A")))

	// Adding Vertex twice should be fine
	v = NewVertexMock("A", 3, "Genesis", 0)
	if err := F.VerifyVertex(v); err != nil {
		assert.Fail(t, err.Error())
	}
	F.AddVertex(v)
	assert.True(t, F.HasVertex(string2Identifyer("A")))
}

// TestLeveledForrest_AcceptingGenesis checks that Levelled forrest accepts vertices
// whose level are at LeveledForrest.LowestLevel without requiring the parent.
// we test this by having the mock.vertex.Parent() panic
func TestLeveledForrest_AcceptingGenesis(t *testing.T) {
	// LeveledForrest.LowestLevel on initial conditions
	F := populateNewForrest(t)
	v1 := &mock.Vertex{}
	v1.On("VertexID").Return(string2Identifyer("Root-Vertex-A_@Level0"))
	v1.On("Level").Return(uint64(0))
	v1.On("Parent").Return(func() (flow.Identifier, uint64) { panic("Parent() should not have been called") })
	assert.NotPanics(t, func() { F.AddVertex(v1) })

	v2 := &mock.Vertex{}
	v2.On("VertexID").Return(string2Identifyer("Root-Vertex-B_@Level0"))
	v2.On("Level").Return(uint64(0))
	v2.On("Parent").Return(func() (flow.Identifier, uint64) { panic("Parent() should not have been called") })
	assert.NotPanics(t, func() { F.AddVertex(v2) })
	assert.NotPanics(t, func() { F.AddVertex(v2) })

	F = populateNewForrest(t)
	F.PruneAtLevel(7) // LeveledForrest.LowestLevel on initial conditions
	v3 := &mock.Vertex{}
	v3.On("VertexID").Return(string2Identifyer("Root-Vertex-A_@Level8"))
	v3.On("Level").Return(uint64(8))
	v3.On("Parent").Return(func() (flow.Identifier, uint64) { panic("Parent() should not have been called") })
	assert.NotPanics(t, func() { F.AddVertex(v3) })

	v4 := &mock.Vertex{}
	v4.On("VertexID").Return(string2Identifyer("Root-Vertex-B_@Level8"))
	v4.On("Level").Return(uint64(8))
	v4.On("Parent").Return(func() (flow.Identifier, uint64) { panic("Parent() should not have been called") })
	assert.NotPanics(t, func() { F.AddVertex(v4) })
	assert.NotPanics(t, func() { F.AddVertex(v4) })
}

// TestLeveledForrest_VerifyVertex checks that invalid Vertices are detected.
// with an ID identical to a known reference BUT whose level is not consistent with the reference
func TestLeveledForrest_VerifyVertex(t *testing.T) {
	F := populateNewForrest(t)

	// KNOWN vertex but with wrong level number
	err := F.VerifyVertex(NewVertexMock("D", 10, "C", 2))
	assert.True(t, err != nil, err.Error())

	// KNOWN vertex whose PARENT references a known vertex but with mismatching level
	err = F.VerifyVertex(NewVertexMock("D", 10, "C", 10))
	assert.True(t, err != nil, err.Error())

	// adding unknown vertex whose PARENT references a known vertex but with mismatching level
	err = F.VerifyVertex(NewVertexMock("F", 4, "Genesis", 10))
	assert.True(t, err != nil, err.Error())
}

// TestLeveledForrest_HasVertex test that vertices as correctly reported as contained in forrest
// NOTE: We consider a vertex added only if it has been directly added through the AddVertex method.
//       Vertices that references bvy known vertices but have not themselves are considered to be not in the tree.
func TestLeveledForrest_HasVertex(t *testing.T) {
	F := populateNewForrest(t)
	assert.True(t, F.HasVertex(string2Identifyer("A")))
	assert.True(t, F.HasVertex(string2Identifyer("B")))
	assert.True(t, F.HasVertex(string2Identifyer("X")))

	assert.False(t, F.HasVertex(string2Identifyer("Genesis")))     // Genesis block never directly added (only referenced) => unknown
	assert.False(t, F.HasVertex(string2Identifyer("NotYetAdded"))) // Block never mentioned before
}

// TestLeveledForrest_GetChildren tests that children are returned properly
func TestLeveledForrest_GetChildren(t *testing.T) {
	F := populateNewForrest(t)

	// testing children for Block that is contained in Tree
	it := F.GetChildren(string2Identifyer("X"))
	expectedChildren := []*mock.Vertex{
		TestVertices["Y"],
		TestVertices["Z"],
	}
	assert.ElementsMatch(t, expectedChildren, children2List(&it))

	// testing children for referenced Block that is NOT contained in Tree
	it = F.GetChildren(string2Identifyer("Genesis"))
	expectedChildren = []*mock.Vertex{
		TestVertices["A"],
		TestVertices["B"],
	}
	assert.ElementsMatch(t, expectedChildren, children2List(&it))

	// testing children for Block that is contained in Tree but no children are known
	it = F.GetChildren(string2Identifyer("D"))
	assert.False(t, it.HasNext())
}

// TestLeveledForrest_GetNumberOfChildren tests that children are returned properly
func TestLeveledForrest_GetNumberOfChildren(t *testing.T) {
	F := populateNewForrest(t)

	// testing children for Block that is contained in Tree
	assert.Equal(t, 2, F.GetNumberOfChildren(string2Identifyer("X")))

	// testing children for referenced Block that is NOT contained in Tree
	assert.Equal(t, 2, F.GetNumberOfChildren(string2Identifyer("Genesis")))

	// testing children for Block that is contained in Tree but no children are known
	assert.Equal(t, 0, F.GetNumberOfChildren(string2Identifyer("D")))
}

// TestLeveledForrest_GetVerticesAtLevel tests that Vertex blob is returned properly
func TestLeveledForrest_GetVerticesAtLevel(t *testing.T) {
	F := populateNewForrest(t)

	// testing vertices for level that are contained in Tree
	it := F.GetVerticesAtLevel(6)
	expectedChildren := []*mock.Vertex{
		TestVertices["Y"],
		TestVertices["Z"],
	}
	assert.ElementsMatch(t, expectedChildren, children2List(&it))

	// testing vertices for level that are not in Tree but referenced by vertices in the Tree
	it = F.GetVerticesAtLevel(0)
	assert.ElementsMatch(t, []*mock.Vertex{}, children2List(&it))

	// testing vertices for level with a mixture of referenced but unknown vertices and known vertices
	it = F.GetVerticesAtLevel(4)
	expectedChildren = []*mock.Vertex{
		TestVertices["W"],
	}
	assert.ElementsMatch(t, expectedChildren, children2List(&it))

	// testing vertices for level that are neither in Tree nor referenced by vertices in the Tree
	it = F.GetVerticesAtLevel(100000)
	assert.ElementsMatch(t, []*mock.Vertex{}, children2List(&it))
}

// TestLeveledForrest_GetNumberOfVerticesAtLevel tests that the number of vertices at a specified level is reported correctly. 
func TestLeveledForrest_GetNumberOfVerticesAtLevel(t *testing.T) {
	F := populateNewForrest(t)

	// testing vertices for level that are contained in Tree
	assert.Equal(t, 2, F.GetNumberOfVerticesAtLevel(6))

	// testing vertices for level that are not in Tree but referenced by vertices in the Tree
	assert.Equal(t, 0, F.GetNumberOfVerticesAtLevel(0))

	// testing vertices for level with a mixture of referenced but unknown vertices and known vertices
	assert.Equal(t, 1, F.GetNumberOfVerticesAtLevel(4))

	// testing vertices for level that are neither in Tree nor referenced by vertices in the Tree
	assert.Equal(t, 0, F.GetNumberOfVerticesAtLevel(100000))
}

// TestLeveledForrest_GetVertex tests that Vertex blob is returned properly
func TestLeveledForrest_GetVertex(t *testing.T) {
	F := populateNewForrest(t)
	v, exists := F.GetVertex(string2Identifyer("D"))
	assert.Equal(t, TestVertices["D"], v)
	assert.True(t, exists)

	v, exists = F.GetVertex(string2Identifyer("X"))
	assert.Equal(t, TestVertices["X"], v)
	assert.True(t, exists)

	v, exists = F.GetVertex(string2Identifyer("Genesis"))
	assert.Equal(t, (Vertex)(nil), v)
	assert.False(t, exists)
}

// TestLeveledForrest_GetVertex tests that Vertex blob is returned properly
func TestLeveledForrest_PruneAtLevel(t *testing.T) {
	F := populateNewForrest(t)
	err := F.PruneAtLevel(0)
	assert.False(t, err != nil)
	assert.False(t, F.HasVertex(string2Identifyer("Genesis")))
	assert.True(t, F.HasVertex(string2Identifyer("A")))
	assert.True(t, F.HasVertex(string2Identifyer("B")))
	assert.True(t, F.HasVertex(string2Identifyer("C")))
	assert.True(t, F.HasVertex(string2Identifyer("D")))
	assert.True(t, F.HasVertex(string2Identifyer("X")))
	assert.True(t, F.HasVertex(string2Identifyer("Y")))
	assert.True(t, F.HasVertex(string2Identifyer("Z")))

	err = F.PruneAtLevel(2)
	assert.False(t, err != nil)
	assert.False(t, F.HasVertex(string2Identifyer("Genesis")))
	assert.True(t, F.HasVertex(string2Identifyer("A")))
	assert.False(t, F.HasVertex(string2Identifyer("B")))
	assert.False(t, F.HasVertex(string2Identifyer("C")))
	assert.True(t, F.HasVertex(string2Identifyer("D")))
	assert.True(t, F.HasVertex(string2Identifyer("X")))
	assert.True(t, F.HasVertex(string2Identifyer("Y")))
	assert.True(t, F.HasVertex(string2Identifyer("Z")))

	err = F.PruneAtLevel(5)
	assert.False(t, err != nil)
	assert.False(t, F.HasVertex(string2Identifyer("Genesis")))
	assert.False(t, F.HasVertex(string2Identifyer("A")))
	assert.False(t, F.HasVertex(string2Identifyer("B")))
	assert.False(t, F.HasVertex(string2Identifyer("C")))
	assert.False(t, F.HasVertex(string2Identifyer("D")))
	assert.False(t, F.HasVertex(string2Identifyer("X")))
	assert.True(t, F.HasVertex(string2Identifyer("Y")))
	assert.True(t, F.HasVertex(string2Identifyer("Z")))

	// pruning at same level repeatedly should be fine
	err = F.PruneAtLevel(5)
	assert.False(t, err != nil)
	assert.True(t, F.HasVertex(string2Identifyer("Y")))
	assert.True(t, F.HasVertex(string2Identifyer("Z")))

	// checking that pruning at lower level than what is already pruned results in error
	err = F.PruneAtLevel(4)
	assert.True(t, err != nil)
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~ Helper Functions ~~~~~~~~~~~~~~~~~~~~~~~~~~~~ //

func populateNewForrest(t *testing.T) *LeveledForrest {
	F := NewLeveledForrest()
	for _, v := range TestVertices {
		if err := F.VerifyVertex(v); err != nil {
			assert.Fail(t, err.Error())
		}
		F.AddVertex(v)
	}
	return F
}

func children2List(it *VertexIterator) []*mock.Vertex {
	l := []*mock.Vertex{}
	for it.HasNext() {
		// Vertex interface is implemented by mock.Vertex POINTER!
		// Hence, the concrete type is *mock.Vertex
		v := it.NextVertex().(*mock.Vertex)
		l = append(l, v)
	}
	return l
}

func string2Identifyer(s string) flow.Identifier {
	var identifier flow.Identifier
	copy(identifier[:], []byte(s))
	return identifier
}
