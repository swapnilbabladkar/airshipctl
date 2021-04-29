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

package pull_test

import (
	"io/ioutil"

	"path"
	"strings"
	"testing"

	fixtures "github.com/go-git/go-git-fixtures/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"opendev.org/airship/airshipctl/pkg/config"
	"opendev.org/airship/airshipctl/pkg/document/pull"
	"opendev.org/airship/airshipctl/pkg/document/repo"
	"opendev.org/airship/airshipctl/pkg/util"
	"opendev.org/airship/airshipctl/testutil"
)

func mockConfigFactory(t *testing.T, testGitDir string, chkOutOpts *config.RepoCheckout, tmpDir string) config.Factory {
	return func() (*config.Config, error) {
		cfg := testutil.DummyConfig()
		currentManifest, err := cfg.CurrentContextManifest()
		require.NoError(t, err)
		currentManifest.Repositories = map[string]*config.Repository{
			currentManifest.PhaseRepositoryName: {
				URLString:       testGitDir,
				CheckoutOptions: chkOutOpts,
				Auth: &config.RepoAuth{
					Type: "http-basic",
				},
			},
		}

		currentManifest.TargetPath = tmpDir

		_, err = repo.NewRepository(
			".",
			currentManifest.Repositories[currentManifest.PhaseRepositoryName],
		)
		require.NoError(t, err)

		return cfg, nil
	}
}

func TestPull(t *testing.T) {
	require := require.New(t)
	assert := assert.New(t)

	tests := []struct {
		name         string
		url          string
		checkoutOpts *config.RepoCheckout
		error        error
	}{
		{
			name: "TestCloneRepositoriesValidOpts",
			checkoutOpts: &config.RepoCheckout{
				Branch:        "master",
				LocalBranch:   true,
				ForceCheckout: false,
			},
			error: nil,
		},
		{
			name:  "TestCloneRepositoriesMissingCheckoutOptions",
			error: nil,
		},
		{
			name: "TestCloneRepositoriesInvalidOpts",
			checkoutOpts: &config.RepoCheckout{
				Branch:        "master",
				Tag:           "someTag",
				ForceCheckout: false,
			},
			error: config.ErrMutuallyExclusiveCheckout{},
		},
	}

	testGitDir := fixtures.Basic().One().DotGit().Root()
	dirNameFromURL := util.GitDirNameFromURL(testGitDir)
	globalTmpDir, cleanup := testutil.TempDir(t, "airshipctlCloneTest-")
	defer cleanup(t)

	for _, tt := range tests {
		tmpDir := path.Join(globalTmpDir, tt.name)
		expectedErr := tt.error
		chkOutOpts := tt.checkoutOpts
		t.Run(tt.name, func(t *testing.T) {
			cfgFactory := mockConfigFactory(t, testGitDir, tt.checkoutOpts, tmpDir)
			cfg, err := cfgFactory()
			require.NoError(err)
			currentManifest, err := cfg.CurrentContextManifest()
			require.NoError(err)

			err = pull.Pull(cfgFactory, false)
			if expectedErr != nil {
				assert.NotNil(err)
				assert.Equal(expectedErr, err)
			} else {
				require.NoError(err)
				assert.FileExists(path.Join(currentManifest.TargetPath, dirNameFromURL, "go/example.go"))
				assert.FileExists(path.Join(currentManifest.TargetPath, dirNameFromURL, ".git/HEAD"))
				contents, err := ioutil.ReadFile(path.Join(currentManifest.TargetPath, dirNameFromURL, ".git/HEAD"))
				require.NoError(err)
				if chkOutOpts == nil {
					assert.Equal(
						"ref: refs/heads/master",
						strings.TrimRight(string(contents), "\t \n"),
					)
				} else {
					assert.Equal(
						"ref: refs/heads/"+chkOutOpts.Branch,
						strings.TrimRight(string(contents), "\t \n"),
					)
				}
			}
		})
	}
	testutil.CleanUpGitFixtures(t)
}
