# Generic Platform

This document describes the generic compute, network, and storage needs of a Kubernetes self-hosted cluster. In particular the master, infra, and worker nodes.

**Key User Question:** What is the minimum amount of compute/network architecture that is required to enable self-driving that a user cannot change?

## DNS Setup

- API server will generally be `<cluster-name>-k8s.<subdomain>`
- Console, Dex, etc will generally be `<cluster-name>.<subdomain>`

## Nodes

### Master Nodes

Master nodes run most, if not all, control plane components including the API server, etcd, scheduler, controller manager, etc.

- **Count:** At least one dedicated master node for etcd and Kubernetes components, 3 or 5 nodes for HA 
- **Network:**
   - Ingress 
      - MUST allow tcp port 22 (ssh) from user network 
      - MUST allow port 4789 (UDP) from masters & workers for flannel
      - MUST allow 32000-32002 from all for: Tectonic ingress (if using node ports for ingress like on AWS, otherwise use host ports on workers) 
      - SHOULD allow port 9100 from masters & workers for: Prometheus Node Exporter metrics 
      - MAY have tcp/udp port 30000-32767 (node port range open)  
      - MAY have tcp port 2379 (etcd client API when not using external etcd) 
      - MAY allow port 10255 from all for: read-only kubelet status & LB health checks for k8s API 
    - Egress 
      - MUST have 443 to download gcr, quay, and docker hub images 
      - MAY have 2379 to external etcd cluster 
      - MAY allow 2379-2380 for self-hosted etcd pods if using experimental etcd operator
      - MAY allow 12379-12380 for temporary etcd pod if using experimental etcd operator

- **Access Control**
  - MUST have any necessary API access for k8s cloud plugin functionality (i.e. AWS node IAM Role) 

- **Cloud-init/Ignition:**
  - MUST run the kubelet with label node-role.kubernetes.io/master 
  - MAY run the kubelet with label with label node-role.kubernetes.io/node 


### Infra Nodes (Optional) or Load Balancers

Infra nodes are special nodes that must be deployed when a load balancer isn't available in an environment.  

- **Count:** 0 or more 
- **Network:**
  - **Ingress**
     - MUST have tcp port 443 (Ingress controller for Tectonic Console and Dex)
   - **Egress**
     - MUST have 443 to download gcr, quay, and docker hub images 
- **Cloud-init/Ignition:**
    - MUST run the kubelet 
    - MAY have a taint annotation to avoid user workloads from running here 

### Worker Nodes

Worked nodes run all of the user applications. The only component they must run on-boot is the kubelet.

- **Count:** User specified N 
- **Network:**
    - **Ingress**
        - MUST allow all ports open to master nodes (TODO: be more specific) 
        - MUST have 30000 to 32767 host port range access open 
        - MUST allow port 4789 (UDP) from masters & workers for: VXLAN (flannel) 
        - SHOULD allow port 10250 from masters for k8s features: port-forward, exec, proxy 
        - SHOULD allow port 9100 from masters & workers for: Prometheus Node Exporter metrics 
        - SHOULD allow port 4194 from masters for: Heapster connections to CAdvisor 
        - MAY allow 2379-2380 for self-hosted etcd pods if using experimental etcd operator

    - **Egress**
        - MUST have 443 to download gcr, quay, and docker hub images 

- **Cloud-init/Ignition**
    - MUST run the kubelet with label node-role.kubernetes.io/node=true
