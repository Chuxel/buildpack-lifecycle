// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/buildpacks/lifecycle/internal/layer (interfaces: MetadataRestorer)

// Package testmock is a generated GoMock package.
package testmock

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"

	buildpack "github.com/buildpacks/lifecycle/buildpack"
	layer "github.com/buildpacks/lifecycle/internal/layer"
	platform "github.com/buildpacks/lifecycle/platform"
)

// MockMetadataRestorer is a mock of MetadataRestorer interface.
type MockMetadataRestorer struct {
	ctrl     *gomock.Controller
	recorder *MockMetadataRestorerMockRecorder
}

// MockMetadataRestorerMockRecorder is the mock recorder for MockMetadataRestorer.
type MockMetadataRestorerMockRecorder struct {
	mock *MockMetadataRestorer
}

// NewMockMetadataRestorer creates a new mock instance.
func NewMockMetadataRestorer(ctrl *gomock.Controller) *MockMetadataRestorer {
	mock := &MockMetadataRestorer{ctrl: ctrl}
	mock.recorder = &MockMetadataRestorerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMetadataRestorer) EXPECT() *MockMetadataRestorerMockRecorder {
	return m.recorder
}

// Restore mocks base method.
func (m *MockMetadataRestorer) Restore(arg0 []buildpack.GroupBuildable, arg1 platform.LayersMetadata, arg2 platform.CacheMetadata, arg3 layer.SHAStore) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Restore", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(error)
	return ret0
}

// Restore indicates an expected call of Restore.
func (mr *MockMetadataRestorerMockRecorder) Restore(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Restore", reflect.TypeOf((*MockMetadataRestorer)(nil).Restore), arg0, arg1, arg2, arg3)
}
