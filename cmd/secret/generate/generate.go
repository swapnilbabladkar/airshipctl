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

package generate

import (
	"github.com/spf13/cobra"

	"opendev.org/airship/airshipctl/cmd/secret/generate/encryptionkey"
)

// NewGenerateCommand creates a new command for generating secret information
func NewGenerateCommand() *cobra.Command {
	generateRootCmd := &cobra.Command{
		Use: "generate",
		// TODO(howell): Make this more expressive
		Short: "Generate various secrets",
	}

	generateRootCmd.AddCommand(encryptionkey.NewGenerateEncryptionKeyCommand())

	return generateRootCmd
}
