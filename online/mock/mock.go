package mock

import (
	"time"

	"github.com/src-d/terraform-provider-online-net/online"
	"github.com/stretchr/testify/mock"
)

// OnlineClientMock can be used to unit test against the online.net client using testify/mock
type OnlineClientMock struct {
	mock.Mock
}

// Server is a mock call
func (o *OnlineClientMock) Server(id int) (*online.Server, error) {
	args := o.Called(id)
	return args.Get(0).(*online.Server), args.Error(1)
}

// SetServer is a mock call
func (o *OnlineClientMock) SetServer(s *online.Server) error {
	args := o.Called(s)
	return args.Error(0)
}

// GetRescueImages is a mock call
func (o *OnlineClientMock) GetRescueImages(serverID int) ([]string, error) {
	args := o.Called(serverID)
	return args.Get(0).([]string), args.Error(1)
}

// ListRPNv2 is a mock call
func (o *OnlineClientMock) ListRPNv2() ([]*online.RPNv2, error) {
	args := o.Called()
	return args.Get(0).([]*online.RPNv2), args.Error(1)
}

// RPNv2 is a mock call
func (o *OnlineClientMock) RPNv2(id int) (*online.RPNv2, error) {
	args := o.Called(id)
	return args.Get(0).(*online.RPNv2), args.Error(1)
}

// RPNv2ByName is a mock call
func (o *OnlineClientMock) RPNv2ByName(name string) (*online.RPNv2, error) {
	args := o.Called(name)
	return args.Get(0).(*online.RPNv2), args.Error(1)
}

// SetRPNv2 is a mock call
func (o *OnlineClientMock) SetRPNv2(r *online.RPNv2, wait time.Duration) error {
	args := o.Called(r, wait)
	return args.Error(0)
}

// DeleteRPNv2 is a mock call
func (o *OnlineClientMock) DeleteRPNv2(id int, wait time.Duration) error {
	args := o.Called(id, wait)
	return args.Error(0)
}

// BootRescueMode is a mock call
func (o *OnlineClientMock) BootRescueMode(serverID int, image string) (*online.RescueCredentials, error) {
	args := o.Called(serverID, image)
	return args.Get(0).(*online.RescueCredentials), args.Error(1)
}

// BootNormalMode is a mock call
func (o *OnlineClientMock) BootNormalMode(serverID int) error {
	args := o.Called(serverID)
	return args.Error(0)
}
