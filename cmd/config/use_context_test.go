/*
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package config_test

import (
	"errors"
	"testing"

	cmd "opendev.org/airship/airshipctl/cmd/config"
	"opendev.org/airship/airshipctl/pkg/config"
	"opendev.org/airship/airshipctl/testutil"
)

func TestConfigUseContext(t *testing.T) {
	settings := func() (*config.Config, error) {
		return testutil.DummyConfig(), nil
	}
	cmdTests := []*testutil.CmdTest{
		{
			Name:    "config-use-context",
			CmdLine: "dummy_context",
			Cmd:     cmd.NewUseContextCommand(settings),
		},
		{
			Name:    "config-use-context-no-args",
			CmdLine: "",
			Cmd:     cmd.NewUseContextCommand(settings),
			Error:   errors.New("accepts 1 arg(s), received 0"),
		},
		{
			Name:    "config-use-context-does-not-exist",
			CmdLine: "foo",
			Cmd:     cmd.NewUseContextCommand(settings),
			Error:   errors.New("missing configuration: context with name 'foo'"),
		},
	}

	for _, tt := range cmdTests {
		testutil.RunTest(t, tt)
	}
}
