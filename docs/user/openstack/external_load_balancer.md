# Deploying OpenShift with a user managed load balancer

This document explains how to deploy OpenShift with a user managed load balancer, rather than with the self-hosted, OpenShift managed, load balancer based on HAproxy & Keepalived.

## Table of Contents
- [Deploying OpenShift with a user managed load balancer](#deploying-openshift-with-an-user-managed-load-balancer)
  - [Table of Contents](#table-of-contents)
  - [Common prerequisites](#common-prerequisites)
  - [Deploy the load balancer](#deploy-the-load-balancer)
  - [Deploy OpenShift](#deploy-openshift)
  - [Known limitations](#known-limitations)
  - [Notes](#notes)

## Common prerequisites

* When deploying OpenShift with a user managed load balancer, it's required to bring your own network(s) before
  the deployment. This can be a tenant network or a provider network.
* The load balancer(s) will have to be deployed before installing OpenShift.
  * It has be connected to the network(s) where OpenShift will be deployed.
  * If it's on a server managed by OpenStack, allowed address pairs have to be configured
    for the port that will serve API and Ingress traffic, otherwise the VIP traffic 
    will be rejected by the OpenStack SDN when port security is enabled.
  * The firewall has to allow the following traffic on the load balancer (if a server managed by OpenStack, create a security group):
    * 22/TCP - SSH (to allow Ansible to perform remote tasks from the host where it runs). This rule can be removed after the load balancer is deployed.
    * 6443/TCP (from within and outside the cluster) - OpenShift API
    * 80/TCP (from within and outside the cluster) - Ingress HTTP
    * 443/TCP (from within and outside the cluster) - Ingress HTTPS
    * 22623/TCP (from within the OCP network) - Machine Config Server

## Deploy the load balancer

Before you install OpenShift, you must provision at least one load balancer.
The load balancer will manage the VIPs for API, Ingress and Machine Config Server services.
If you deploy in production, at least two load balancers should be deployed per network fabric for high availability.

You can use your own solution that suits your needs, or you can use this [Ansible role](https://github.com/shiftstack/ansible-role-routed-lb)
that has been created for testing purpose and can be used as an example. This role has not been tested in production,
therefore we can't recommand to use it outside of testing environments.

## Deploy OpenShift

Now that your load balancer(s) are ready, you can deploy OpenShift.
Here is an example of an `install-config.yaml`:

```yaml
apiVersion: v1
baseDomain: mydomain.test 
compute:
- name: worker
  platform:
    openstack:
      type: m1.xlarge
  replicas: 3
controlPlane:
  name: master
  platform:
    openstack:
      type: m1.xlarge
  replicas: 3
metadata:
  name: mycluster
networking:
  clusterNetwork:
  - cidr: 10.128.0.0/14
    hostPrefix: 23
  machineNetwork:
  - cidr: 192.168.10.0/24
platform:
  openstack:
    cloud: mycloud
    controlPlanePort:
      fixedIPs:
      - subnet:
          id: 8586bf1a-cc3c-4d40-bdf6-c243decc603a
    apiVIPs:
    - 192.168.10.5
    ingressVIPs:
    - 192.168.10.7
    loadBalancer:
      type: UserManaged
```

There are some important things to note here:

* `loadBalancer` is a new stanza created in OCP 4.13. The default type is `OpenShiftManagedDefault` (which will deploy HAproxy and Keepalived in OCP, known as the OpenShift managed load balancer). Setting it to `UserManaged` will allow a user managed load balancer to replace the OpenShift managed one.
* `platform.openstack.controlPlanePort.fixedIPs.subnet.id` is the subnet ID where both the OpenShift cluster and the user managed load balancer are deployed.
* In OCP 4.13 the feature had to be enabled as Technology Preview. This can be done by adding featureSet: `TechPreviewNoUpgrade` into the install-config.yaml.


Deploy the cluster:
```bash
openshift-install create cluster
```

## Known limitations

These limitations will eventually be addressed in our roadmap:

* Deploying OpenShift with static IPs for the machines is only supported on Baremetal platform.
* Changing the IP address for any OpenShift control plane VIP (API + Ingress) is currently not supported: once the user managed LB and the OpenShift cluster is deployed, the VIPs can’t be changed.
* Migrating an OpenShift cluster from the OpenShift managed LB to a user managed LB is currently not supported.


## Notes

* In combination with `FailureDomains`, this feature allows customers to deploy their OpenShift cluster across multiple subnets.
* Using a user managed load balancer has been proven to reduce the load on Kube API and help large scale deployments to perform better. This is because there are less
  pods fetching the API every few seconds (haproxy and keepalived monitors).
