#!/usr/bin/env bash
set -euoE pipefail ## -E option will cause functions to inherit trap

echo "Running bootkube bootstrap-in-place post reboot"
export KUBECONFIG=/etc/kubernetes/bootstrap-secrets/kubeconfig

function wait_for_api {
  until oc get csr &> /dev/null
  do
    echo "Waiting for api ..."
    sleep 5
  done
}

# This is required since the progress service (https://github.com/openshift/installer/blob/dd9047c4c119e942331f702a4b7da85c60042da5/data/data/bootstrap/files/usr/local/bin/report-progress.sh#L22-L33),
# usually dedicated to creating the bootstrap ConfigMap, will fail to create this configmap in case of bootstrap-in-place single node deployment,
# due to the lack of a control plane when bootkube is complete
function signal_bootstrap_complete {
  until oc get cm bootstrap -n kube-system &> /dev/null
  do
    echo "Creating bootstrap configmap ..."
    oc create cm bootstrap -n kube-system --from-literal status=complete || true
    sleep 5
  done
}

function release_lease {
  local ns="$1"
  local lease="$2"

  until [ "$(oc get leases -n "${ns}" "${lease}" | grep -c "${lease}")" -eq 0 ];
  do
    echo "Deleting ${ns} ${lease} lease"
    oc delete leases -n "${ns}" "${lease}" || sleep 5
  done
  until [ "$(oc get cm -n "${ns}" "${lease}" | grep -c "${lease}")" -eq 0 ];
  do
    echo "Deleting ${ns} ${lease} cm"
    oc delete cm -n "${ns}" "${lease}" || sleep 5
  done
}

function release_cvo_lease {
  if [ ! -f /opt/openshift/release_cvo_lease.done ]
  then
    release_lease openshift-cluster-version version
    touch /opt/openshift/release_cvo_lease.done
  fi
}

function release_cpc_lease {
  if [ ! -f /opt/openshift/release_cpc_lease.done ]
  then
    release_lease kube-system cluster-policy-controller-lock
    touch /opt/openshift/release_cpc_lease.done
  fi
}

function restart_kubelet {
  echo "Waiting for kube-apiserver-operator"
  until [ "$(oc get pod -n openshift-kube-apiserver-operator --selector='app=kube-apiserver-operator' -o jsonpath='{.items[*].status.conditions[?(@.type=="Ready")].status}' | grep -c "True")" -eq 1 ];
  do
    echo "Waiting for kube-apiserver-operator ready condition to be True"
    oc get --raw='/readyz' &> /dev/null || echo "Api is not available"
    sleep 10
  done

  echo "Restarting kubelet"
  # daemon-reload is required because /etc/systemd/system/kubelet.service.d/20-nodenet.conf is added after kubelet started
  systemctl daemon-reload
  systemctl restart kubelet

  while grep bootstrap-kube-apiserver /etc/kubernetes/manifests/kube-apiserver-pod.yaml;
  do
    echo "Waiting for kube-apiserver to apply the new static pod configuration"
    sleep 10
  done

  echo "Restarting kubelet"
  systemctl restart kubelet
}

function approve_csr {
  echo "Waiting for node to report ready status"
  # use [*] and not [0] in the jsonpath because the node resource may not have been created yet
  until [ "$(oc get nodes --selector='node-role.kubernetes.io/master' -o jsonpath='{.items[*].status.conditions[?(@.type=="Ready")].status}' | grep -c "True")" -eq 1 ];
  do
    echo "Approving csrs ..."
    oc get csr -o go-template='{{range .items}}{{if not .status}}{{.metadata.name}}{{"\n"}}{{end}}{{end}}' | xargs --no-run-if-empty oc adm certificate approve &> /dev/null || true
    sleep 10
  done
  echo "node is ready"
}

function wait_for_cvo {
  echo "Waiting for cvo"
  until [ "$(oc get clusterversion -o jsonpath='{.items[*].status.conditions[?(@.type=="Available")].status}')" == "True" ];
  do
    echo "Still waiting for cvo ..."
    # print the not ready operators names and message
    oc get clusteroperator -o jsonpath='{range .items[*]}{@.metadata.name}: {range @.status.conditions[?(@.type=="Available")]}{@.type}={@.status} {@.message}{"\n"}{end}{end}' | grep -v "Available=True" || true
    sleep 20
  done
}

function restore_cvo_overrides {
    echo "Restoring CVO overrides"
    until \
        oc patch clusterversion.config.openshift.io version \
            --type=merge \
            --patch-file=/opt/openshift/original_cvo_overrides.patch
    do
        echo "Trying again to restore CVO overrides ..."
        sleep 10
    done
}

function clean {
  if [ -d "/etc/kubernetes/bootstrap-secrets" ]; then
     rm -rf /etc/kubernetes/bootstrap-*
  fi

  rm -rf /usr/local/bin/installer-gather.sh
  rm -rf /usr/local/bin/installer-masters-gather.sh
  rm -rf /var/log/log-bundle-bootstrap.tar.gz
  rm -rf /opt/openshift

  systemctl disable bootkube.service
}

wait_for_api
signal_bootstrap_complete
release_cvo_lease
release_cpc_lease
restore_cvo_overrides
approve_csr
restart_kubelet
wait_for_cvo
clean
