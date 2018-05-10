# Libvirt howto

Tectonic has limited support for installing a Libvirt cluster. This is useful especially
for operator development.

## HOW TO:

### Setup and preparation
* Make sure you have the `virsh` binary on your path.
* Install the libvirt terraform provider:
    ```
    go get github.com/dmacvicar/terraform-provider-libvirt
    mkdir -p ~/.terraform.d/plugins
    cp $GOPATH/bin/terraform-provider-libvirt ~/.terraform.d/plugins/
    ```
* Decide on an IP range. In this example, `192.168.124.0/24`
* Decide on a domain. In this example, `tt.testing`
* Download the latest CoreOS image. This is not done automatically to avoid unnecessary downloads. e.g.
```
wget https://beta.release.core-os.net/amd64-usr/current/coreos_production_qemu_image.img.bz2
bunzip2 coreos_production_qemu_image.img.bz2
```

Now, copy `examples/tectonic.libvirt.yaml` and customize it. You're ready to begin! The workflow is the same, but only the `install assets` and `install bootstrap` steps are supported.


## Differences between libvirt and aws:

1. We use the Libvirt DNS server. So, if you want to resolve those names on your host, you'll need to configure NetworkManager's dns overlay mode (dnsmasq mode)
1. There isn't a load balancer. We need to manually remap port 6443 to 443
1. We may not support changing the number of workers.

## Remaining tasks
1. Provision the masters and update the DNS names
1. Provision the workers and update the ingress names
