// Code generated by mockery v2.12.1. DO NOT EDIT.

package mocks

import (
	hotstuff "github.com/onflow/flow-go/consensus/hotstuff"
	flow "github.com/onflow/flow-go/model/flow"

	mock "github.com/stretchr/testify/mock"

	testing "testing"
)

// Packer is an autogenerated mock type for the Packer type
type Packer struct {
	mock.Mock
}

// Pack provides a mock function with given fields: blockID, sig
func (_m *Packer) Pack(blockID flow.Identifier, sig *hotstuff.BlockSignatureData) ([]flow.Identifier, []byte, error) {
	ret := _m.Called(blockID, sig)

	var r0 []flow.Identifier
	if rf, ok := ret.Get(0).(func(flow.Identifier, *hotstuff.BlockSignatureData) []flow.Identifier); ok {
		r0 = rf(blockID, sig)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]flow.Identifier)
		}
	}

	var r1 []byte
	if rf, ok := ret.Get(1).(func(flow.Identifier, *hotstuff.BlockSignatureData) []byte); ok {
		r1 = rf(blockID, sig)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).([]byte)
		}
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(flow.Identifier, *hotstuff.BlockSignatureData) error); ok {
		r2 = rf(blockID, sig)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// Unpack provides a mock function with given fields: blockID, signerIDs, sigData
func (_m *Packer) Unpack(blockID flow.Identifier, signerIDs []flow.Identifier, sigData []byte) (*hotstuff.BlockSignatureData, error) {
	ret := _m.Called(blockID, signerIDs, sigData)

	var r0 *hotstuff.BlockSignatureData
	if rf, ok := ret.Get(0).(func(flow.Identifier, []flow.Identifier, []byte) *hotstuff.BlockSignatureData); ok {
		r0 = rf(blockID, signerIDs, sigData)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*hotstuff.BlockSignatureData)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(flow.Identifier, []flow.Identifier, []byte) error); ok {
		r1 = rf(blockID, signerIDs, sigData)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewPacker creates a new instance of Packer. It also registers the testing.TB interface on the mock and a cleanup function to assert the mocks expectations.
func NewPacker(t testing.TB) *Packer {
	mock := &Packer{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
