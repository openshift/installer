# Adding a node via the node-joiner tool

## Deprecated

For OpenShift 4.17 and later, the node-joiner.sh and node-joiner-monitor.sh scripts cannot be used. The node-joiner scripts were created for OpenShift 4.16 as a development preview. They have been deprecated and replaced by the "oc adm node-image" command. See [Adding worker nodes to an on-premise cluster](https://docs.openshift.com/container-platform/4.17/nodes/nodes/nodes-nodes-adding-node-iso.html) for more details.

## Pre-requisites
1. The `oc` tool must be available in the execution environment (the "user host").
2. Ensure `oc` is properly configured to manage the cluster to be expanded, i.e. by either:
   - running `oc login` command with the required cluster credentials, or
   - running `export KUBECONFIG=<cluster `kubeconfig` file>`

## Setup
1. Download the [node-joiner.sh](./node-joiner.sh) script in a working directory in
   the user host (the "assets folder").
2. Create a `nodes-config.yaml` in the assets folder. This configuration file must contain the 
   list of all the nodes that the user wants to add to the target cluster. At minimum, the name and primary interface MAC address must be specified. For example:
```
hosts:
    - hostname: extra-worker-0
      interfaces:
        - name: eth0
          macAddress: 00:02:46:e3:9e:7c
    - hostname: extra-worker-1
      interfaces:
        - name: eth0
          macAddress: 00:02:46:e3:9e:8c
    - hostname: extra-worker-2
      interfaces:
        - name: eth0
          macAddress: 00:02:46:e3:9e:9c
```
3. Optionally, it's possible to specify - for each node - an `NMState` configuration block denoted below as `networkConfig`
   (it will be applied during the first boot), for example:
```
hosts:
    - hostname: extra-worker-0
      interfaces:
        - name: eth0
          macAddress: 00:02:46:e3:9e:7c
      networkConfig:
        interfaces:
          - name: eth0
            type: ethernet
            state: up
            mac-address: 00:02:46:e3:9e:7c
            ipv4:
              enabled: true
              address:
                - ip: 192.168.111.90
                  prefix-length: 24
              dhcp: false
        dns-resolver:
          config:
            server:
              - 192.168.111.1
        routes:
          config:
            - destination: 0.0.0.0/0 
              next-hop-address: 192.168.111.1
              next-hop-interface: eth0
              table-id: 254
    - hostname: extra-worker-1
      interfaces:
        - name: eth0
          macAddress: 00:02:46:e3:9e:8c
    - hostname: extra-worker-2
      interfaces:
        - name: eth0
          macAddress: 00:02:46:e3:9e:9c
```

## ISO generation
Run [node-joiner.sh](./node-joiner.sh):
```bash
$ ./node-joiner.sh
```
The script will generate a temporary namespace prefixed with `openshift-node-joiner` in the target cluster,
where a pod will be launched to execute the effective node-joiner workload.
In case of success, the `node.x86_64.iso` ISO image will be downloaded in the assets folder.

### Configuration file name
By default the script looks for a configuration file named `nodes-config.yaml`. It's possible to specify a 
different config file name, as the first parameter of the script:

```bash
$ ./node-joiner.sh config.yaml
```

## Nodes joining
Use the iso image to boot all the nodes listed in the configuration file, and wait for the related
certificate signing requests (CSRs) to appear. When adding a new node to the cluster, two pending CSRs will
be generated, and they must be manually approved by the user.

Use the following command or [node-joiner-monitor.sh](./node-joiner-monitor.sh) described below to monitor the pending certificates:
```
$ oc get csr
```
Use the `oc` `approve` command to approve them:
```
$ oc adm certificate approve <csr_name>
```
Once all the pendings certificates will be approved, then the new node will become available:
```
$ oc get nodes
NAME            STATUS   ROLES                  AGE   VERSION
extra-worker-0  Ready    worker                  1h   v1.29.3+8628c3c                                        
master-0        Ready    control-plane,master   31h   v1.29.3+8628c3c
master-1        Ready    control-plane,master   32h   v1.29.3+8628c3c
master-2        Ready    control-plane,master   32h   v1.29.3+8628c3c
```

# Monitoring
After a node is booted using the ISO image, progress can be monitored using the node-joiner-monitor.sh script. 

Download the [node-joiner-monitor.sh](./node-joiner-monitor.sh) script to a local directory.

The script requires the IP address of the node to monitor.

Run [node-joiner-monitor.sh](./node-joiner-monitor.sh):
```bash
$ ./node-joiner-monitor.sh 192.168.111.90
```

The script will execute a command to monitor the node using a temporary namespace with
prefix `openshift-node-joiner-monitor` in the target cluster. The output of this command
is printed out to stdout.

The script shows useful information about the node as it joins the cluster.  
* Pre-flight validations. In case the node does not pass one or more validations, the installation will not start. The output of the failed validations are reported to allow users to fix the problem(s) when required.
* Installation progress indicating the current stage is shown. For example, writing of the image to disk, and initial reboot are reported.
* CSRs requiring the user's approval are shown.

The script exits either after the node has joined the cluster and is in ready state or after 90 minutes have elapsed.

Sample monitoring output:
```
INFO[2024-04-29T22:45:39-04:00] Monitoring IPs: [192.168.111.90]             
INFO[2024-04-29T22:48:17-04:00] Node 192.168.111.90: Assisted Service API is available 
INFO[2024-04-29T22:48:17-04:00] Node 192.168.111.90: Cluster is adding hosts 
INFO[2024-04-29T22:48:17-04:00] Node 192.168.111.90: Updated image information (Image type is "full-iso", SSH public key is set) 
INFO[2024-04-29T22:48:22-04:00] Node 192.168.111.90: Host ca241aa5-4f86-42bf-95a3-6b7ab7d4d66a: Successfully registered 
WARNING[2024-04-29T22:48:32-04:00] Node 192.168.111.90: Host couldn't synchronize with any NTP server 
WARNING[2024-04-29T22:48:32-04:00] Node 192.168.111.90: Host extraworker-0: updated status from discovering to insufficient (Host does not meet the minimum hardware requirements: Host couldn't synchronize with any NTP server) 
INFO[2024-04-29T22:49:28-04:00] Node 192.168.111.90: Host extraworker-0: updated status from known to installing (Installation is in progress) 
INFO[2024-04-29T22:50:28-04:00] Node 192.168.111.90: Host: extraworker-0, reached installation stage Writing image to disk: 5% 
INFO[2024-04-29T22:50:33-04:00] Node 192.168.111.90: Host: extraworker-0, reached installation stage Writing image to disk: 16% 
INFO[2024-04-29T22:50:38-04:00] Node 192.168.111.90: Host: extraworker-0, reached installation stage Writing image to disk: 28% 
INFO[2024-04-29T22:50:43-04:00] Node 192.168.111.90: Host: extraworker-0, reached installation stage Writing image to disk: 40% 
INFO[2024-04-29T22:50:48-04:00] Node 192.168.111.90: Host: extraworker-0, reached installation stage Writing image to disk: 51% 
INFO[2024-04-29T22:50:53-04:00] Node 192.168.111.90: Host: extraworker-0, reached installation stage Writing image to disk: 67% 
INFO[2024-04-29T22:50:58-04:00] Node 192.168.111.90: Host: extraworker-0, reached installation stage Writing image to disk: 77% 
INFO[2024-04-29T22:51:03-04:00] Node 192.168.111.90: Host: extraworker-0, reached installation stage Writing image to disk: 88% 
INFO[2024-04-29T22:51:08-04:00] Node 192.168.111.90: Host: extraworker-0, reached installation stage Writing image to disk: 93% 
INFO[2024-04-29T22:51:13-04:00] Node 192.168.111.90: Host: extraworker-0, reached installation stage Rebooting 
INFO[2024-04-29T22:56:35-04:00] Node 192.168.111.90: Kubelet is running      
INFO[2024-04-29T22:56:45-04:00] Node 192.168.111.90: First CSR Pending approval 
INFO[2024-04-29T22:56:45-04:00] CSR csr-257ms with signerName kubernetes.io/kube-apiserver-client-kubelet and username system:serviceaccount:openshift-machine-config-operator:node-bootstrapper is Pending and awaiting approval 
INFO[2024-04-29T22:58:50-04:00] Node 192.168.111.90: Second CSR Pending approval 
INFO[2024-04-29T22:58:50-04:00] CSR csr-tc8xt with signerName kubernetes.io/kubelet-serving and username system:node:extraworker-0 is Pending and awaiting approval 
INFO[2024-04-29T22:58:50-04:00] Node 192.168.111.90: Node joined cluster     
INFO[2024-04-29T23:00:00-04:00] Node 192.168.111.90: Node is Ready           
```

