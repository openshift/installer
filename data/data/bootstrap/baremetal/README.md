# Bare Metal IPI Bootstrap Assets

The `baremetal` platform (IPI for Bare Metal hosts) includes some additional
assets on the bootstrap node for automating some infrastructure requirements
that would have normally been handled by some cloud infrastructure service.
The [Bare Metal IPI Networking Infrastructure design
document](../../../../docs/design/baremetal/networking-infrastructure.md)
covers the high-level background, and this document explains these
bootstrap assets in more detail.

## API failover from bootstrap to control plane machines

`keepalived` is used to manage the failover of a VIP (Virtual IP) for the API
server. This VIP first resides on the bootstrap VM. Once the master nodes come
up, the VIP will move to the control plane machines.

## Internal DNS

The bootstrap assets relating to DNS automate as much of the DNS requirements
internal to the cluster as possible.

There is a DNS VIP managed by `keepalived` in a manner similar to the API VIP
discussed above.

`coredns` runs with a custom `mdns` plugin (`coredns-mdns`).

Relevant files:
* **[files/etc/dhcp/dhclient.conf](files/etc/dhcp/dhclient.conf)** - Sepcify
  that the bootstrap VM should use `localhost` as its primary DNS server.
