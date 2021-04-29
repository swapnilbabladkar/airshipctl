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

package document

import (
	"github.com/spf13/cobra"

	"opendev.org/airship/airshipctl/pkg/config"
	"opendev.org/airship/airshipctl/pkg/document/pull"
)

// NewPullCommand creates a new command for pulling airship document repositories
func NewPullCommand(cfgFactory config.Factory) *cobra.Command {
	var noCheckout bool
	documentPullCmd := &cobra.Command{
		Use:   "pull",
		Short: "Pulls documents from remote git repository",
		RunE: func(cmd *cobra.Command, args []string) error {
			return pull.Pull(cfgFactory, noCheckout)
		},
	}

	documentPullCmd.Flags().BoolVarP(&noCheckout, "no-checkout", "n", false,
		"No checkout is performed after the clone is complete.")

	return documentPullCmd
}
