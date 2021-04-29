/*
 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

     https://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

package baremetal_test

import (
	"testing"

	"opendev.org/airship/airshipctl/cmd/baremetal"
	"opendev.org/airship/airshipctl/pkg/inventory"
	"opendev.org/airship/airshipctl/testutil"
)

func TestBaremetal(t *testing.T) {
	tests := []*testutil.CmdTest{
		{
			Name:    "baremetal-with-help",
			CmdLine: "-h",
			Cmd:     baremetal.NewBaremetalCommand(nil),
		},
		{
			Name:    "baremetal-ejectmedia-with-help",
			CmdLine: "-h",
			Cmd:     baremetal.NewEjectMediaCommand(nil, &inventory.CommandOptions{}),
		},
		{
			Name:    "baremetal-poweroff-with-help",
			CmdLine: "-h",
			Cmd:     baremetal.NewPowerOffCommand(nil, &inventory.CommandOptions{}),
		},
		{
			Name:    "baremetal-poweron-with-help",
			CmdLine: "-h",
			Cmd:     baremetal.NewPowerOnCommand(nil, &inventory.CommandOptions{}),
		},
		{
			Name:    "baremetal-powerstatus-with-help",
			CmdLine: "-h",
			Cmd:     baremetal.NewPowerStatusCommand(nil, &inventory.CommandOptions{}),
		},
		{
			Name:    "baremetal-reboot-with-help",
			CmdLine: "-h",
			Cmd:     baremetal.NewRebootCommand(nil, &inventory.CommandOptions{}),
		},
		{
			Name:    "baremetal-remotedirect-with-help",
			CmdLine: "-h",
			Cmd:     baremetal.NewRemoteDirectCommand(nil, &inventory.CommandOptions{}),
		},
	}

	for _, tt := range tests {
		testutil.RunTest(t, tt)
	}
}
