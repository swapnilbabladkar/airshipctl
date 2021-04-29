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

	"opendev.org/airship/airshipctl/cmd/cluster/checkexpiration"
	"opendev.org/airship/airshipctl/cmd/cluster/resetsatoken"
	"opendev.org/airship/airshipctl/pkg/config"
)

const (
	// TODO: (kkalynovskyi) Add more description when more subcommands are added
	clusterLong = `
This command provides capabilities for interacting with a Kubernetes cluster,
such as getting status and deploying initial infrastructure.
`
)

// NewClusterCommand creates a command for interacting with a Kubernetes cluster.
func NewClusterCommand(cfgFactory config.Factory) *cobra.Command {
	clusterRootCmd := &cobra.Command{
		Use:   "cluster",
		Short: "Manage Kubernetes clusters",
		Long:  clusterLong[1:],
	}

	clusterRootCmd.AddCommand(NewStatusCommand(cfgFactory))
	clusterRootCmd.AddCommand(resetsatoken.NewResetCommand(cfgFactory))
	clusterRootCmd.AddCommand(checkexpiration.NewCheckCommand(cfgFactory))
	clusterRootCmd.AddCommand(NewGetKubeconfigCommand(cfgFactory))
	clusterRootCmd.AddCommand(NewListCommand(cfgFactory))

	return clusterRootCmd
}
