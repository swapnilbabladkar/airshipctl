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

package cluster_test

import (
	"bytes"
	"fmt"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"opendev.org/airship/airshipctl/pkg/cluster"
	"opendev.org/airship/airshipctl/pkg/document"
	"opendev.org/airship/airshipctl/pkg/k8s/client/fake"
)

type mockStatusOptions struct{}

func (o mockStatusOptions) GetStatusMapDocs() (*cluster.StatusMap, []document.Document, error) {
	fakeClient := fake.NewClient(
		fake.WithCRDs(makeResourceCRD(annotationValidStatusCheck())),
		fake.WithDynamicObjects(makeResource("stable-resource", "stable")))
	fakeSM, err := cluster.NewStatusMap(fakeClient)
	if err != nil {
		return nil, nil, err
	}
	fakeDocBundle, err := document.NewBundleByPath("testdata/statusmap")
	if err != nil {
		return nil, nil, err
	}

	fakeDocs, err := fakeDocBundle.GetAllDocuments()
	if err != nil {
		return nil, nil, err
	}
	return fakeSM, fakeDocs, nil
}

func TestStatusRunner(t *testing.T) {
	statusOptions := mockStatusOptions{}
	b := bytes.NewBuffer(nil)
	err := cluster.StatusRunner(statusOptions, b)
	require.NoError(t, err)
	expectedOutput := fmt.Sprintf("Kind Name Status Resource stable-resource Stable ")
	space := regexp.MustCompile(`\s+`)
	str := space.ReplaceAllString(b.String(), " ")
	assert.Equal(t, expectedOutput, str)
}
