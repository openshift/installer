# Installer Test Scripts

These tools automate using the installer and provide simple tools to work with the cluster once its launched.

## Basic Usage
These commands bring up a cluster, wait for it to be ready, setup and start bootkube, and check for cluster availability.

```bash
# set AWS credential environment variables
export AWS_ACCESS_KEY_ID=[actual key]
export AWS_SECRET_ACCESS_KEY=[actual key]

# set cluster environment variables for new environment
source ./default.env.sh

# use installer to provision cluster
./launch_cluster_aws.sh

# wait until cluster is available
./wait_for_dns.sh

# transfer configuration and start bootkube
./setup_bootkube.sh
sleep 5

# check for API server availability and Pod status
kubectl --kubeconfig=./output/${CLUSTER_NAME}/assets/auth/kubeconfig get pods --all-namespaces
kubectl --kubeconfig=./output/${CLUSTER_NAME}/assets/auth/kubeconfig get nodes
```

## Configuration
**aws_payload.tmpl.sh**
The JSON body of the request to the Tectonic installer

**default.env.sh**
The default values for new clusters, options are documented there.

### Console auth
Username: `admin@example.com`
Password: `tectonicTestPass11042016`

### Switching clusters
If you are working with multiple clusters and want to switch the one you are working on, each cluster profile has an 'env' file with the correct environment variables to use that profile. To use it set the variables in your shell by running:
```bash
source ./output/[Cluster Name]/env
```
Replace "Cluster Name" with the name of your cluster.
