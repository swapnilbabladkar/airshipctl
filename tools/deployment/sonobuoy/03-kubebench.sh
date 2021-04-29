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

set -xe
: ${KUBECONFIG:="$HOME/.airship/kubeconfig"}
: ${KUBEBENCH_MASTER_PLUGIN:="https://raw.githubusercontent.com/vmware-tanzu/sonobuoy-plugins/master/cis-benchmarks/kube-bench-master-plugin.yaml"}
: ${KUBEBENCH_WORKER_PLUGIN:="https://raw.githubusercontent.com/vmware-tanzu/sonobuoy-plugins/master/cis-benchmarks/kube-bench-plugin.yaml"}
: ${TARGET_CLUSTER_CONTEXT:="target-cluster"}
# This shouldnot include minor version
: ${KUBEBENCH_K8S_VERSION:=1.18}
: ${TIMEOUT:=300}

mkdir -p /tmp/sonobuoy_snapshots/kubebench
cd /tmp/sonobuoy_snapshots/kubebench

# Run aggregator, and default plugins e2e and systemd-logs
sonobuoy run \
--kubeconfig ${KUBECONFIG} \
--context ${TARGET_CLUSTER_CONTEXT} \
--plugin ${KUBEBENCH_MASTER_PLUGIN} \
--plugin ${KUBEBENCH_WORKER_PLUGIN} \
--plugin-env kube-bench-master.KUBERNETES_VERSION=${KUBEBENCH_K8S_VERSION} \
--plugin-env kube-bench-master.KUBERNETES_VERSION=${KUBEBENCH_K8S_VERSION} \
--wait --timeout ${TIMEOUT} \
--log_dir /tmp/sonobuoy_snapshots/kubebench

# Get information on pods
kubectl get all -n sonobuoy --kubeconfig ${KUBECONFIG} --context ${TARGET_CLUSTER_CONTEXT}

# Check sonobuoy status
sonobuoy status --kubeconfig ${KUBECONFIG} --context ${TARGET_CLUSTER_CONTEXT}

# Get logs
sonobuoy logs

# Store Results
results=$(sonobuoy retrieve --kubeconfig ${KUBECONFIG} --context ${TARGET_CLUSTER_CONTEXT})
echo "Results: ${results}"

# Display Results
sonobuoy results $results
ls -ltr /tmp/sonobuoy_snapshots/kubebench