apiVersion: v1
kind: Config
clusters:
- name: local
  cluster:
    server: ${server}
    certificate-authority-data: ${ca_cert}
users:
- name: kubelet
  user:
    client-certificate-data: ${kubelet_cert}
    client-key-data: ${kubelet_key}
contexts:
- context:
    cluster: local
    user: kubelet
