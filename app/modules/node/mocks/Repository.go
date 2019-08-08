// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import entities "gitlab.com/velo-labs/cen/app/entities"
import mock "github.com/stretchr/testify/mock"

// Repository is an autogenerated mock type for the Repository type
type Repository struct {
	mock.Mock
}

// SaveCredit provides a mock function with given fields: credit
func (_m *Repository) SaveCredit(credit entities.Credit) error {
	ret := _m.Called(credit)

	var r0 error
	if rf, ok := ret.Get(0).(func(entities.Credit) error); ok {
		r0 = rf(credit)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
