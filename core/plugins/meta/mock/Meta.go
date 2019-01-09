// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import meta "github.com/secure2work/nori/core/plugins/meta"
import mock "github.com/stretchr/testify/mock"

// Meta is an autogenerated mock type for the Meta type
type Meta struct {
	mock.Mock
}

// GetAuthor provides a mock function with given fields:
func (_m *Meta) GetAuthor() meta.Author {
	ret := _m.Called()

	var r0 meta.Author
	if rf, ok := ret.Get(0).(func() meta.Author); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(meta.Author)
	}

	return r0
}

// GetCore provides a mock function with given fields:
func (_m *Meta) GetCore() meta.Core {
	ret := _m.Called()

	var r0 meta.Core
	if rf, ok := ret.Get(0).(func() meta.Core); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(meta.Core)
	}

	return r0
}

// GetDependencies provides a mock function with given fields:
func (_m *Meta) GetDependencies() []meta.Dependency {
	ret := _m.Called()

	var r0 []meta.Dependency
	if rf, ok := ret.Get(0).(func() []meta.Dependency); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]meta.Dependency)
		}
	}

	return r0
}

// GetDescription provides a mock function with given fields:
func (_m *Meta) GetDescription() meta.Description {
	ret := _m.Called()

	var r0 meta.Description
	if rf, ok := ret.Get(0).(func() meta.Description); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(meta.Description)
	}

	return r0
}

// GetInterface provides a mock function with given fields:
func (_m *Meta) GetInterface() meta.Interface {
	ret := _m.Called()

	var r0 meta.Interface
	if rf, ok := ret.Get(0).(func() meta.Interface); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(meta.Interface)
	}

	return r0
}

// GetLicense provides a mock function with given fields:
func (_m *Meta) GetLicense() meta.License {
	ret := _m.Called()

	var r0 meta.License
	if rf, ok := ret.Get(0).(func() meta.License); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(meta.License)
	}

	return r0
}

// GetLinks provides a mock function with given fields:
func (_m *Meta) GetLinks() []meta.Link {
	ret := _m.Called()

	var r0 []meta.Link
	if rf, ok := ret.Get(0).(func() []meta.Link); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]meta.Link)
		}
	}

	return r0
}

// GetTags provides a mock function with given fields:
func (_m *Meta) GetTags() []string {
	ret := _m.Called()

	var r0 []string
	if rf, ok := ret.Get(0).(func() []string); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	return r0
}

// Id provides a mock function with given fields:
func (_m *Meta) Id() meta.ID {
	ret := _m.Called()

	var r0 meta.ID
	if rf, ok := ret.Get(0).(func() meta.ID); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(meta.ID)
	}

	return r0
}
