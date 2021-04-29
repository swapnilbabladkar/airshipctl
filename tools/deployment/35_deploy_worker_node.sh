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

set -e

#Default wait timeout is 3600 seconds
export TIMEOUT=${TIMEOUT:-3600}
export KUBECONFIG=${KUBECONFIG:-"$HOME/.airship/kubeconfig"}
export KUBECONFIG_TARGET_CONTEXT=${KUBECONFIG_TARGET_CONTEXT:-"target-cluster"}
WORKER_NODE=${WORKER_NODE:-"node03"}
EPHEMERAL_DOMAIN_NAME="air-ephemeral"

# all vms. This can be removed once sushy tool is fixed
if type "virsh" > /dev/null; then
  for vm in $(sudo virsh list --all --name --state-running |grep ${EPHEMERAL_DOMAIN_NAME})
  do
    echo "Stop ephemeral node '$vm'"
    sudo virsh destroy $vm
  done
fi

node_timeout () {
  end=$(($(date +%s) + $TIMEOUT))
  for worker in $WORKER_NODE
  do
    while true; do
      if (kubectl --request-timeout 20s --kubeconfig $KUBECONFIG --context $KUBECONFIG_TARGET_CONTEXT get $1 $worker | grep -qw $2) ; then
        if [ "$1" = "node" ]; then
          kubectl --kubeconfig $KUBECONFIG --context $KUBECONFIG_TARGET_CONTEXT label nodes $worker node-role.kubernetes.io/worker=""
        fi

        echo -e "\nGet $1 status"
        kubectl --kubeconfig $KUBECONFIG --context $KUBECONFIG_TARGET_CONTEXT get $1
        break
      else
        now=$(date +%s)
        if [ $now -gt $end ]; then
          echo -e "\n$1 is not ready before TIMEOUT."
          exit 1
        fi
        echo -n .
        sleep 15
      fi
    done
  done
}

echo "Deploy worker node"
airshipctl phase run workers-target --debug

echo "Waiting $TIMEOUT seconds for bmh to be in ready state."
node_timeout bmh ready

echo "Classify and provision worker node"
airshipctl phase run  workers-classification --debug

echo "Waiting $TIMEOUT seconds for node to be provisioned."
node_timeout node Ready
