#!/bin/bash

while test $# -gt 0; do
  case "$1" in
    -h|--help)
      echo "Create UPI OCP 4.X on Azure MAG, see Readme"
      echo " "
      echo "options:"
      echo "-h, --help                show brief help"
      echo "-w, --worker-nodes=number specify number of worker nodes, 3 default"
      exit 0
      ;;
    -w)
      shift
      if test $# -gt 0; then
        export WORKERNODES=$1
      else
        export WORKERNODES=3
      fi
      shift
      ;;
    --worker-nodes*)
      export PROCESS=`echo $1 | sed -e 's/^[^=]*=//g'`
      shift
      ;;
    *)
      break
      ;;
  esac
done

export CLUSTER_NAME=`yq -r .metadata.name install-config.yaml`
export AZURE_REGION=`yq -r .platform.azure.region install-config.yaml`
export SSH_KEY=`yq -r .sshKey install-config.yaml | xargs`
export BASE_DOMAIN=`yq -r .baseDomain install-config.yaml`
export BASE_DOMAIN_RESOURCE_GROUP=`yq -r .platform.azure.baseDomainResourceGroupName install-config.yaml`

az cloud set -n AzureUSGovernment
az login

python3 -c '
import yaml;
path = "install-config.yaml";
data = yaml.full_load(open(path));
data["compute"][0]["replicas"] = 0;
open(path, "w").write(yaml.dump(data, default_flow_style=False))'

openshift-install create manifests

read -p "press follow step 6 in azuremag quickstart guide, press [enter] to continue..." dummy

rm -fv openshift/99_openshift-cluster-api_master-machines-*.yaml
rm -fv openshift/99_openshift-cluster-api_worker-machineset-*.yaml


python3 -c '
import yaml;
path = "manifests/cluster-scheduler-02-config.yml";
data = yaml.full_load(open(path));
data["spec"]["mastersSchedulable"] = False;
open(path, "w").write(yaml.dump(data, default_flow_style=False))'

python3 -c '
import yaml;
path = "manifests/cluster-dns-02-config.yml";
data = yaml.full_load(open(path));
del data["spec"]["publicZone"];
del data["spec"]["privateZone"];
open(path, "w").write(yaml.dump(data, default_flow_style=False))'

export INFRA_ID=`yq -r '.status.infrastructureName' manifests/cluster-infrastructure-02-config.yml`
export RESOURCE_GROUP=`yq -r '.status.platformStatus.azure.resourceGroupName' manifests/cluster-infrastructure-02-config.yml`

openshift-install create ignition-configs

tree

az group create --name $RESOURCE_GROUP --location $AZURE_REGION
az identity create -g $RESOURCE_GROUP -n ${INFRA_ID}-identity
az storage account create -g $RESOURCE_GROUP --location $AZURE_REGION --name ${CLUSTER_NAME}sa --kind Storage --sku Standard_LRS
export ACCOUNT_KEY=`az storage account keys list -g $RESOURCE_GROUP --account-name ${CLUSTER_NAME}sa --query "[0].value" -o tsv`

export VHD_URL=`curl -s https://raw.githubusercontent.com/openshift/installer/release-4.3/data/data/rhcos.json | jq -r .azure.url`

az storage container create --name vhd --account-name ${CLUSTER_NAME}sa
az storage container create --name files --account-name ${CLUSTER_NAME}sa --public-access blob
az storage blob upload --account-name ${CLUSTER_NAME}sa --account-key $ACCOUNT_KEY -c "files" -f "bootstrap.ign" -n "bootstrap.ign"
az storage blob copy start --account-name ${CLUSTER_NAME}sa --account-key $ACCOUNT_KEY --destination-blob "rhcos.vhd" --destination-container vhd --source-uri "$VHD_URL"

export PRINCIPAL_ID=`az identity show -g $RESOURCE_GROUP -n ${INFRA_ID}-identity --query principalId --out tsv`
export RESOURCE_GROUP_ID=`az group show -g $RESOURCE_GROUP --query id --out tsv`
az role assignment create --assignee "$PRINCIPAL_ID" --role 'Contributor' --scope "$RESOURCE_GROUP_ID"

az network dns zone create -g $RESOURCE_GROUP -n ${CLUSTER_NAME}.${BASE_DOMAIN}

nsServer0=`az network dns zone show --resource-group $RESOURCE_GROUP --name ${CLUSTER_NAME}.${BASE_DOMAIN} --query "nameServers[0]" -o tsv`
nsServer1=`az network dns zone show --resource-group $RESOURCE_GROUP --name ${CLUSTER_NAME}.${BASE_DOMAIN} --query "nameServers[1]" -o tsv`

az network dns record-set ns create --name ${CLUSTER_NAME} --resource-group $BASE_DOMAIN_RESOURCE_GROUP --zone-name ${BASE_DOMAIN}
az network dns record-set ns add-record --nsdname $nsServer0 -n ${CLUSTER_NAME} -g $BASE_DOMAIN_RESOURCE_GROUP -z ${BASE_DOMAIN}
az network dns record-set ns add-record --nsdname $nsServer1 -n ${CLUSTER_NAME} -g $BASE_DOMAIN_RESOURCE_GROUP -z ${BASE_DOMAIN}

az network private-dns zone create -g $RESOURCE_GROUP -n ${CLUSTER_NAME}.${BASE_DOMAIN}

status="unknown"
while [ "$status" != "success" ]
do
  status=`az storage blob show --container-name vhd --name "rhcos.vhd" --account-name ${CLUSTER_NAME}sa --account-key $ACCOUNT_KEY -o tsv --query properties.copy.status`
  echo $status
done

az deployment group create -g $RESOURCE_GROUP \
  --template-file "01_vnet.json" \
  --parameters baseName="$INFRA_ID"

az network private-dns link vnet create -g $RESOURCE_GROUP -z ${CLUSTER_NAME}.${BASE_DOMAIN} -n ${INFRA_ID}-network-link -v "${INFRA_ID}-vnet" -e false

export VHD_BLOB_URL=`az storage blob url --account-name ${CLUSTER_NAME}sa --account-key $ACCOUNT_KEY -c vhd -n "rhcos.vhd" -o tsv`

az deployment group create -g $RESOURCE_GROUP \
  --template-file "02_storage.json" \
  --parameters vhdBlobURL="$VHD_BLOB_URL" \
  --parameters baseName="$INFRA_ID"

az deployment group create -g $RESOURCE_GROUP \
  --template-file "03_infra.json" \
  --parameters privateDNSZoneName="${CLUSTER_NAME}.${BASE_DOMAIN}" \
  --parameters baseName="$INFRA_ID"

export PUBLIC_IP=`az network public-ip list -g $RESOURCE_GROUP --query "[?name=='${INFRA_ID}-master-pip'] | [0].ipAddress" -o tsv`
az network dns record-set a add-record -g $RESOURCE_GROUP -z ${CLUSTER_NAME}.${BASE_DOMAIN} -n api -a $PUBLIC_IP --ttl 60

export BOOTSTRAP_URL=`az storage blob url --account-name ${CLUSTER_NAME}sa --account-key $ACCOUNT_KEY -c "files" -n "bootstrap.ign" -o tsv`
export BOOTSTRAP_IGNITION=`jq -rcnM --arg v "2.2.0" --arg url $BOOTSTRAP_URL '{ignition:{version:$v,config:{replace:{source:$url}}}}' | base64 -w0`

az deployment group create -g $RESOURCE_GROUP \
  --template-file "04_bootstrap.json" \
  --parameters bootstrapIgnition="$BOOTSTRAP_IGNITION" \
  --parameters sshKeyData="$SSH_KEY" \
  --parameters baseName="$INFRA_ID"

export MASTER_IGNITION=`cat master.ign | base64`

az deployment group create -g $RESOURCE_GROUP \
  --template-file "05_masters.json" \
  --parameters masterIgnition="$MASTER_IGNITION" \
  --parameters sshKeyData="$SSH_KEY" \
  --parameters privateDNSZoneName="${CLUSTER_NAME}.${BASE_DOMAIN}" \
  --parameters baseName="$INFRA_ID"


openshift-install wait-for bootstrap-complete --log-level debug

az network nsg rule delete -g $RESOURCE_GROUP --nsg-name ${INFRA_ID}-controlplane-nsg --name bootstrap_ssh_in
az vm stop -g $RESOURCE_GROUP --name ${INFRA_ID}-bootstrap
az vm deallocate -g $RESOURCE_GROUP --name ${INFRA_ID}-bootstrap
az vm delete -g $RESOURCE_GROUP --name ${INFRA_ID}-bootstrap --yes
az disk delete -g $RESOURCE_GROUP --name ${INFRA_ID}-bootstrap_OSDisk --no-wait --yes
az network nic delete -g $RESOURCE_GROUP --name ${INFRA_ID}-bootstrap-nic --no-wait
az storage blob delete --account-key $ACCOUNT_KEY --account-name ${CLUSTER_NAME}sa --container-name files --name bootstrap.ign
az network public-ip delete -g $RESOURCE_GROUP --name ${INFRA_ID}-bootstrap-ssh-pip

export KUBECONFIG="$PWD/auth/kubeconfig"
oc get nodes
oc get clusteroperator

export WORKER_IGNITION=`cat worker.ign | base64`
az deployment group create -g $RESOURCE_GROUP \
  --template-file "06_workers.json" \
  --parameters workerIgnition="$WORKER_IGNITION" \
  --parameters sshKeyData="$SSH_KEY" \
  --parameters baseName="$INFRA_ID" \
  --parameters numberOfNodes="$WORKERNODES"