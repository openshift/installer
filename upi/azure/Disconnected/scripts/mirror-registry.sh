#!/bin/bash

echo "Mirroring Registry https://docs.openshift.com/container-platform/4.3/installing/install_config/installing-restricted-networks-preparations.html#cli-installing-cli_installing-restricted-networks-preparations"

yum -y install podman httpd-tools
mkdir -pv /opt/registry/{auth,certs,data}

lvextend -L+20G /dev/rootvg/optlv
xfs_growfs /dev/rootvg/optlv

cd /opt/registry/certs
openssl req -newkey rsa:2048 -nodes -sha256 -keyout domain.key -x509 -days 365 -out domain.crt -subj "/C=US/ST=Virginia/L=Lynchburg/O=CloudFit/OU=IT Department/CN=${HOSTNAME}"
htpasswd -bBc /opt/registry/auth/htpasswd registry registry
cd ..

podman run --name mirror-registry -p 5000:5000 \
     -v /opt/registry/data:/var/lib/registry:z \
     -v /opt/registry/auth:/auth:z \
     -e "REGISTRY_AUTH=htpasswd" \
     -e "REGISTRY_AUTH_HTPASSWD_REALM=Registry Realm" \
     -e REGISTRY_AUTH_HTPASSWD_PATH=/auth/htpasswd \
     -v /opt/registry/certs:/certs:z \
     -e REGISTRY_HTTP_TLS_CERTIFICATE=/certs/domain.crt \
     -e REGISTRY_HTTP_TLS_KEY=/certs/domain.key \
     -d docker.io/library/registry:2

firewall-cmd --add-port=5000/tcp --zone=internal --permanent
firewall-cmd --add-port=5000/tcp --zone=public --permanent 
firewall-cmd --reload

cp /opt/registry/certs/domain.crt /etc/pki/ca-trust/source/anchors/
update-ca-trust

curl -u registry:registry -k https://vmmirrorregistrytest001:5000/v2/_catalog 


cd ~
echo ${RH_PULLSECRET} > ~/.pull_secret.txt
cat ~/.pull_secret.txt | jq . > ~/pull-secret.json
export base64registry=`echo -n 'registry:registry' | base64 -w0`

jq ".auths += [{${HOSTNAME}:5000:{auth}}]"

export OCP_RELEASE=4.3.0-x86_64
export LOCAL_REGISTRY='vmmirrorregistrytest001:5000' 
export LOCAL_REPOSITORY='ocp4/openshift430' 
export PRODUCT_REPO='openshift-release-dev' 
export LOCAL_SECRET_JSON='/root/pull-secret.json' 
export RELEASE_NAME="ocp-release" 

oc adm -a ${LOCAL_SECRET_JSON} release mirror \
     --from=quay.io/${PRODUCT_REPO}/${RELEASE_NAME}:${OCP_RELEASE} \
     --to=${LOCAL_REGISTRY}/${LOCAL_REPOSITORY} \
     --to-release-image=${LOCAL_REGISTRY}/${LOCAL_REPOSITORY}:${OCP_RELEASE}