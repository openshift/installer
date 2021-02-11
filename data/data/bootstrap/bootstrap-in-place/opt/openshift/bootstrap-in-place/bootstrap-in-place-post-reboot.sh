#!/usr/bin/env bash
set -euoE pipefail ## -E option will cause functions to inherit trap

export KUBECONFIG=/etc/kubernetes/bootstrap-secrets/kubeconfig

function wait_for_api {
  until oc get csr &> /dev/null
  do
      echo "Waiting for api ..."
      sleep 5
  done
}

function restart_kubelet {
  echo "Restarting kubelet"
  until [ "$(oc get pod -n openshift-kube-apiserver-operator --selector='app=kube-apiserver-operator'  -o jsonpath='{.items[0].status.conditions[?(@.type=="Ready")].status}' | grep -c "True")" -eq 1 ];
  do
    echo "Waiting for kube-apiserver-operator ready condition to be True"
    sleep 10
  done
  # daemon-reload is required because /etc/systemd/system/kubelet.service.d/20-nodenet.conf is added after kubelet started
  systemctl daemon-reload
  systemctl restart kubelet

  while grep  bootstrap-kube-apiserver /etc/kubernetes/manifests/kube-apiserver-pod.yaml;
  do
    echo "Waiting for kube-apiserver to apply the new static pod configuration"
    sleep 10
  done
  systemctl restart kubelet
}

function approve_csr {
  echo "Approving csrs ..."
  until [ "$(oc get nodes --selector='node-role.kubernetes.io/master' -o jsonpath='{.items[0].status.conditions[?(@.type=="Ready")].status}' | grep -c "True")" -eq 1 ];
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

function clean {
  if [ -d "/etc/kubernetes/bootstrap-secrets" ]; then
     rm -rf /etc/kubernetes/bootstrap-*
  fi

  rm -rf /usr/local/bin/installer-gather.sh
  rm -rf /usr/local/bin/installer-masters-gather.sh
  rm -rf /var/log/log-bundle-bootstrap.tar.gz

  systemctl disable bootstrap-in-place-post-reboot.service
}

wait_for_api
approve_csr
restart_kubelet
wait_for_cvo
clean
