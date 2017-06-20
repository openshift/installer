# Troubleshooting etcd nodes using SSH

Tectonic etcd nodes are not assigned a public IP address, only the master node are. To debug an etcd node, SSH to it through a master (bastion host) or use a VPN connected to the internal network.

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

## Connect to a master node

SSH to a master node with its `EXTERNAL-IP`, providing the `-A` flag to forward the local `ssh-agent`. Add the `-i` option giving the location of the ssh key known to Tectonic:

```bash
$ ssh -A core@10.0.23.37 -i /path/to/tectonic/cluster/ssh/key
```

## Get the IP address of the etcd nodes

Run the following command on the master instance:

```sh
core@ip-10-0-23-37 ~ $ grep etcd /opt/tectonic/manifests/kube-apiserver.yaml 
        - --etcd-servers=http://10.0.23.31:2379
```

## Connect to an etcd node

```sh
# From the master node
$ ssh core@10.0.23.31
```

To investigate issues with etcd, execute:
```sh
$ systemctl status etcd-member && journalctl etcd-member
```
