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

package implementations_test

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	versionclient "k8s.io/apimachinery/pkg/util/version"
	clusterctlv1 "sigs.k8s.io/cluster-api/cmd/clusterctl/api/v1alpha3"
	"sigs.k8s.io/yaml"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"opendev.org/airship/airshipctl/pkg/clusterctl/implementations"
)

type Version struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec VersionSpec `json:"spec"`
}

type VersionSpec struct {
	Version string `json:"version"`
}

func TestNewRepository(t *testing.T) {
	tests := []struct {
		name           string
		root           string
		versions       map[string]string
		defaultVersion string
		Error          error
	}{
		{
			name: "simple repository success",
			root: "testdata",
			versions: map[string]string{
				"v0.0.1": "functions/capi",
			},
			Error:          nil,
			defaultVersion: "v0.0.1",
		},
		{
			name: "invalid version",
			root: "testdata",
			versions: map[string]string{
				"malformed-version": "functions/capi",
			},
			Error: implementations.ErrNoVersionsAvailable{Versions: map[string]string{}},
		},
		{
			name: "multiple repository versions",
			root: "testdata",
			versions: map[string]string{
				"v0.0.1": "functions/1",
				"v0.2.3": "functions/2",
				"v7.0.2": "functions/3",
				"v0.3.2": "functions/4",
				"v4.0.2": "functions/5",
				"v1.0.2": "functions/6",
			},
			Error:          nil,
			defaultVersion: "v7.0.2",
		},
		{
			name:     "Empty version",
			root:     "testdata",
			versions: map[string]string{},
			Error:    implementations.ErrNoVersionsAvailable{Versions: map[string]string{}},
		},
	}

	for _, tt := range tests {
		repo, err := implementations.NewRepository(tt.root, tt.versions)
		expectedErr := tt.Error
		defaultVersion := tt.defaultVersion
		t.Run(tt.name, func(t *testing.T) {
			if expectedErr != nil {
				assert.Equal(t, expectedErr, err)
				assert.Nil(t, repo)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, repo)
				assert.Equal(t, defaultVersion, repo.DefaultVersion())
			}
		})
	}
}

func TestGetFile(t *testing.T) {
	tests := []struct {
		name           string
		root           string
		versions       map[string]string
		expectErr      bool
		resultVersion  string
		versionToUse   string
		fileToUse      string
		resultContract string
	}{
		{
			name: "single version",
			root: "testdata",
			versions: map[string]string{
				"v0.0.1": "functions/1",
			},
			expectErr:     false,
			resultVersion: "v0.0.1",
			versionToUse:  "v0.0.1",
		},
		{
			name: "latest version",
			root: "testdata",
			versions: map[string]string{
				"v0.0.1": "functions/1",
				"v0.0.2": "functions/2",
				"v0.0.3": "functions/3",
			},
			expectErr:     false,
			resultVersion: "v0.0.3",
			versionToUse:  "latest",
		},
		{
			name: "failed to bundle",
			root: "testdata",
			versions: map[string]string{
				"v1.3.2": "does-not-exist",
			},
			versionToUse: "v1.3.2",
			expectErr:    true,
		},
		{
			name: "multiple repository versions",
			root: "testdata",
			versions: map[string]string{
				"v0.0.1": "functions/1",
				"v0.2.3": "functions/2",
				"v7.0.2": "functions/3",
			},
			expectErr:     false,
			resultVersion: "v0.0.2",
			versionToUse:  "v0.2.3",
		},
		{
			name: "version doesn't exist",
			root: "testdata",
			versions: map[string]string{
				"v1.3.2": "does-not-exist",
			},
			versionToUse: "v1.3.3",
			expectErr:    true,
		},
		{
			name: "test valid metadata",
			root: "testdata",
			versions: map[string]string{
				"v0.2.3": "functions/2",
			},
			expectErr:      false,
			versionToUse:   "v0.2.3",
			fileToUse:      "metadata.yaml",
			resultContract: "v1alpha2",
		},
		{
			name: "test valid metadata",
			root: "testdata",
			versions: map[string]string{
				"v0.2.0": "functions/1",
			},
			expectErr:      true,
			versionToUse:   "v0.2.3",
			fileToUse:      "metadata.yaml",
			resultContract: "v1alpha2",
		},
	}
	for _, tt := range tests {
		root := tt.root
		versions := tt.versions
		resultVersion := tt.resultVersion
		versionToUse := tt.versionToUse
		expectErr := tt.expectErr
		fileToUse := tt.fileToUse
		resultContract := tt.resultContract
		t.Run(tt.name, func(t *testing.T) {
			repo, err := implementations.NewRepository(root, versions)
			require.NoError(t, err)
			assert.NotNil(t, repo)
			b, err := repo.GetFile(versionToUse, fileToUse)
			if expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				if fileToUse == "metadata.yaml" {
					gotMetadata := metadata(t, b)
					parsedVersion, err := versionclient.ParseSemantic(versionToUse)
					require.NoError(t, err)
					assert.Equal(t, resultContract, gotMetadata.GetReleaseSeriesForVersion(parsedVersion).Contract)
				} else {
					gotVersion := version(t, b)
					assert.Equal(t, resultVersion, gotVersion.Spec.Version)
				}
			}
		})
	}
}

func version(t *testing.T, versionBytes []byte) *Version {
	t.Helper()
	ver := &Version{}

	err := yaml.Unmarshal(versionBytes, ver)
	require.NoError(t, err)
	return ver
}

func metadata(t *testing.T, metadataBytes []byte) *clusterctlv1.Metadata {
	t.Helper()
	m := &clusterctlv1.Metadata{}
	err := yaml.Unmarshal(metadataBytes, m)
	require.NoError(t, err)
	return m
}

func TestComponentsPath(t *testing.T) {
	versions := map[string]string{
		"v0.0.1": "functions/1",
	}
	repo, err := implementations.NewRepository("testdata", versions)
	require.NoError(t, err)
	assert.NotEmpty(t, repo.ComponentsPath())
}

func TestRootPath(t *testing.T) {
	versions := map[string]string{
		"v0.0.1": "functions/1",
	}
	repo, err := implementations.NewRepository("testdata", versions)
	require.NoError(t, err)
	assert.Equal(t, "testdata", repo.RootPath())
}

func TestGetVersions(t *testing.T) {
	tests := []struct {
		name             string
		root             string
		versions         map[string]string
		expectedVersions []string
	}{
		{
			name: "single version",
			root: "testdata",
			versions: map[string]string{
				"v0.0.1": "functions/1",
			},
			expectedVersions: []string{"v0.0.1"},
		},
		{
			name: "multiple repository versions",
			root: "testdata",
			versions: map[string]string{
				"v0.0.1":            "functions/1",
				"v0.2.3":            "functions/2",
				"v7.0.2":            "functions/3",
				"malformed-version": "doesn't matter",
			},
			expectedVersions: []string{"v0.0.1", "v0.2.3", "v7.0.2"},
		},
	}

	for _, tt := range tests {
		root := tt.root
		versions := tt.versions
		expectedVersions := tt.expectedVersions
		t.Run(tt.name, func(t *testing.T) {
			repo, err := implementations.NewRepository(root, versions)
			require.NoError(t, err)
			actualVersions, err := repo.GetVersions()
			assert.NoError(t, err)
			// this will make sure that slices are sorted in a same way
			sort.Strings(expectedVersions)
			sort.Strings(actualVersions)
			assert.Equal(t, expectedVersions, actualVersions)
		})
	}
}
