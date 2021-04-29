#!/usr/bin/env bash

# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# Example Usage
# ./tools/deployment/provider_common/30_deploy_controlplane.sh

set -xe

export KUBECONFIG=${KUBECONFIG:-"$HOME/.airship/kubeconfig"}
export KUBECONFIG_TARGET_CONTEXT=${KUBECONFIG_TARGET_CONTEXT:-"target-cluster"}
export KUBECONFIG_EPHEMERAL_CONTEXT=${KUBECONFIG_EPHEMERAL_CONTEXT:-"ephemeral-cluster"}

echo "create control plane"
airshipctl phase run controlplane-ephemeral --debug --wait-timeout 1000s

airshipctl cluster get-kubeconfig > ~/.airship/kubeconfig-tmp

mv ~/.airship/kubeconfig-tmp "${KUBECONFIG}"

echo "apply cni as a part of initinfra-networking"
airshipctl phase run initinfra-networking-target --debug

echo "Check nodes status"
kubectl --kubeconfig "${KUBECONFIG}" --context "${KUBECONFIG_TARGET_CONTEXT}" wait --for=condition=Ready nodes --all --timeout 4000s
kubectl get nodes --kubeconfig "${KUBECONFIG}" --context "${KUBECONFIG_TARGET_CONTEXT}"

echo "Waiting for  pods to come up"
kubectl --kubeconfig "${KUBECONFIG}" --context "${KUBECONFIG_TARGET_CONTEXT}" wait --for=condition=ready pods --all --timeout=4000s -A
kubectl --kubeconfig "${KUBECONFIG}" --context "${KUBECONFIG_TARGET_CONTEXT}" get pods -A

echo "Check machine status"
kubectl get machines --kubeconfig ${KUBECONFIG} --context "${KUBECONFIG_EPHEMERAL_CONTEXT}"

echo "Get cluster state for target workload cluster "
kubectl --kubeconfig ${KUBECONFIG} --context "${KUBECONFIG_EPHEMERAL_CONTEXT}" get cluster
