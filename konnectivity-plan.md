High level goal:
kube-apiserver (KAS) running in the bootstrap environment cannot connect to a
webhook whose endpoint is in the pod network. I would like to make an
architectural change which permits the bootstrap KAS to be able to access
webhooks hosted in the cluster. Non-bootstrap KAS instances are not affected.


Architecture change:
Provide connectivity to the bootstrap KAS using Konnectivity.
* Deploy a Konnectivity server in the bootstrap environment
* Deploy the bootstrap KAS with an EgresSelectorConfiguration to proxy cluster
  traffic to the local Konnectivity server.
* Deploy a DaemonSet which deploys a Konnectivity agent on all Nodes which
  connects back to the Konnectivity server in the bootstrap environment.
* During teardown of the bootstrap environment, also remove the Konnectivity
  agent DaemonSet.


To validate:
I believe this will permit the bootstrap KAS to access a webhook hosted in the
running cluster.

I believe this Konnectivity deployment will not impact the non-bootstrap KAS
instances, which will continue to route normally to cluster-hosted webhooks.


Investigation tasks:
Find the relevant code in installer which creates the bootstrap environment.

Determine what is different about the bootstrap KAS's networking environment
which means it cannot route to the pod network.

Validate my assumptions above (I believe...)

Determine if Konnectivity is already available in the OpenShift payload. I
believe HyperShift uses Konnectivity already. The code is available at
/home/mbooth/src/openshift/hypershift.

Determine how the Konnectivity agent would authenticate to the Konnectivity
server.
