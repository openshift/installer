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

Now, copy `examples/tectonic.libvirt.yaml` and customize it. You're ready to begin! The workflow is the same:

```
tectonic init --config=<path-to-config>
tectonic install --dir=<clustername>
```

The cluster should be up and running in about 10-20 minutes, depending on how quickly the container images are downloaded.


## Differences between libvirt and aws:

1. We use the Libvirt DNS server. So, if you want to resolve those names on your host, you'll need to configure NetworkManager's dns overlay mode (dnsmasq mode):
    1. Edit `/etc/NetworkManager/NetworkManager.conf` and set `dns=dnsmasq` in section `main`
    2. Tell dnsmasq to use your cluster. For me, this is: `echo server=/tt.testing/192.168.124.1 
     sudo tee /etc/NetworkManager/dnsmasq.d/tectonic.conf`
    3. restart NetworkManager
1. There isn't a load balancer. This means:
    1. We need to manually remap ports that the loadbalancer would
    2. Only the first server (e.g. master) is actually used. If you want to reach another, you have to manually update the domain name.
