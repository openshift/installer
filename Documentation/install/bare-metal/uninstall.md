## Uninstall Tectonic

To uninstall Tectonic from a bare-metal cluster, set cluster nodes to boot from your previous, standard PXE image.

## Reinstall bare-metal Tectonic cluster

Tectonic bare metal clusters should be reprovisioned using the latest [Tectonic bare-metal][install-bare-metal] documentation. If you've set up your infrastructure as described, skip to step 4, *Tectonic Installer*. Use the Tectonic Installer application, then set your servers to PXE on the next boot and reboot:

```
ipmitool -H node.example.com -U USER -P PASS chassis bootdev pxe
```

[install-bare-metal]: index.md
