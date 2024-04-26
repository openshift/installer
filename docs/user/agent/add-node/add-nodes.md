# Adding a node via the node-joiner tool

## Pre-requisites
1. The `oc` tool must be available in the execution environment (the "user host").
2. The user host has a valid network connection to the target OpenShift cluster to be expanded.

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

## ISO generation
Run the [node-joiner.sh](./node-joiner.sh):
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
Use the following command to monitor the pending certificates:
```
$ oc get csr
```
User the `oc` `approve` command to approve them:
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