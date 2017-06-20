# Troubleshooting worker nodes using SSH

Tectonic worker nodes are not assigned a public IP address, only the master node. To debug a worker node, SSH to it through a master (bastion host) or use a VPN connected to the internal network.

To do so, perform the following:

## Set up SSH agent forwarding

Once a passphrase of the local ssh key is added to `ssh-agent`, you will not be prompted for the credentials the next time connecting to nodes via SSH or SCP. The following instructions outline adding a passphrase to the `ssh-agent` on the system.

  1. At the terminal, enter:

  `$ eval ssh-agent`

  2. Run the following:

  `$ ssh-add`

  The `ssh-add` command prompts for a private key passphrase and adds it to the list maintained by `ssh-agent`.

  3. Enter your private key passphrase.

  4. Before logging out, run the following:

  `$ kill $SSH_AGENT_PID`

  To automatically run this command when logging out, place it in the `.logout` file if you are using csh or tcsh. Place the command in the `.bash_logout` file if you are using bash.

## Get the IP address of the worker nodes

Run the following command:

`$ kubectl get nodes -o wide`

A table of nodes and their IP addresses is displayed:

```bash
NAME                                       STATUS    AGE      EXTERNAL-IP
ip-192-0-2-18.us-west-2.compute.internal   Ready     3d       <none>
ip-192-0-2-10.us-west-2.compute.internal   Ready     3d       203.0.113.3
ip-192-0-2-12.us-west-2.compute.internal   Ready     3d       <none>
```

## Connect to a master node

SSH to a master node with its `EXTERNAL-IP`, providing the `-A` flag to forward the local `ssh-agent`. Add the `-i` option giving the location of the ssh key known to Tectonic:

```bash
$ ssh -A core@203.0.113.3 -i /path/to/tectonic/cluster/ssh/key
```

## Connect to a worker node

The worker node is accessible from the master because both machines are on the same private network, but the master is the only public entry point into the cluster. From a master, reach worker nodes on their internal cluster IP addresses. This address is encoded in the node's host name by convention. In this example, the worker node `ip-192-0-2-18.us-west-2.compute.internal` listed by `kubectl get nodes -o wide` has the internal IP 192.0.2.18.

Having connected to a master, ssh from there to the target worker node's internal IP:

```bash
# From master node
$ ssh core@192.0.2.18
```

View logs on the worker node by using `journalctl -xe` or similar tools. [Reading the system log][journalctl] has more information.

To examine the kubelet logs, execute:
```sh
# From worker node
$ journalctl -u kubelet
```

To examine the status and logs of potentially failed containers, execute:
```sh
$ docker ps -a | grep -v pause 
...
$ docker logs ...
```

[journalctl]: https://github.com/coreos/docs/blob/master/os/reading-the-system-log.md
