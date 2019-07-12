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

Relevant files:
* **[files/etc/keepalived/keepalived.conf.tmpl](files/etc/keepalived/keepalived.conf.tmpl)** -
  `keepalived` configuration template
  * *NOTE:* The extension is `.tmpl` instead of `.template`, as the templating
    is done by `envsubst` in `keepalived.sh` and not via go templating in the
    installer.
* **[files/usr/local/bin/keepalived.sh](files/usr/local/bin/keepalived.sh)** -
  This script runs before `keepalived` starts and generates the `keepalived`
  configuration file from the template.
* **[systemd/units/keepalived.service](systemd/units/keepalived.service)** -
  systemd unit file for `keepalived`. This runs `keepalived.sh` to generate the
  proper configuration from the template and then runs podman to launch
  `keepalived`.
* **[files/usr/local/bin/fletcher8](files/usr/local/bin/fletcher8)** - Script
  that uses the
  [fletcher8](https://en.wikipedia.org/wiki/Fletcher%27s_checksum) algorithm to
  generate a hash from an input string. This is used by `keepalived.sh` to
  generate a hash based on the cluster name to generate VRRP ids for use with
  Keepalived and to ensure those IDs do not clash with another cluster on the
  same network.
* **[files/usr/local/bin/get_vip_subnet_cidr](files/usr/local/bin/get_vip_subnet_cidr)** -
  Script to determine the network CIDR for a given VIP.  This is used by
  `keepalived.sh` to determine which local interface is on the VIP’s network.

## Internal DNS

The bootstrap assets relating to DNS automate as much of the DNS requirements
internal to the cluster as possible.

There is a DNS VIP managed by `keepalived` in a manner similar to the API VIP
discussed above.

`coredns` runs with a custom `mdns` plugin (`coredns-mdns`) that can
dynamically generate the DNS SRV record for `etcd`, as well as resolve the
`etcd` hostnames.

Relevant files:
* **[files/etc/keepalived/keepalived.conf.tmpl](files/etc/keepalived/keepalived.conf.tmpl)** -
  `keepalived` configuration template, includes configuration for the DNS VIP
  * *NOTE:* The extension is `.tmpl` instead of `.template`, as the templating
    is done by `envsubst` in `keepalived.sh` and not via go templating in the
    installer.
* **[files/etc/dhcp/dhclient.conf](files/etc/dhcp/dhclient.conf)** - Sepcify
  that the bootstrap VM should use `localhost` as its primary DNS server.
* **[files/etc/coredns/Corefile](files/etc/coredns/Corefile)** - This is a
  template for a CoreDNS configuration file.
  * Includes a static entry for `api-int`
  * Enables the `mdns` plugin for dynamic `etcd` resolution
  * Forwards other queries along to the originally configured DNS server for
    that host
* **[files/usr/local/bin/coredns.sh](files/usr/local/bin/coredns.sh)** - A
  script to prepare the CoreDNS configuration
* **[systemd/units/coredns.service](systemd/units/coredns.service)** - systemd
  unit that launches CoreDNS via podman after first running `coredns.sh` to
  prepare the configuration
