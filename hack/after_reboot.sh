#!/bin/bash -x
export KUBECONFIG=/etc/kubernetes/bootstrap-secrets/kubeconfig
function copy_file {
  ETCD_DIR="/boot/master/etcd"
  echo "Copy $ETCD_DIR"
  if [ -d "$ETCD_DIR" ]; then
    cp -r $ETCD_DIR /var/lib
    rm -rf $ETCD_DIR
  fi
  KUBE_DIR="/boot/master/kubernetes/"
  echo "Copy $KUBE_DIR"
  if [ -d "$KUBE_DIR" ]; then
    cp -r $KUBE_DIR/* /etc/kubernetes
    rm -rf $KUBE_DIR
  fi
}
function wait_for_api {
  until oc get csr &> /dev/null
    do
        echo "Waiting for api ..."
        sleep 30
    done
}
function restart_kubelet {
  echo "Restarting kubelet"
  while cat /etc/kubernetes/manifests/kube-apiserver-pod.yaml  | grep bootstrap-kube-apiserver; do
    echo "Waiting for kube-apiserver to apply the new static pod configuration"
    sleep 10
  done
  systemctl daemon-reload
  systemctl restart kubelet
}
function approve_csr {
  echo "Approving csrs ..."
  needed_to_approve=false
  until [ $(oc get nodes | grep master | grep -v NotReady | grep Ready | wc -l) -eq 1 ]; do
      needed_to_approve=true
      echo "Approving csrs ..."
     oc get csr -o go-template='{{range .items}}{{if not .status}}{{.metadata.name}}{{"\n"}}{{end}}{{end}}' | xargs oc adm certificate approve &> /dev/null || true
     sleep 30
    done
  # Restart kubelet only if node was added
  if $needed_to_approve ; then
    sleep 60
    restart_kubelet
  fi
}
function wait_for_cvo {
  echo "Waiting for cvo"
  until [ "$(oc get clusterversion -o jsonpath='{.items[0].status.conditions[?(@.type=="Available")].status}')" == "True" ]; do
      echo "Still waiting for cvo ..."
     sleep 30
    done
}
function clean {
  if [ -d "/etc/kubernetes/bootstrap-secrets" ]; then
     rm -rf /etc/kubernetes/bootstrap-*
  fi
}
copy_file
wait_for_api
approve_csr
wait_for_cvo
clean