#!/bin/bash
set -e
echo "Start Deployment"
az group deployment create \
   --name $1 \
   --resource-group $1 \
   --template-uri "https://raw.githubusercontent.com/openshift/installer/master/upi/azure/azuredeploy.json" \
     --parameters "runit.parameters.json"
if [ $? -ne 0 ]
then 
   echo "Deployment Failed"
   exit $?
fi
./openshift-install --dir=gw wait-for bootstrap-complete --log-level debug
az vm stop --resource-group $1 --name bootstrap-0
az vm deallocate --resource-group $1 --name bootstrap-0 --no-wait
ACCOUNT_KEY=$(az storage account keys list --account-name sa${1} --resource-group $1 --query "[0].value" -o tsv)
az storage blob delete --account-key $ACCOUNT_KEY --account-name sa${1} --container-name files --name bootstrap.ign
./openshift-install --dir=gw wait-for install-complete --log-level debug
