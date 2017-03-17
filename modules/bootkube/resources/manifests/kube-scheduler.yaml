apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  name: kube-scheduler
  namespace: kube-system
  labels:
    k8s-app: kube-scheduler
spec:
  replicas: 2
  template:
    metadata:
      labels:
        k8s-app: kube-scheduler
    spec:
      nodeSelector:
        master: "true"
      containers:
      - name: kube-scheduler
        image: ${hyperkube_image}
        command:
        - ./hyperkube
        - scheduler
        - --leader-elect=true
