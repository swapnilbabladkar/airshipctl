// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package redfishutils

import (
	"context"

	"github.com/stretchr/testify/mock"

	"opendev.org/airship/airshipctl/pkg/remote/power"
	"opendev.org/airship/airshipctl/pkg/remote/redfish"
)

// MockClient is a fake Redfish client for unit testing.
type MockClient struct {
	mock.Mock
	nodeID string
}

// NodeID provides a stubbed method that can be mocked to test functions that use the Redfish client without
// making any Redfish API calls or requiring the appropriate Redfish client settings.
//
//     Example usage:
//         client := redfishutils.NewClient()
//         client.On("NodeID").Return(<return values>)
//
//         err := client.NodeID()
func (m *MockClient) NodeID() string {
	args := m.Called()
	return args.String(0)
}

// EjectVirtualMedia provides a stubbed method that can be mocked to test functions that use the
// Redfish client without making any Redfish API calls or requiring the appropriate Redfish client
// settings.
//
//     Example usage:
//         client := redfishutils.NewClient()
//         client.On("EjectVirtualMedia").Return(<return values>)
//
//         err := client.EjectEphemeralVirtualMedia(<args>)
func (m *MockClient) EjectVirtualMedia(ctx context.Context) error {
	args := m.Called()
	return args.Error(0)
}

// RebootSystem provides a stubbed method that can be mocked to test functions that use the Redfish client without
// making any Redfish API calls or requiring the appropriate Redfish client settings.
//
//     Example usage:
//         client := redfishutils.NewClient()
//         client.On("RebootSystem").Return(<return values>)
//
//         err := client.RebootSystem(<args>)
func (m *MockClient) RebootSystem(ctx context.Context) error {
	args := m.Called()
	return args.Error(0)
}

// SetBootSourceByType provides a stubbed method that can be mocked to test functions that use the
// Redfish client without making any Redfish API calls or requiring the appropriate Redfish client settings.
//
//     Example usage:
//         client := redfishutils.NewClient()
//         client.On("SetBootSourceByType").Return(<return values>)
//
//         err := client.SetBootSourceByType(<args>)
func (m *MockClient) SetBootSourceByType(ctx context.Context) error {
	args := m.Called()
	return args.Error(0)
}

// SetVirtualMedia provides a stubbed method that can be mocked to test functions that use the
// Redfish client without making any Redfish API calls or requiring the appropriate Redfish client settings.
//
//     Example usage:
//         client := redfishutils.NewClient()
//         client.On("SetVirtualMedia").Return(<return values>)
//
//         err := client.SetVirtualMedia(<args>)
func (m *MockClient) SetVirtualMedia(ctx context.Context, isoPath string) error {
	args := m.Called()
	return args.Error(0)
}

// SystemPowerOff provides a stubbed method that can be mocked to test functions that use the
// Redfish client without making any Redfish API calls or requiring the appropriate Redfish client settings.
//
//     Example usage:
//         client := redfishutils.NewClient()
//         client.On("SystemPowerOff").Return(<return values>)
//
//         err := client.SystemPowerOff(<args>)
func (m *MockClient) SystemPowerOff(ctx context.Context) error {
	args := m.Called()
	return args.Error(0)
}

// SystemPowerOn provides a stubbed method that can be mocked to test functions that use the
// Redfish client without making any Redfish API calls or requiring the appropriate Redfish client settings.
//
//     Example usage:
//         client := redfishutils.NewClient()
//         client.On("SystemPowerOn").Return(<return values>)
//
//         err := client.SystemPowerOn(<args>)
func (m *MockClient) SystemPowerOn(ctx context.Context) error {
	args := m.Called()
	return args.Error(0)
}

// SystemPowerStatus provides a stubbed method that can be mocked to test functions that use the
// Redfish client without making any Redfish API calls or requiring the appropriate Redfish client settings.
//
//     Example usage:
//         client := redfishutils.NewClient()
//         client.On("SystemPowerStatus").Return(<return values>)
//
//         err := client.SystemPowerStatus(<args>)
func (m *MockClient) SystemPowerStatus(ctx context.Context) (power.Status, error) {
	args := m.Called()
	powerStatus, ok := args.Get(0).(power.Status)
	if !ok {
		return power.StatusUnknown, args.Error(2)
	}

	return powerStatus, args.Error(1)
}

// RemoteDirect mocks remote client interface
func (m *MockClient) RemoteDirect(ctx context.Context, isoURL string) error {
	if isoURL == "" {
		return redfish.ErrRedfishMissingConfig{What: "isoURL"}
	}
	return m.Called().Error(0)
}

// NewClient returns a mocked Redfish client in order to test functions that use the Redfish client without making any
// Redfish API calls.
func NewClient(redfishURL string, insecure bool, useProxy bool, username string,
	password string) (*MockClient, error) {
	if redfishURL == "" {
		return nil, redfish.ErrRedfishMissingConfig{What: "Redfish URL"}
	}

	// Retrieve system ID from end of Redfish URL
	systemID := redfish.GetResourceIDFromURL(redfishURL)
	if len(systemID) == 0 {
		return nil, redfish.ErrRedfishMissingConfig{What: "management URL system ID"}
	}

	m := &MockClient{nodeID: systemID}
	return m, nil
}
