# Control Plane Replicas

The number of control plane replicas are either set in install-config.yaml or agent-cluster-install.yaml if ZTP manifests are used.

install-config.yaml
```
controlPlane:
  name: master
  replicas: 3
```

agent-cluster-install.yaml
```
spec:
  provisionRequirements:
    controlPlaneAgents: 3
```

# 4.18 update

Starting with 4.18, replicas of 5, 4, 3, or 1 are supported.

# Prior to 4.18

Only control plane replicas of either 3 or 1 are supported.

