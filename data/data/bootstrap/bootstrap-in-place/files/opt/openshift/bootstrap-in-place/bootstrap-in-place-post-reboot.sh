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

function restart_kubelet {
  echo "Restarting kubelet"
  until [ "$(oc get pod -n openshift-kube-apiserver-operator --selector='app=kube-apiserver-operator' -o jsonpath='{.items[0].status.conditions[?(@.type=="Ready")].status}' | grep -c "True")" -eq 1 ];
  do
    echo "Waiting for kube-apiserver-operator ready condition to be True"
    sleep 10
  done
  # daemon-reload is required because /etc/systemd/system/kubelet.service.d/20-nodenet.conf is added after kubelet started
  systemctl daemon-reload
  systemctl restart kubelet

  while grep bootstrap-kube-apiserver /etc/kubernetes/manifests/kube-apiserver-pod.yaml;
  do
    echo "Waiting for kube-apiserver to apply the new static pod configuration"
    sleep 10
  done
  systemctl restart kubelet
}

function approve_csr {
  echo "Approving csrs ..."
  # use [*] and not [0] in the jsonpath because the node resource may not have been created yet
  until [ "$(oc get nodes --selector='node-role.kubernetes.io/master' -o jsonpath='{.items[*].status.conditions[?(@.type=="Ready")].status}' | grep -c "True")" -eq 1 ];
  do
    echo "Approving csrs ..."
    oc get csr -o go-template='{{range .items}}{{if not .status}}{{.metadata.name}}{{"\n"}}{{end}}{{end}}' | xargs --no-run-if-empty oc adm certificate approve &> /dev/null || true
    sleep 30
  done
}

function wait_for_cvo {
  echo "Waiting for cvo"
  until [ "$(oc get clusterversion -o jsonpath='{.items[0].status.conditions[?(@.type=="Available")].status}')" == "True" ];
  do
    echo "Still waiting for cvo ..."
    sleep 30
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

  systemctl disable bootkube.service
}

wait_for_api
signal_bootstrap_complete
restore_cvo_overrides
approve_csr
restart_kubelet
wait_for_cvo
clean
