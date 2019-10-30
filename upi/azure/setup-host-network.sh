cat > gw/manifests/ingress-controller-02-default.yaml <<EOF
apiVersion: operator.openshift.io/v1
kind: IngressController
metadata:
  finalizers:
  - ingresscontroller.operator.openshift.io/finalizer-ingresscontroller
  name: default
  namespace: openshift-ingress-operator
spec:
  endpointPublishingStrategy: 
    type: HostNetwork
  replicas: 3
EOF

