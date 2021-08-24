// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package mocks

import (
	context "context"

	solaredge "github.com/clambin/solaredge"
	mock "github.com/stretchr/testify/mock"

	time "time"
)

// API is an autogenerated mock type for the API type
type API struct {
	mock.Mock
}

// GetPower provides a mock function with given fields: ctx, id, from, to
func (_m *API) GetPower(ctx context.Context, id int, from time.Time, to time.Time) ([]solaredge.PowerMeasurement, error) {
	ret := _m.Called(ctx, id, from, to)

	var r0 []solaredge.PowerMeasurement
	if rf, ok := ret.Get(0).(func(context.Context, int, time.Time, time.Time) []solaredge.PowerMeasurement); ok {
		r0 = rf(ctx, id, from, to)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]solaredge.PowerMeasurement)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int, time.Time, time.Time) error); ok {
		r1 = rf(ctx, id, from, to)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetPowerOverview provides a mock function with given fields: ctx, id
func (_m *API) GetPowerOverview(ctx context.Context, id int) (float64, float64, float64, float64, float64, error) {
	ret := _m.Called(ctx, id)

	var r0 float64
	if rf, ok := ret.Get(0).(func(context.Context, int) float64); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Get(0).(float64)
	}

	var r1 float64
	if rf, ok := ret.Get(1).(func(context.Context, int) float64); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Get(1).(float64)
	}

	var r2 float64
	if rf, ok := ret.Get(2).(func(context.Context, int) float64); ok {
		r2 = rf(ctx, id)
	} else {
		r2 = ret.Get(2).(float64)
	}

	var r3 float64
	if rf, ok := ret.Get(3).(func(context.Context, int) float64); ok {
		r3 = rf(ctx, id)
	} else {
		r3 = ret.Get(3).(float64)
	}

	var r4 float64
	if rf, ok := ret.Get(4).(func(context.Context, int) float64); ok {
		r4 = rf(ctx, id)
	} else {
		r4 = ret.Get(4).(float64)
	}

	var r5 error
	if rf, ok := ret.Get(5).(func(context.Context, int) error); ok {
		r5 = rf(ctx, id)
	} else {
		r5 = ret.Error(5)
	}

	return r0, r1, r2, r3, r4, r5
}

// GetSiteIDs provides a mock function with given fields: ctx
func (_m *API) GetSiteIDs(ctx context.Context) ([]int, error) {
	ret := _m.Called(ctx)

	var r0 []int
	if rf, ok := ret.Get(0).(func(context.Context) []int); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]int)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
