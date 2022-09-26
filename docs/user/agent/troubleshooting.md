# Troubleshooting

This document describes good places to look to troubleshoot problems with the agent installer.

## Services

After your agent ISO has booted, the following services are started in roughly the listed order to prepare and start your cluster installation. If your cluster does not start to install, check none of the services below have failed. If one of them has failed then checking its service log would be a good first step.

For example, to see the service log for create-cluster-and-infraenv.service:

````
journalctl -u create-cluster-and-infraenv.service
````

### Services that are started on all nodes
* selinux.service
  - This service applies a selinux policy to enable chronyd to run.
* set-hostname.service 
  - If a hostname is defined for a host in the agent-config.yaml configuration file, this service determines if there is a hostname to apply and applies it.
* pre-network-manage-config.service
  - If NMStateConfigs have been defined, this service applies the network configurations to the host.
* node-zero.service
  - The node-zero service determines if the current host is the rendezvous host (aka node0). The rendezvous host is where assisted-service runs. It coordinates host validation and coordinates the cluster installation when all hosts have been validated and requirements have been met. On the rendezvous host, the service will write the rendezvous host IP to /etc/assisted-service/node0.
* agent.service
  - The agent service registers this host with the assisted-service REST-API. It also receives commands from assisted-service to perform host validation and installation.

### Services that are started on only the rendezvous host (node0)
* assisted-service-pod.service
  - This pod is where the assisted-service and its database containers are deployed.
* assisted-service.service
  - This service working in tandem with agent.service coordinates host validation and checks requirements are met before the cluster can be installed. The agent installer uses the assisted-service REST-API to create and start the cluster installation.
* assisted-service-db.service
  - The database assisted-service uses to save information about the cluster and the installation.
* create-cluster-and-infraenv.service
  - This service registers the cluster and infraenv with the assisted-service REST-API.
* apply-host-config.service
  - If the agent-config.yaml contains role and root device hints, this service determines if there is a role and/or root device hint to apply for this host and applies them.
* start-cluster-installation.service
  - This service waits for the cluster to become ready to be installed. To be ready, all hosts that will form the cluster must register with assisted-service and their validations must pass. Once the cluster status reaches "ready", this service calls the assisted-service REST-API to start the cluster installation.

## View cluster events:

From the rendezvous host (node0), to list the events that have happened on the cluster being installed:

````
export CLUSTER_ID=$(curl "http://127.0.0.1:8090/api/assisted-install/v2/clusters/" | jq '.[0].id' | tr -d '"')
curl "http://127.0.0.1:8090/api/assisted-install/v2/events?cluster_id=$CLUSTER_ID" | jq '.'
````