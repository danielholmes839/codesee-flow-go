// Code generated by mockery v1.0.0. DO NOT EDIT.

package mock

import flow "github.com/dapperlabs/flow-go/model/flow"
import mock "github.com/stretchr/testify/mock"
import virtualmachine "github.com/dapperlabs/flow-go/engine/execution/computation/virtualmachine"

// VirtualMachine is an autogenerated mock type for the VirtualMachine type
type VirtualMachine struct {
	mock.Mock
}

// ASTCache provides a mock function with given fields:
func (_m *VirtualMachine) ASTCache() virtualmachine.ASTCache {
	ret := _m.Called()

	var r0 virtualmachine.ASTCache
	if rf, ok := ret.Get(0).(func() virtualmachine.ASTCache); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(virtualmachine.ASTCache)
		}
	}

	return r0
}

// NewBlockContext provides a mock function with given fields: b, blocks
func (_m *VirtualMachine) NewBlockContext(b *flow.Header, blocks virtualmachine.Blocks) virtualmachine.BlockContext {
	ret := _m.Called(b, blocks)

	var r0 virtualmachine.BlockContext
	if rf, ok := ret.Get(0).(func(*flow.Header, virtualmachine.Blocks) virtualmachine.BlockContext); ok {
		r0 = rf(b, blocks)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(virtualmachine.BlockContext)
		}
	}

	return r0
}
