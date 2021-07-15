# Setting FeatureGate during installation

To test new unsupported features developers can enable a custom feature set at the time of installation. To do so, they need to specify `OPENSHIFT_INSTALL_CUSTOM_FEATURESET_NAME` environment variable to define a desired [feature set](https://github.com/openshift/api/blob/master/config/v1/types_feature.go#L102) and begin the installation. After that the installer will generate a [FeatureGate](https://docs.openshift.com/container-platform/4.6/nodes/clusters/nodes-cluster-enabling-features.html) manifest that will be applied during bootstrap process.

Example:

```sh
export OPENSHIFT_INSTALL_CUSTOM_FEATURESET_NAME="IPv6DualStackNoUpgrade"
./openshift-install create cluster --dir ostest
```

It will generate the next manifest:

```yaml
apiVersion: config.openshift.io/v1
kind: FeatureGate
metadata:
  creationTimestamp: null
  name: cluster
spec:
  featureSet: IPv6DualStackNoUpgrade
status: {}
```

## Enabling or disabling custom Kubernetes feature gates

Additionally developers are able to define custom Kubernetes [feature gates](https://kubernetes.io/docs/reference/command-line-tools-reference/feature-gates/). To accomplish this they first have to set `OPENSHIFT_INSTALL_CUSTOM_FEATURESET_NAME` value to `CustomNoUpgrade`. Then they can define two additional environment variables: `OPENSHIFT_INSTALL_CUSTOM_FEATUREGATES_ENABLED` and `OPENSHIFT_INSTALL_CUSTOM_FEATUREGATES_DISABLED`. They both contain a comma-separated list of strings that correspond to enabled and disabled Kubernetes feature gates.

Example:

```sh
export OPENSHIFT_INSTALL_CUSTOM_FEATURESET_NAME="CustomNoUpgrade"
export OPENSHIFT_INSTALL_CUSTOM_FEATUREGATES_ENABLED="CSIMigrationAWS,CSIServiceAccountToken"
export OPENSHIFT_INSTALL_CUSTOM_FEATUREGATES_DISABLED="DefaultPodTopologySpread,DisableAcceleratorUsageMetrics"
./openshift-install create cluster --dir ostest
```

```yaml
apiVersion: config.openshift.io/v1
kind: FeatureGate
metadata:
  creationTimestamp: null
  name: cluster
spec:
  customNoUpgrade:
    disabled:
    - DefaultPodTopologySpread
    - DisableAcceleratorUsageMetrics
    enabled:
    - CSIMigrationAWS
    - CSIServiceAccountToken
  featureSet: CustomNoUpgrade
status: {}
```

**NOTE**: If `OPENSHIFT_INSTALL_CUSTOM_FEATURESET_NAME` value is not `CustomNoUpgrade` then `OPENSHIFT_INSTALL_CUSTOM_FEATUREGATES_ENABLED` and `OPENSHIFT_INSTALL_CUSTOM_FEATUREGATES_DISABLED` are ignored.
