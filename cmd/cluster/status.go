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
	"opendev.org/airship/airshipctl/pkg/errors"
)

// NewStatusCommand creates a command which reports the statuses of a cluster's deployed components.
func NewStatusCommand(cfgFactory config.Factory) *cobra.Command {
	var kubeconfig string
	statusCmd := &cobra.Command{
		Use:   "status",
		Short: "Retrieve statuses of deployed cluster components",
		RunE:  clusterStatusRunE,
	}

	statusCmd.Flags().StringVar(
		&kubeconfig,
		"kubeconfig",
		"",
		"Path to kubeconfig associated with cluster being managed")

	return statusCmd
}

func clusterStatusRunE(cmd *cobra.Command, args []string) error {
	return errors.ErrNotImplemented{What: "Cluster Status"}
}
