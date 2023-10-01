# Install on Data Disks


```bash
export OPENSHIFT_INSTALL_RELEASE_IMAGE_OVERRIDE="quay.io/openshift-release-dev/ocp-release:4.14.0-rc.6-x86_64"
export OPENSHIFT_INSTALL_EXPERIMENTAL_ETCD_DEDICATED=true

CLUSTER_NAME=az-a414notetcd2
INSTALL_DIR=${PWD}/installdir-azure-${CLUSTER_NAME}
mkdir $INSTALL_DIR
cat << EOF > $INSTALL_DIR/install-config.yaml
apiVersion: v1
metadata:
  name: $CLUSTER_NAME
publish: External
pullSecret: '$(cat ~/.openshift/pull-secret-latest.json)'
sshKey: |
  $(cat ~/.ssh/id_rsa.pub)
baseDomain: splat.azure.devcluster.openshift.com
platform:
  azure:
    baseDomainResourceGroupName: os4-common
    cloudName: AzurePublicCloud
    outboundType: Loadbalancer
    region: eastus
EOF

./openshift-install create manifests --dir $INSTALL_DIR

./openshift-install create cluster --dir $INSTALL_DIR

```