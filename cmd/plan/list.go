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

package plan

import (
	"github.com/spf13/cobra"

	"opendev.org/airship/airshipctl/pkg/config"
	"opendev.org/airship/airshipctl/pkg/phase"
)

const (
	listLong = `
List life-cycle plans which were defined in document model.
`
)

// NewListCommand creates a command which prints available phase plans
func NewListCommand(cfgFactory config.Factory) *cobra.Command {
	planCmd := &phase.PlanListCommand{Factory: cfgFactory}

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List plans",
		Long:  listLong[1:],
		RunE: func(cmd *cobra.Command, args []string) error {
			planCmd.Writer = cmd.OutOrStdout()
			return planCmd.RunE()
		},
	}
	return listCmd
}
