apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  name: kube-controller-manager
  namespace: kube-system
  labels:
    k8s-app: kube-controller-manager
spec:
  replicas: 2
  template:
    metadata:
      labels:
        k8s-app: kube-controller-manager
    spec:
      nodeSelector:
        master: "true"
      containers:
      - name: kube-controller-manager
        image: ${hyperkube_image}
        command:
        - ./hyperkube
        - controller-manager
        - --allocate-node-cidrs=true
        - --configure-cloud-routes=false
        - --cluster-cidr=${cluster_cidr}
        - --root-ca-file=/etc/kubernetes/secrets/ca.crt
        - --service-account-private-key-file=/etc/kubernetes/secrets/service-account.key
        - --leader-elect=true
        - --cloud-provider=${cloud_provider}
        volumeMounts:
        - name: secrets
          mountPath: /etc/kubernetes/secrets
          readOnly: true
        - name: ssl-host
          mountPath: /etc/ssl/certs
          readOnly: true
      volumes:
      - name: secrets
        secret:
          secretName: kube-controller-manager
      - name: ssl-host
        hostPath:
          path: /usr/share/ca-certificates
      dnsPolicy: Default # Don't use cluster DNS.
