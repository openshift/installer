# Installer assets bundle

All installs of Tectonic produce an `assets.zip` bundle used to bootstrap the cluster. Though this archive is used to create the initial components, it contains several useful assets for administrators concerned with cluster lifecycle.

## Contents of assets.zip

* A TLS client certificate, and accompanying kubeconfig, with root-level access to the cluster to bypass Tectonic Identity.
* TLS assets used by Kubelets and control plane components.
* Manifests for all Tectonic components.
* Provider specific assets such as raw cloudconfig files.

## Using assets.zip to troubleshoot Tectonic

The main use of these files is to troubleshoot failing components.

### Troubleshooting authentication issues in Tectonic Identity

If [Tectonic Identity][tectonic-identity] has been misconfigured, the kubeconfig and manifests in `assets.zip` can be used to regain access and reset Identity to its default configuration.

```
unzip assets.zip
export KUBECONFIG=$PWD/assets/auth/kubeconfig
kubectl apply -f ./assets/tectonic/identity-config.yaml
```

Trigger a rolling update with `kubectl`:

```
kubectl patch deployment tectonic-identity \
    --patch "{\"spec\":{\"template\":{\"metadata\":{\"annotations\":{\"date\":\"`date +'%s'`\"}}}}}" \
    --namespace tectonic-system
```

## Sensitive information in assets.zip

Administrators should take proper precautions to ensure these assets remain secure.

[tectonic-identity]: user-management.md
