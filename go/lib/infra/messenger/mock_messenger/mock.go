// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/scionproto/scion/go/lib/infra/messenger (interfaces: LocalSVCRouter,Resolver)

// Package mock_messenger is a generated GoMock package.
package mock_messenger

import (
	context "context"
	net "net"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	addr "github.com/scionproto/scion/go/lib/addr"
	snet "github.com/scionproto/scion/go/lib/snet"
	svc "github.com/scionproto/scion/go/lib/svc"
)

// MockLocalSVCRouter is a mock of LocalSVCRouter interface.
type MockLocalSVCRouter struct {
	ctrl     *gomock.Controller
	recorder *MockLocalSVCRouterMockRecorder
}

// MockLocalSVCRouterMockRecorder is the mock recorder for MockLocalSVCRouter.
type MockLocalSVCRouterMockRecorder struct {
	mock *MockLocalSVCRouter
}

// NewMockLocalSVCRouter creates a new mock instance.
func NewMockLocalSVCRouter(ctrl *gomock.Controller) *MockLocalSVCRouter {
	mock := &MockLocalSVCRouter{ctrl: ctrl}
	mock.recorder = &MockLocalSVCRouterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockLocalSVCRouter) EXPECT() *MockLocalSVCRouterMockRecorder {
	return m.recorder
}

// GetUnderlay mocks base method.
func (m *MockLocalSVCRouter) GetUnderlay(arg0 addr.HostSVC) (*net.UDPAddr, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUnderlay", arg0)
	ret0, _ := ret[0].(*net.UDPAddr)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUnderlay indicates an expected call of GetUnderlay.
func (mr *MockLocalSVCRouterMockRecorder) GetUnderlay(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUnderlay", reflect.TypeOf((*MockLocalSVCRouter)(nil).GetUnderlay), arg0)
}

// MockResolver is a mock of Resolver interface.
type MockResolver struct {
	ctrl     *gomock.Controller
	recorder *MockResolverMockRecorder
}

// MockResolverMockRecorder is the mock recorder for MockResolver.
type MockResolverMockRecorder struct {
	mock *MockResolver
}

// NewMockResolver creates a new mock instance.
func NewMockResolver(ctrl *gomock.Controller) *MockResolver {
	mock := &MockResolver{ctrl: ctrl}
	mock.recorder = &MockResolverMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockResolver) EXPECT() *MockResolverMockRecorder {
	return m.recorder
}

// LookupSVC mocks base method.
func (m *MockResolver) LookupSVC(arg0 context.Context, arg1 snet.Path, arg2 addr.HostSVC) (*svc.Reply, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LookupSVC", arg0, arg1, arg2)
	ret0, _ := ret[0].(*svc.Reply)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LookupSVC indicates an expected call of LookupSVC.
func (mr *MockResolverMockRecorder) LookupSVC(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LookupSVC", reflect.TypeOf((*MockResolver)(nil).LookupSVC), arg0, arg1, arg2)
}
