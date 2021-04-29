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

package cluster

import (
	"github.com/spf13/cobra"

	"opendev.org/airship/airshipctl/pkg/config"
	"opendev.org/airship/airshipctl/pkg/phase"
)

const (
	listShort   = "Retrieve the list of defined clusters"
	listExample = `
# Retrieve cluster list
airshipctl cluster list --airshipconf /tmp/airconfig
airshipctl cluster list -o table
airshipctl cluster list -o name
`
)

// NewListCommand creates a command which retrieves list of clusters
func NewListCommand(cfgFactory config.Factory) *cobra.Command {
	o := &phase.ClusterListCommand{Factory: cfgFactory}
	cmd := &cobra.Command{
		Use:     "list",
		Short:   listShort,
		Example: listExample[1:],
		RunE:    listRunE(o),
	}
	flags := cmd.Flags()
	flags.StringVarP(&o.Format,
		"output", "o", "name", "'table' "+
			"and 'name' are available "+
			"output formats")

	return cmd
}

// listRunE returns a function to cobra command to be executed in runtime
func listRunE(o *phase.ClusterListCommand) func(
	cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		o.Writer = cmd.OutOrStdout()
		return o.RunE()
	}
}
