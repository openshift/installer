#!/bin/bash
set -e
echo "Using resource group $1"
export AZREGION="centralus"
rm -r -f gw
mkdir gw
cp install-config.yaml gw
./openshift-install create manifests --dir=gw
mkdir -p gw/archive/manifests/original
cp gw/manifests/* gw/archive/manifests/original
rm -f gw/openshift/99_openshift-cluster-api_master-machines-*
rm -f gw/openshift/99_openshift-cluster-api_worker-machineset-*
python3 setup-manifests.py $1
./setup-host-network.sh
cp gw/manifests/* gw/archive/manifests/
./openshift-install create ignition-configs --dir=gw
mkdir -p ~/.kube
cp gw/auth/kubeconfig ~/.kube/config
echo "Delete old resource group"
az group delete --name $1 --yes
echo "Create new resource group"
az group create --name $1 --location $AZREGION
az identity create -g $1 -n ${1}_userid
echo "Copy RHCOS to resource group"
export VHD_URL=https://rhcos.blob.core.windows.net/imagebucket/
export VHD_NAME=rhcos-42.80.20191002.0.vhd
az storage account create --location $AZREGION --name sa${1} --kind Storage --resource-group $1  --sku Standard_LRS
az storage container create --name vhd --account-name sa${1}
export ACCOUNT_KEY=$(az storage account keys list --account-name sa${1} --resource-group $1 --query "[0].value" -o tsv)
az storage blob copy start --account-name "sa${1}" --account-key "$ACCOUNT_KEY" --destination-blob "rhcos.vhd" --destination-container vhd --source-uri ${VHD_URL}${VHD_NAME}
echo "Waiting on copy of vhd"
status="unknown"
while [ "$status" != "success" ]
    do
    status=$(az storage blob show --container-name vhd --name "rhcos.vhd" --account-name "sa${1}"  --account-key "$ACCOUNT_KEY" -o json --query properties.copy.status | sed -e 's/^"//' -e 's/"$//')
    done
echo "Copy of vhd complete"


echo "Configure template with ignition files"
az storage container create --name files --account-name sa${1} --public-access blob
ACCOUNT_KEY=$(az storage account keys list --account-name sa${1} --resource-group $1 --query "[0].value" -o tsv)
az storage blob upload --account-name sa$1 --account-key $ACCOUNT_KEY -c "files" -f "gw/bootstrap.ign" -n "bootstrap.ign"
BOOTSTRAPURL=$(az storage blob url --account-name sa$1 --account-key $ACCOUNT_KEY -c "files" -n "bootstrap.ign" -o tsv)
python3 setup-variables.py $BOOTSTRAPURL $1

az network public-ip create -g $1 -n $1 --allocation-method static
az network public-ip create -g $1 -n $1app --allocation-method static

