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

package ifc

import (
	"opendev.org/airship/airshipctl/pkg/api/v1alpha1"
	"opendev.org/airship/airshipctl/pkg/cluster/clustermap"
	"opendev.org/airship/airshipctl/pkg/document"
	"opendev.org/airship/airshipctl/pkg/inventory/ifc"
)

// Helper is a phase helper that provides phases with additional config related information
type Helper interface {
	TargetPath() string
	PhaseRepoDir() string
	DocEntryPointPrefix() string
	WorkDir() string
	Phase(phaseID ID) (*v1alpha1.Phase, error)
	Plan(planID ID) (*v1alpha1.PhasePlan, error)
	ListPhases(o ListPhaseOptions) ([]*v1alpha1.Phase, error)
	ListPlans() ([]*v1alpha1.PhasePlan, error)
	ClusterMapAPIobj() (*v1alpha1.ClusterMap, error)
	ClusterMap() (clustermap.ClusterMap, error)
	ExecutorDoc(phaseID ID) (document.Document, error)
	PhaseBundleRoot() string
	Inventory() ifc.Inventory
	PhaseEntryPointBasePath() string
	PhaseConfigBundle() document.Bundle
}
