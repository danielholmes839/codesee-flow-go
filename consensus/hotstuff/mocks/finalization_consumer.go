// Code generated by mockery v2.12.1. DO NOT EDIT.

package mocks

import (
	testing "testing"

	model "github.com/onflow/flow-go/consensus/hotstuff/model"
	mock "github.com/stretchr/testify/mock"
)

// FinalizationConsumer is an autogenerated mock type for the FinalizationConsumer type
type FinalizationConsumer struct {
	mock.Mock
}

// OnBlockIncorporated provides a mock function with given fields: _a0
func (_m *FinalizationConsumer) OnBlockIncorporated(_a0 *model.Block) {
	_m.Called(_a0)
}

// OnDoubleProposeDetected provides a mock function with given fields: _a0, _a1
func (_m *FinalizationConsumer) OnDoubleProposeDetected(_a0 *model.Block, _a1 *model.Block) {
	_m.Called(_a0, _a1)
}

// OnFinalizedBlock provides a mock function with given fields: _a0
func (_m *FinalizationConsumer) OnFinalizedBlock(_a0 *model.Block) {
	_m.Called(_a0)
}

// NewFinalizationConsumer creates a new instance of FinalizationConsumer. It also registers the testing.TB interface on the mock and a cleanup function to assert the mocks expectations.
func NewFinalizationConsumer(t testing.TB) *FinalizationConsumer {
	mock := &FinalizationConsumer{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
