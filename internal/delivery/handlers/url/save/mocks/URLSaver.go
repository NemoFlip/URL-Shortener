// Code generated by mockery v2.53.3. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// URLSaver is an autogenerated mock type for the URLSaver type
type URLSaver struct {
	mock.Mock
}

// SaveURL provides a mock function with given fields: urlToSave, alias
func (_m *URLSaver) SaveURL(urlToSave string, alias string) error {
	ret := _m.Called(urlToSave, alias)

	if len(ret) == 0 {
		panic("no return value specified for SaveURL")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string) error); ok {
		r0 = rf(urlToSave, alias)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewURLSaver creates a new instance of URLSaver. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewURLSaver(t interface {
	mock.TestingT
	Cleanup(func())
}) *URLSaver {
	mock := &URLSaver{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
