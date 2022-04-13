# Azure Disk Encryption with Customer Managed Keys

The installer can use an existing disk encryption set. The disk encryption set
should be created in a separate resource group than the installer to avoid
losing access to your keys after the cluster is destroyed.

### Set basic variables
```
export RESOURCEGROUP="<Resource Group>"
export LOCATION="<Location>"
export KEYVAULT_NAME="<Keyvault Name>"
export KEYVAULT_KEY_NAME="<Keyvault Key Name>"
export DISK_ENCRYPTION_SET_NAME="<Disk Encryption Set Name>"
export CLUSTER_SP_ID="<Service Principal ID>"
```

### Register features

Encryption at host is in preview, and must be enabled to be used.

```
az feature register --namespace "Microsoft.Compute" --name "EncryptionAtHost"
az feature show --namespace Microsoft.Compute --name EncryptionAtHost
az provider register -n Microsoft.Compute
```

### Create resource group
```
az group create --name $RESOURCEGROUP --location $LOCATION
```

### Create key vault
```
az keyvault create -n $KEYVAULT_NAME -g $RESOURCEGROUP -l $LOCATION \
        --enable-purge-protection true --enable-soft-delete true

az keyvault key create --vault-name $KEYVAULT_NAME -n $KEYVAULT_KEY_NAME \
        --protection software 
```

### Create an instance of a DiskEncryptionSet
```
KEYVAULT_ID=$(az keyvault show --name $KEYVAULT_NAME --query "[id]" -o tsv)
    
KEYVAULT_KEY_URL=$(az keyvault key show --vault-name $KEYVAULT_NAME --name \
        $KEYVAULT_KEY_NAME --query "[key.kid]" -o tsv)
    
az disk-encryption-set create -n $DISK_ENCRYPTION_SET_NAME -l $LOCATION -g \
        $RESOURCEGROUP --source-vault $KEYVAULT_ID --key-url $KEYVAULT_KEY_URL
```

### Grant the DiskEncryptionSet resource access to the key vault
```
DES_IDENTITY=$(az disk-encryption-set show -n $DISK_ENCRYPTION_SET_NAME -g \
        $RESOURCEGROUP --query "[identity.principalId]" -o tsv)

az keyvault set-policy -n $KEYVAULT_NAME -g $RESOURCEGROUP --object-id \
        $DES_IDENTITY --key-permissions wrapkey unwrapkey get
```

### Grant service principal reader permissions to the DiskEncryptionSet
```
DES_RESOURCE_ID=$(az disk-encryption-set show -n $DISK_ENCRYPTION_SET_NAME -g \
        $RESOURCEGROUP --query "[id]" -o tsv)

az role assignment create --assignee $CLUSTER_SP_ID --role Owner --scope \
        $DES_RESOURCE_ID -o jsonc
```

### Configure the disk encryption set in your install-config.yaml
```yaml
apiVersion: v1
baseDomain: my.basedomain.com
compute:
- architecture: amd64
  hyperthreading: Enabled
  name: worker
  platform:
    azure:
      encryptionAtHost: true
      osDisk:
        diskEncryptionSet:
          resourceGroup: diskEncryptionResourceGroup
          name: diskEncryptionSetName
  replicas: 3
controlPlane:
  architecture: amd64
  hyperthreading: Enabled
  name: master
  platform:
    azure:
      encryptionAtHost: true
      osDisk:
        diskEncryptionSet:
          resourceGroup: diskEncryptionResourceGroup
          name: diskEncryptionSetName
  replicas: 3
...
```

### Run the installer
```
openshift-install create cluster --log-level debug
```


### Update the annotations on the default storage class
```
oc patch storageclass managed-premium -p '{"metadata": {"annotations":{"storageclass.kubernetes.io/is-default-class":"false"}}}'
```

### Grab the identity of resource group created by the installer
```
az identity list --resource-group "<Cluster Resource Group>"
```

### Grant resource group permissions to the DiskEncryptionSet

The cluster's resource group needs permissions to write to the disk encryption
set. You could, for example, make the resource group identity an Owner.

```
az role assignment create --role Owner --assignee "<Identity Name>"
```
