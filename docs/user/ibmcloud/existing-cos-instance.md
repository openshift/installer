# Using an Existing IBM Cloud Object Storage Instance

The OpenShift installer can use an existing IBM Cloud Object Storage (COS) instance instead of creating a new one during installation. This is useful when you want to:

- Use a pre-configured COS instance with specific settings
- Comply with organizational policies that require using existing resources
- Reduce installation time by reusing existing infrastructure

## Prerequisites

1. An existing IBM Cloud Object Storage instance
2. The CRN (Cloud Resource Name) of the COS instance
3. Appropriate IAM permissions to access the COS instance

## Finding Your COS Instance CRN

You can find your COS instance CRN using the IBM Cloud CLI:

```bash
ibmcloud resource service-instances --service-name cloud-object-storage --output json | jq -r '.[] | select(.name=="<your-cos-instance-name>") | .id'
```

Or through the IBM Cloud Console:
1. Navigate to your COS instance
2. Click on "Service credentials" or "Overview"
3. Copy the CRN value (format: `crn:v1:bluemix:public:cloud-object-storage:global:a/<account_id>:<instance_id>::`)

## Configuration

Add the `cosInstanceCRN` field to your `install-config.yaml` under the `platform.powerVS` section:

**Important**: The CRN value must be enclosed in quotes because it contains colons.

```yaml
apiVersion: v1
baseDomain: example.com
metadata:
  name: my-cluster
platform:
  powerVS:
    region: us-south
    cosInstanceCRN: "crn:v1:bluemix:public:cloud-object-storage:global:a/1234567890abcdef:abcd-1234-efgh-5678::"
    # ... other PowerVS platform settings
```

## What Happens During Installation

When you specify an existing COS instance:

1. **Validation**: The installer validates that the CRN is properly formatted and that the COS instance exists
2. **RHCOS Image Storage**: The installer uploads the RHCOS image to a bucket in your existing COS instance
3. **Bootstrap Ignition**: The installer stores bootstrap ignition data in a bucket in your existing COS instance
4. **No Instance Creation**: The installer skips creating a new COS instance

## What Happens During Uninstallation

When you destroy a cluster that used an existing COS instance:

1. **Bucket Cleanup**: The installer removes the buckets it created (RHCOS image bucket and bootstrap ignition bucket)
2. **Instance Preservation**: The installer **does not** delete the COS instance itself, preserving your existing resource
3. **Other Buckets**: Any other buckets in the COS instance remain untouched

## Important Notes

- The COS instance must be in the same IBM Cloud account as your cluster
- The installer will create and manage buckets within the COS instance
- Ensure your IAM credentials have sufficient permissions to create buckets and upload objects to the COS instance
- The COS instance CRN is stored in the cluster metadata for proper cleanup during uninstallation

## Example Complete Install Config

```yaml
apiVersion: v1
baseDomain: example.com
metadata:
  name: my-cluster
compute:
- architecture: amd64
  hyperthreading: Enabled
  name: worker
  platform: {}
  replicas: 3
controlPlane:
  architecture: amd64
  hyperthreading: Enabled
  name: master
  platform: {}
  replicas: 3
networking:
  clusterNetwork:
  - cidr: 10.128.0.0/14
    hostPrefix: 23
  machineNetwork:
  - cidr: 10.0.0.0/16
  networkType: OVNKubernetes
  serviceNetwork:
  - 172.30.0.0/16
platform:
  powerVS:
    region: us-south
    resourceGroupName: my-resource-group
    cosInstanceCRN: "crn:v1:bluemix:public:cloud-object-storage:global:a/1234567890abcdef:abcd-1234-efgh-5678::"
    vpcName: my-vpc
publish: External
pullSecret: '<your-pull-secret>'
sshKey: '<your-ssh-public-key>'
```

## Troubleshooting

### CRN Validation Fails

If you receive an error about an invalid CRN:
- Verify the CRN format matches: `crn:v1:bluemix:public:cloud-object-storage:global:a/<account_id>:<instance_id>::`
- Ensure the COS instance exists and is accessible with your credentials
- Check that you have the correct IAM permissions

### Permission Errors

If you encounter permission errors during installation:
- Verify your IAM API key has `Writer` or `Manager` role for the COS instance
- Ensure the service-to-service authorization policy allows VPC to access COS (for image import)

### Bucket Already Exists

If a bucket with the same name already exists:
- The installer will attempt to use the existing bucket
- Ensure the bucket is in the correct region
- Check that the bucket doesn't contain conflicting data