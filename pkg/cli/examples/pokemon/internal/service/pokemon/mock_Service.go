// Code generated by mockery v2.52.2. DO NOT EDIT.

package pokemon

import mock "github.com/stretchr/testify/mock"

// MockService is an autogenerated mock type for the Service type
type MockService struct {
	mock.Mock
}

// GetPokemons provides a mock function with no fields
func (_m *MockService) GetPokemons() (Pokemons, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetPokemons")
	}

	var r0 Pokemons
	var r1 error
	if rf, ok := ret.Get(0).(func() (Pokemons, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() Pokemons); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(Pokemons)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SearchPokemon provides a mock function with given fields: name
func (_m *MockService) SearchPokemon(name string) (*Pokemon, error) {
	ret := _m.Called(name)

	if len(ret) == 0 {
		panic("no return value specified for SearchPokemon")
	}

	var r0 *Pokemon
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (*Pokemon, error)); ok {
		return rf(name)
	}
	if rf, ok := ret.Get(0).(func(string) *Pokemon); ok {
		r0 = rf(name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*Pokemon)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewMockService creates a new instance of MockService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockService(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockService {
	mock := &MockService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
