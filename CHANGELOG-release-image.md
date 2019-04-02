# 4.0.0-0.9

Created: 2019-04-02 07:49:07 +0000 UTC

Image Digest: `sha256:29a057c89cc5cf9c7853c52b815ba95bd66d1a66cb274d1991ca4951f5920cb2`

Promoted from registry.svc.ci.openshift.org/ocp/release:4.0.0-0.nightly-2019-04-02-081046


## Changes from 4.0.0-0.8

### Components

* Kubernetes 1.12.4


### New images

* [baremetal-machine-controllers](https://github.com/openshift/cluster-api-provider-baremetal) git [5d49346b](https://github.com/openshift/cluster-api-provider-baremetal/commit/5d49346b95a2a69a07c186d7cf64a9c3d4e655b4) `sha256:629f9f5739163b64c40a880b11b2df9d83b78dddabe64c62ebcc1ea0fa372e7c`
* [kube-etcd-signer-server](https://github.com/openshift/kubecsr) git [9d3e068b](https://github.com/openshift/kubecsr/commit/9d3e068b3fc7de3df1cb8000492e43efc8200d4b) `sha256:8355fe31ce4eba134caf82df5167c5be778227a1ce7023c258db8d52d4212d02`
* [ovn-kubernetes](https://github.com/openshift/ose-ovn-kubernetes) git [4105e730](https://github.com/openshift/ose-ovn-kubernetes/commit/4105e7303a13114932b7c2d1ad376b5162383c78) `sha256:7cb4d054492be04b648205b714f61bcf3bea0e59e18fd2d4627bc00cd926565a`
* [sriov-cni](https://github.com/openshift/sriov-cni) git [1b427e82](https://github.com/openshift/sriov-cni/commit/1b427e82fdc7ad157d866556ff27ed6495bd5c7e) `sha256:2661cab8b8d351e25c4d475a9d58ccab48db8c3266d6320a7165ccb9d7a6894a`
* [sriov-network-device-plugin](https://github.com/openshift/sriov-network-device-plugin) git [0902e690](https://github.com/openshift/sriov-network-device-plugin/commit/0902e690c82248805d9468ec07c7348fe89fd221) `sha256:93ad0471d28b5e6d264c0d455ef1ba454df29aa6e798e4fe68a029c678a1341c`


### Rebuilt images without code change

* machine-os-content `sha256:d762ceee9f46a141f54cc4dc9689fa19048c1df9b26aae5d8016d6d44995a08d`


### [aws-machine-controllers](https://github.com/openshift/cluster-api-provider-aws)

* Sync with the latest openshift/cluster-api [#181](https://github.com/openshift/cluster-api-provider-aws/pull/181)
* Do not be explicit about iops constraints [#183](https://github.com/openshift/cluster-api-provider-aws/pull/183)
* K8s 1.13 [#185](https://github.com/openshift/cluster-api-provider-aws/pull/185)
* verify staleness of vendored cluster-api-actuator-pkg [#186](https://github.com/openshift/cluster-api-provider-aws/pull/186)
* Drop unused labels from installer machineSets [#182](https://github.com/openshift/cluster-api-provider-aws/pull/182)
* [Full changelog](https://github.com/openshift/cluster-api-provider-aws/compare/995e3e2a6d2b4a06ca07a61279b2131b1e487344...4d953241bc7f62785e0ff9f759315f386e790ba2)


### [cli, deployer, hyperkube, hypershift, node, tests](https://github.com/openshift/ose)

* switch integration tests to external kubeclient [#22382](https://github.com/openshift/ose/pull/22382)
* UPSTREAM: 00000: fix closure error when calling cert rotation Run func [#22395](https://github.com/openshift/ose/pull/22395)
* customresourcevalidation: switch Features to FeatureGate [#22389](https://github.com/openshift/ose/pull/22389)
* extended: validate release image used and image pull policy [#22230](https://github.com/openshift/ose/pull/22230)
* [Bug 1691857](https://bugzilla.redhat.com/show_bug.cgi?id=1691857): Use protobuf from kubelet to retrieve resources, not just send them [#22386](https://github.com/openshift/ose/pull/22386)
* bootstrap-policy: add /readyz to system:discovery [#22399](https://github.com/openshift/ose/pull/22399)
* extended: exit early if unable to read release payload [#22402](https://github.com/openshift/ose/pull/22402)
* UPSTREAM: <carry>: bootstrap-policy: add /readyz system:openshift:public-info-viewer policy [#22401](https://github.com/openshift/ose/pull/22401)
* Retry when checking prometheus working target for router-internal-default job [#22061](https://github.com/openshift/ose/pull/22061)
* change the sdn pod to run the iptables command on the host [#22376](https://github.com/openshift/ose/pull/22376)
* UPSTREAM: <carry>: allow running bare kube-controller-manager [#22398](https://github.com/openshift/ose/pull/22398)
* minor fix in "oc get --help" text re. "list of common resources" [#22372](https://github.com/openshift/ose/pull/22372)
* UPSTREAM: 71804: Use UnmountMountPoint util to clean up subpaths [#22396](https://github.com/openshift/ose/pull/22396)
* Build kube-proxy image [#22155](https://github.com/openshift/ose/pull/22155)
* Bump (*) deps [#22391](https://github.com/openshift/ose/pull/22391)
* Revert "UPSTREAM: 71804: Use UnmountMountPoint util to clean up subpaths" [#22412](https://github.com/openshift/ose/pull/22412)
* UPSTREAM: <drop>: enable secure serving for scheduler [#22414](https://github.com/openshift/ose/pull/22414)
* UPSTREAM: 71804: Use UnmountMountPoint util to clean up subpaths [#22417](https://github.com/openshift/ose/pull/22417)
* [Bug 1689061](https://bugzilla.redhat.com/show_bug.cgi?id=1689061): Report build failures due to eviction [#22344](https://github.com/openshift/ose/pull/22344)
* operators: pull release image with pull secret from cluster [#22408](https://github.com/openshift/ose/pull/22408)
* Check NewForConfig err in WhoAmI [#22415](https://github.com/openshift/ose/pull/22415)
* Removing utf-8 file in extended build test [#22421](https://github.com/openshift/ose/pull/22421)
* e2e: delete worker machines and ensure recovery [#22090](https://github.com/openshift/ose/pull/22090)
* Temporarily disable RBR until we move it to a CRD [#22416](https://github.com/openshift/ose/pull/22416)
* Clear conntrack entries for externalIPs [#22345](https://github.com/openshift/ose/pull/22345)
* [Full changelog](https://github.com/openshift/ose/compare/461e7d39741f996fad13203ccdc8c1a55ad6c44a...048c5dbc3a3df0a9004d395d1df1ab67bb47bfd2)


### [cloud-credential-operator](https://github.com/openshift/cloud-credential-operator)

* remove v1beta1 cloudcredentials objects [#46](https://github.com/openshift/cloud-credential-operator/pull/46)
* [Bug 1689442](https://bugzilla.redhat.com/show_bug.cgi?id=1689442): Rename ClusterOperator to cloud-credential. [#47](https://github.com/openshift/cloud-credential-operator/pull/47)
* Remove hardcoded amd64 for multi-arch [#48](https://github.com/openshift/cloud-credential-operator/pull/48)
* [Full changelog](https://github.com/openshift/cloud-credential-operator/compare/2560a997b6712c240339a92109780ea36b9cf30f...9c0803e2da486c102d055f557a67f8cb148e849f)


### [cluster-authentication-operator](https://github.com/openshift/cluster-authentication-operator)

* redeploy on operatorConfig.spec changes [#98](https://github.com/openshift/cluster-authentication-operator/pull/98)
* Improve error messages [#101](https://github.com/openshift/cluster-authentication-operator/pull/101)
* Tighten APIs [#99](https://github.com/openshift/cluster-authentication-operator/pull/99)
* [Bug 1691488](https://bugzilla.redhat.com/show_bug.cgi?id=1691488): Fix route healthz to tolerate custom wildcard certs [#102](https://github.com/openshift/cluster-authentication-operator/pull/102)
* [Full changelog](https://github.com/openshift/cluster-authentication-operator/compare/a0b0f41eadd3bd153a1d1ac20c7fceb3f9921f05...522c921df8db8fec15ed303b591a9f0db40cc42d)


### [cluster-autoscaler](https://github.com/openshift/kubernetes-autoscaler)

* UPSTREAM: <carry>: openshift add unit tests to resize nodegroup [#77](https://github.com/openshift/kubernetes-autoscaler/pull/77)
* UPSTREAM: <carry>: openshift: Add fmt, lint, vet scripts/Makefile [#75](https://github.com/openshift/kubernetes-autoscaler/pull/75)
* [Full changelog](https://github.com/openshift/kubernetes-autoscaler/compare/3b6f5dfa8bf38e7b0b0dc6c19f4b79fccd3eef0f...d9417ed48ee888d97839bcc8ab2f0eee56ff3206)


### [cluster-autoscaler-operator](https://github.com/openshift/cluster-autoscaler-operator)

* Add additional logging [#77](https://github.com/openshift/cluster-autoscaler-operator/pull/77)
* verify staleness of vendored cluster-api-actuator-pkg [#78](https://github.com/openshift/cluster-autoscaler-operator/pull/78)
* Make StatusReporter a Runnable added to the manager [#79](https://github.com/openshift/cluster-autoscaler-operator/pull/79)
* [Full changelog](https://github.com/openshift/cluster-autoscaler-operator/compare/54672e5e562be1c55db424a744da445fc715f7b3...61c6966085488d29ac2af4edaab8385f167e0b67)


### [cluster-bootstrap](https://github.com/openshift/cluster-bootstrap)

* manifests/image-references: add kube-etcd-signer-server [#22](https://github.com/openshift/cluster-bootstrap/pull/22)
* [Full changelog](https://github.com/openshift/cluster-bootstrap/compare/fb8d5d44671462d9356a8a79ce94f5eea4045c55...9cc6349d2963336c500356ffcb03859a186fa738)


### [cluster-config-operator](https://github.com/openshift/cluster-config-operator)

* Remove regitry CRD, bump library-go and update generated CRD's descriptions [#27](https://github.com/openshift/cluster-config-operator/pull/27)
* bump deps, update oauth crd and infrastructure crd [#29](https://github.com/openshift/cluster-config-operator/pull/29)
* [Full changelog](https://github.com/openshift/cluster-config-operator/compare/db0343ca7ea7ea4b34b3499a7341e6aa106c5ae9...66006332df09a182380e46e1bc1af32a91fee759)


### [cluster-image-registry-operator](https://github.com/openshift/cluster-image-registry-operator)

* api: Replace time.Duration by metav1.Duration [#245](https://github.com/openshift/cluster-image-registry-operator/pull/245)
* Add reasons to conditions [#246](https://github.com/openshift/cluster-image-registry-operator/pull/246)
* [Full changelog](https://github.com/openshift/cluster-image-registry-operator/compare/f5370d9a9d5f55b2b75bd99968b704b8ce24e08c...a93c2103454dfd8a898ddd82fec150e334589308)


### [cluster-ingress-operator](https://github.com/openshift/cluster-ingress-operator)

* deploymentConfigChanged: Fix volume comparison [#182](https://github.com/openshift/cluster-ingress-operator/pull/182)
* manifests: Simplify RouterNamespace [#180](https://github.com/openshift/cluster-ingress-operator/pull/180)
* [Bug 1683765](https://bugzilla.redhat.com/show_bug.cgi?id=1683765): Ensure IngressController domain is unique [#175](https://github.com/openshift/cluster-ingress-operator/pull/175)
* [Bug 1690333](https://bugzilla.redhat.com/show_bug.cgi?id=1690333): Implement IngressController Available Status Condition [#174](https://github.com/openshift/cluster-ingress-operator/pull/174)
* Remove unused method in manifests [#184](https://github.com/openshift/cluster-ingress-operator/pull/184)
* Use existing RouterDeploymentName() in router deployment controller [#188](https://github.com/openshift/cluster-ingress-operator/pull/188)
* Add instructions for building remotely on a cluster [#181](https://github.com/openshift/cluster-ingress-operator/pull/181)
* Remove unnecessary operator config in manifests factory [#189](https://github.com/openshift/cluster-ingress-operator/pull/189)
* Fix error handling in currentRouterDeployment method [#191](https://github.com/openshift/cluster-ingress-operator/pull/191)
* haproxy: set nbthread to 4 by default [#190](https://github.com/openshift/cluster-ingress-operator/pull/190)
* Dep updates to support generating Ingress CRD [#185](https://github.com/openshift/cluster-ingress-operator/pull/185)
* RouterDeployment method no longer need IngressController parameter [#193](https://github.com/openshift/cluster-ingress-operator/pull/193)
* [Full changelog](https://github.com/openshift/cluster-ingress-operator/compare/e49c483cea90d0360ce653afdc8104e145d67123...2492ac0b17d0deedf01ffefab6fa49a3aecabc9a)


### [cluster-kube-apiserver-operator](https://github.com/openshift/cluster-kube-apiserver-operator)

* pod.yaml: add init container waiting for free port [#346](https://github.com/openshift/cluster-kube-apiserver-operator/pull/346)
* Increase minimal-shutdown-duration to 35s to cope with slowly converging SDN [#352](https://github.com/openshift/cluster-kube-apiserver-operator/pull/352)
* Disallow Unmanaged+Removed for management state [#325](https://github.com/openshift/cluster-kube-apiserver-operator/pull/325)
* audit-policy: exclude /readyz [#353](https://github.com/openshift/cluster-kube-apiserver-operator/pull/353)
* operator/v1: fix pattern for managementState [#354](https://github.com/openshift/cluster-kube-apiserver-operator/pull/354)
* Revert render changes from 0b686ff00295c382f245b0b4103a566d672498c8 [#356](https://github.com/openshift/cluster-kube-apiserver-operator/pull/356)
* Add e2e test for user client-ca [#357](https://github.com/openshift/cluster-kube-apiserver-operator/pull/357)
* Move certrotation configmap to openshift-config namespace [#359](https://github.com/openshift/cluster-kube-apiserver-operator/pull/359)
* [Full changelog](https://github.com/openshift/cluster-kube-apiserver-operator/compare/f9309586c338c9b8332a401f6c9393a42b3c6054...1ef690e485222e559ea0aa27388abcb2620fc372)


### [cluster-kube-controller-manager-operator](https://github.com/openshift/cluster-kube-controller-manager-operator)

* Remove cluster api dependency [#196](https://github.com/openshift/cluster-kube-controller-manager-operator/pull/196)
* library-go: bump to fix pruning and condition reporting [#198](https://github.com/openshift/cluster-kube-controller-manager-operator/pull/198)
* pod.yaml: add init container waiting for free port [#197](https://github.com/openshift/cluster-kube-controller-manager-operator/pull/197)
* disable kcm secure serving [#200](https://github.com/openshift/cluster-kube-controller-manager-operator/pull/200)
* Generate crd schema [#186](https://github.com/openshift/cluster-kube-controller-manager-operator/pull/186)
* cloudprovider: add vsphere provider [#204](https://github.com/openshift/cluster-kube-controller-manager-operator/pull/204)
* Add cert-rotation configmap for setting base rotation interval [#203](https://github.com/openshift/cluster-kube-controller-manager-operator/pull/203)
* Shorten wait time for the free port [#201](https://github.com/openshift/cluster-kube-controller-manager-operator/pull/201)
* [Full changelog](https://github.com/openshift/cluster-kube-controller-manager-operator/compare/6e4acedf260f5339237fb1049d69bea65cc3d322...02bdf056832992818ccdd7f5a896d7445468f7d2)


### [cluster-kube-scheduler-operator](https://github.com/openshift/cluster-kube-scheduler-operator)

* bump(*): library-go check for revision status 0 [#74](https://github.com/openshift/cluster-kube-scheduler-operator/pull/74)
* Generate CRD schemas [#67](https://github.com/openshift/cluster-kube-scheduler-operator/pull/67)
* Fix deps [#84](https://github.com/openshift/cluster-kube-scheduler-operator/pull/84)
* force scheduler insecure [#80](https://github.com/openshift/cluster-kube-scheduler-operator/pull/80)
* [Full changelog](https://github.com/openshift/cluster-kube-scheduler-operator/compare/12ced252a7cd2b17ff4df9b45c009fbe55451e44...f8d432d051b991a42d6e664b46bd9b05863c4ced)


### [cluster-monitoring-operator](https://github.com/openshift/cluster-monitoring-operator)

* pkg/client/condition: respect current and target version [#295](https://github.com/openshift/cluster-monitoring-operator/pull/295)
* jsonnet: bump telemeter [#296](https://github.com/openshift/cluster-monitoring-operator/pull/296)
* Don't assign users to node-exporter SCC directly [#297](https://github.com/openshift/cluster-monitoring-operator/pull/297)
* jsonnet: Decrease KubeClientCertificateExpiration expiration threshold [#275](https://github.com/openshift/cluster-monitoring-operator/pull/275)
* test/e2e: check if telemeter-client is present [#299](https://github.com/openshift/cluster-monitoring-operator/pull/299)
* [Full changelog](https://github.com/openshift/cluster-monitoring-operator/compare/1010aef15fde299122bbaee787eade6b8c973e93...846a61077598df7f5ee35fe50db29f006abbba29)


### [cluster-network-operator](https://github.com/openshift/cluster-network-operator)

* Remove hardcoded 'amd64' from Dockerfile to support multi-arch builds. [#131](https://github.com/openshift/cluster-network-operator/pull/131)
* add additional NetworkType openshiftSRIOV [#84](https://github.com/openshift/cluster-network-operator/pull/84)
* DHCP Daemon [#130](https://github.com/openshift/cluster-network-operator/pull/130)
* Initial support for ovn-kubernetes [#44](https://github.com/openshift/cluster-network-operator/pull/44)
* fix mount path for sriov device plugin socket [#134](https://github.com/openshift/cluster-network-operator/pull/134)
* [Full changelog](https://github.com/openshift/cluster-network-operator/compare/535239c80d41e62b6bacb981131854d97a613293...433e5a322b14790bd2bb97cf018ce6e1aae6bf0a)


### [cluster-openshift-apiserver-operator](https://github.com/openshift/cluster-openshift-apiserver-operator)

* Generate CRD schema [#158](https://github.com/openshift/cluster-openshift-apiserver-operator/pull/158)
* library-go: bump to pick certs, managed state handling and status reporting [#171](https://github.com/openshift/cluster-openshift-apiserver-operator/pull/171)
* [Full changelog](https://github.com/openshift/cluster-openshift-apiserver-operator/compare/aa5be8acac46963da85d538a04e6aed58effcfed...c170817cf98ffe73a25b91c9f49a6c4a3b2ed0a9)


### [cluster-openshift-controller-manager-operator](https://github.com/openshift/cluster-openshift-controller-manager-operator)

* Generate CRD schema [#79](https://github.com/openshift/cluster-openshift-controller-manager-operator/pull/79)
* [Full changelog](https://github.com/openshift/cluster-openshift-controller-manager-operator/compare/69639ade927d81f62d3db21e66976cbe0f86d53b...fa7596cae195f29dbc6f313501711aa693f31cc8)


### [cluster-storage-operator](https://github.com/openshift/cluster-storage-operator)

* Set progressing status just before trying to create SC [#22](https://github.com/openshift/cluster-storage-operator/pull/22)
* Read Infrastructure instead of cluster-config-v1 configmap [#23](https://github.com/openshift/cluster-storage-operator/pull/23)
* [Full changelog](https://github.com/openshift/cluster-storage-operator/compare/382f03f877f5aa2cf79cec205ccb27a2ab5b5a0b...2fd0a79105668422ebe7c1c5363e764e07f4d912)


### [cluster-svcat-apiserver-operator](https://github.com/openshift/cluster-svcat-apiserver-operator)

* add liveness & readiness probes [#41](https://github.com/openshift/cluster-svcat-apiserver-operator/pull/41)
* Bump and crd cleanup [#42](https://github.com/openshift/cluster-svcat-apiserver-operator/pull/42)
* [Full changelog](https://github.com/openshift/cluster-svcat-apiserver-operator/compare/843601592ae408201d2eb5b57e71bbb60781f40c...0e8e9ffc5e874be13ee99c471216f930575c083b)


### [cluster-svcat-controller-manager-operator](https://github.com/openshift/cluster-svcat-controller-manager-operator)

* add liveness & readiness probes [#26](https://github.com/openshift/cluster-svcat-controller-manager-operator/pull/26)
* crd schema validation & vendor bump [#29](https://github.com/openshift/cluster-svcat-controller-manager-operator/pull/29)
* create operand ServiceMonitor during sync [#28](https://github.com/openshift/cluster-svcat-controller-manager-operator/pull/28)
* [Full changelog](https://github.com/openshift/cluster-svcat-controller-manager-operator/compare/83ec5b9fc15ba35db0cad5a395bc4dbe4e9b4b4e...f8ad3be67ae5baaa8c2c5d167ba48611146df560)


### [cluster-version-operator](https://github.com/openshift/cluster-version-operator)

* [Bug 1670727](https://bugzilla.redhat.com/show_bug.cgi?id=1670727): install/0000_00_cluster-version-operator_03_deployment: Set 'strategy: Recreate' [#140](https://github.com/openshift/cluster-version-operator/pull/140)
* test: Avoid collisions on namespace and lengthen shutdown timeout [#143](https://github.com/openshift/cluster-version-operator/pull/143)
* [Bug 1691513](https://bugzilla.redhat.com/show_bug.cgi?id=1691513): CVO should retry initialization until cancelled to avoid wedging because of failing dependencies [#141](https://github.com/openshift/cluster-version-operator/pull/141)
* leaderelect: Prevent double close channel panic when shutting down lease [#144](https://github.com/openshift/cluster-version-operator/pull/144)
* BUG 1693051: install: drop deprecated crd [#149](https://github.com/openshift/cluster-version-operator/pull/149)
* status: Correctly increment history and don't overwrite controller cache objects [#146](https://github.com/openshift/cluster-version-operator/pull/146)
* Remove hardcoded 'amd64' from Dockerfile to support multi-arch builds. [#147](https://github.com/openshift/cluster-version-operator/pull/147)
* cvo: fix max_workers to number of nodes in graph [#151](https://github.com/openshift/cluster-version-operator/pull/151)
* [Full changelog](https://github.com/openshift/cluster-version-operator/compare/5e907f5800858bc3298308496b15a6a996051cf5...b72cf639fd562bd76d7025a60e074b3817974ef9)


### [console](https://github.com/openshift/console)

* Monitoring: Improve graph hover labels to only show those that differ [#1321](https://github.com/openshift/console/pull/1321)
* [Bug 1633127](https://bugzilla.redhat.com/show_bug.cgi?id=1633127): Use SI prefixes in graphs [#1323](https://github.com/openshift/console/pull/1323)
* Fixed prow e2e Test flakes [#1286](https://github.com/openshift/console/pull/1286)
* [Bug 1691694](https://bugzilla.redhat.com/show_bug.cgi?id=1691694): Avoid runtime error when event reason missing [#1322](https://github.com/openshift/console/pull/1322)
* Adding co-resource-link__resource-api class and rules so that active … [#1326](https://github.com/openshift/console/pull/1326)
* Consume New Marketplace APIs [#1327](https://github.com/openshift/console/pull/1327)
* Monitoring: Enable y axis zoom for alert and alerting rule graphs [#1329](https://github.com/openshift/console/pull/1329)
* Monitoring: Use a default graph span that matches a dropdown option [#1331](https://github.com/openshift/console/pull/1331)
* [Bug 1691602](https://bugzilla.redhat.com/show_bug.cgi?id=1691602): should hide password when configure default pull secret for project [#1333](https://github.com/openshift/console/pull/1333)
* Monitoring: Use formatPrometheusDuration() for alerting rule's `for` [#1330](https://github.com/openshift/console/pull/1330)
* Monitoring: Fix graph time span dropdown to align with the text input [#1334](https://github.com/openshift/console/pull/1334)
* Added npm scripts to debug e2e and unit tests [#1317](https://github.com/openshift/console/pull/1317)
* Monitoring: Add Prometheus API timeout option [#1335](https://github.com/openshift/console/pull/1335)
* Fix scenario 'logs in via htpasswd identity provider' for username [#1338](https://github.com/openshift/console/pull/1338)
* Update console for identity provider API changes [#1337](https://github.com/openshift/console/pull/1337)
* Changes to the storage class overview page [#1339](https://github.com/openshift/console/pull/1339)
* Add OpenIDP Connect Form to Cluster Settings OAuth Page [#1328](https://github.com/openshift/console/pull/1328)
* Add Taints and Tolerations [#1301](https://github.com/openshift/console/pull/1301)
* Rich Preview for Operator InstallPlans [#1343](https://github.com/openshift/console/pull/1343)
* Align secrets createItems with deploymentConfig's [#1345](https://github.com/openshift/console/pull/1345)
* Add cluster version flag check to masthead toolbar to prevent failed cluster version requests [#1347](https://github.com/openshift/console/pull/1347)
* Remove Backwards Compatibility Hack for Marketplace APIs [#1351](https://github.com/openshift/console/pull/1351)
* Remove duplication of dropdown items [#1352](https://github.com/openshift/console/pull/1352)
* Show tolerations in overview resource summaries [#1353](https://github.com/openshift/console/pull/1353)
* Don't automatically add an empty row to taints and tolerations dialogs [#1354](https://github.com/openshift/console/pull/1354)
* Add common component for idp mapping method dropdown [#1348](https://github.com/openshift/console/pull/1348)
* Update OpenID IDP form [#1356](https://github.com/openshift/console/pull/1356)
* Monitoring: Use PatternFly input error style for graph span text input [#1340](https://github.com/openshift/console/pull/1340)
* Monitoring: Fix Graph's Reset Zoom button to also reset the Y axis [#1358](https://github.com/openshift/console/pull/1358)
* Improve consistency of remove buttons [#1355](https://github.com/openshift/console/pull/1355)
* Address issue of resource name on details page doesn't wrap and causes resource icon to break [#1350](https://github.com/openshift/console/pull/1350)
* Fix lint warning in crud integration test scenario [#1349](https://github.com/openshift/console/pull/1349)
* Add field level help for IDP mapping method [#1359](https://github.com/openshift/console/pull/1359)
* Fine tuning search all resource types to fix visual issues [#1357](https://github.com/openshift/console/pull/1357)
* Monitoring: Add show / hide toggle button for the graphs [#1346](https://github.com/openshift/console/pull/1346)
* [Bug 1660785](https://bugzilla.redhat.com/show_bug.cgi?id=1660785): REQUEST LIMITS should be RESOURCE REQUEST on EtcdCluster details page [#1364](https://github.com/openshift/console/pull/1364)
* [Bug 1660797](https://bugzilla.redhat.com/show_bug.cgi?id=1660797): LABEL selector on PVC detail page should be PV selector [#1365](https://github.com/openshift/console/pull/1365)
* Use 'default' namespace when creating PVC using 'all projects' namespace [#1366](https://github.com/openshift/console/pull/1366)
* [Full changelog](https://github.com/openshift/console/compare/4cb747009adf95fc9409714a4df52ff1e7522912...59a5d768394658f83c876c7e3e00198213e77988)


### [console-operator](https://github.com/openshift/console-operator)

* Fix Sprintf "deployment" for SyncDeployment() error [#186](https://github.com/openshift/console-operator/pull/186)
* Handle Console Operator LogLevel [#179](https://github.com/openshift/console-operator/pull/179)
* Add infrastructures config to operator related objects [#185](https://github.com/openshift/console-operator/pull/185)
* Use const in configmap test, eliminate stale comments [#184](https://github.com/openshift/console-operator/pull/184)
* Support Spec.UnsupportedConfigOverrides [#162](https://github.com/openshift/console-operator/pull/162)
* Remove hardcoded amd64 from Dockerfile for MultiArch builds [#187](https://github.com/openshift/console-operator/pull/187)
* Split roles into 3 files for clarity & to avoid mistakes [#188](https://github.com/openshift/console-operator/pull/188)
* [Full changelog](https://github.com/openshift/console-operator/compare/e57e22533cccddf7f8a2fdf30152643ada5a72f4...4be0f3ee915ab588b2bb74133d1721c7562c4536)


### [docker-builder](https://github.com/openshift/builder)

* Add verbose output for bsdtar [#55](https://github.com/openshift/builder/pull/55)
* Fix empty name in status update [#56](https://github.com/openshift/builder/pull/56)
* [Full changelog](https://github.com/openshift/builder/compare/a49acf5501281e982dcea3f87e4fc228b4d3fa71...8fc01667c3c5813c0a453a2889a744eeef888a6e)


### [installer](https://github.com/openshift/installer)

* docs/user/customization: Catch up with "Creating infrastructure resources" [#1425](https://github.com/openshift/installer/pull/1425)
* upi/vsphere: support rhcos-latest template [#1451](https://github.com/openshift/installer/pull/1451)
* cmd/openshift-install/create: Log progress on timeout too [#1447](https://github.com/openshift/installer/pull/1447)
* bootstrap: Work around systemd-journal-gateway DynamicUser=yes [#1445](https://github.com/openshift/installer/pull/1445)
* images/installer: add image that can be used to instal UPI platforms [#1456](https://github.com/openshift/installer/pull/1456)
* Remove cluster-api cluster object dependency [#1449](https://github.com/openshift/installer/pull/1449)
* release: Allow release image to be directly substituted into binary [#1422](https://github.com/openshift/installer/pull/1422)
* CHANGELOG: Document changes since 0.14.0 [#1462](https://github.com/openshift/installer/pull/1462)
* Re-add support for trunk ports [#1431](https://github.com/openshift/installer/pull/1431)
* [openstack] Support for tagging workers [#1453](https://github.com/openshift/installer/pull/1453)
* rhcos: Bump bootimage to 410.8.20190325.0 [#1459](https://github.com/openshift/installer/pull/1459)
* pkg/types: add vsphere platform [#1458](https://github.com/openshift/installer/pull/1458)
* Re-group the network operator crd. [#1410](https://github.com/openshift/installer/pull/1410)
* machines: add the authorized keys for a pool using a machine config [#1150](https://github.com/openshift/installer/pull/1150)
* [Bug 1659970](https://bugzilla.redhat.com/show_bug.cgi?id=1659970): terraform/exec/plugins/vendor: Bump terraform-provider-aws to v2.2.0 [#1442](https://github.com/openshift/installer/pull/1442)
* Add initial docs and example implementation for UPI bare-metal [#1416](https://github.com/openshift/installer/pull/1416)
* *: use kube-etcd-cert-signer release image [#1477](https://github.com/openshift/installer/pull/1477)
* pkg/types/aws/validation: Require machine-pool zones in platform region [#1469](https://github.com/openshift/installer/pull/1469)
* chore: Bumping protobuf version to 1.3.1 [#1482](https://github.com/openshift/installer/pull/1482)
* Fixed DNS record description for OpenStack [#1435](https://github.com/openshift/installer/pull/1435)
* upi/vsphere: create apps dns entry [#1476](https://github.com/openshift/installer/pull/1476)
* pkg/asset/machines/worker: Structured workers [#1481](https://github.com/openshift/installer/pull/1481)
* machines: fix machine asset for None, VSphere platform [#1493](https://github.com/openshift/installer/pull/1493)
* aws/permissions: Add s3:GetBucketObjectLockConfiguration to pre-flight checks [#1491](https://github.com/openshift/installer/pull/1491)
* Openstack doc [#1486](https://github.com/openshift/installer/pull/1486)
* Dockerfile.upi.ci: declari cli image as input image [#1497](https://github.com/openshift/installer/pull/1497)
* [Full changelog](https://github.com/openshift/installer/compare/0d891e1119806555871330723497f4ac770ad13a...58a27678c02d28693306e1ff5ff8480d7c9fb84c)


### [jenkins, jenkins-agent-maven, jenkins-agent-nodejs](https://github.com/openshift/jenkins)

* 🔒 Upgrade plugins for Jenkins security advisory 2019-03-25 [#826](https://github.com/openshift/jenkins/pull/826)
* fixes around tbr, 3.11 sec advisories, agent default img [#828](https://github.com/openshift/jenkins/pull/828)
* [Full changelog](https://github.com/openshift/jenkins/compare/1ed72f020c8246d051c378ac82d09d4a88a92d7a...29e953b6fbb6b436747a772755946b90691a862c)


### [kube-client-agent](https://github.com/openshift/kubecsr)

* cherry-pick: *: add support for metrics signer [#6](https://github.com/openshift/kubecsr/pull/6)
* [Full changelog](https://github.com/openshift/kubecsr/compare/ae8f7b57f379689fe3ee412bd57c7cd0e0ef9023...9d3e068b3fc7de3df1cb8000492e43efc8200d4b)


### [libvirt-machine-controllers](https://github.com/openshift/cluster-api-provider-libvirt)

* Vendor latest cluster-api changes [#135](https://github.com/openshift/cluster-api-provider-libvirt/pull/135)
* Update example files to providerSpec [#136](https://github.com/openshift/cluster-api-provider-libvirt/pull/136)
* [Full changelog](https://github.com/openshift/cluster-api-provider-libvirt/compare/12a147c2332c3e572862d556654d6bae3e88dcc7...a286b41a60c1ae7e1b02e1abc03e6b8dce81faba)


### [machine-api-operator](https://github.com/openshift/machine-api-operator)

* Add support for managed bare metal platform. [#235](https://github.com/openshift/machine-api-operator/pull/235)
* extend actuator-pkg staleness check [#259](https://github.com/openshift/machine-api-operator/pull/259)
* Introduce azure support [#260](https://github.com/openshift/machine-api-operator/pull/260)
* Don't expose kubemark image to clusters [#263](https://github.com/openshift/machine-api-operator/pull/263)
* operator: use infrastructure.config.openshfit.io for platformtype [#264](https://github.com/openshift/machine-api-operator/pull/264)
* Reuse configv1 for handling platforms [#266](https://github.com/openshift/machine-api-operator/pull/266)
* Drop azure from payload temporally [#267](https://github.com/openshift/machine-api-operator/pull/267)
* [Full changelog](https://github.com/openshift/machine-api-operator/compare/c4a348a011c07db7e3f864380b0134b318112d19...2340bf9b136a8d2e220c7968532992ec12d48ebe)


### [machine-config-controller, machine-config-daemon, machine-config-operator, machine-config-server, setup-etcd-environment](https://github.com/openshift/machine-config-operator)

* Add validation for blocked registries [#569](https://github.com/openshift/machine-config-operator/pull/569)
* docs: Directories and links aren't supported [#568](https://github.com/openshift/machine-config-operator/pull/568)
* Use CRI-O’s pause_image_auth_file option [#540](https://github.com/openshift/machine-config-operator/pull/540)
* Add FeatureGate support to Kubelet Config Controller [#553](https://github.com/openshift/machine-config-operator/pull/553)
* [Bug 1691660](https://bugzilla.redhat.com/show_bug.cgi?id=1691660): pkg/controller: always use the OSImageURL from the CVO [#475](https://github.com/openshift/machine-config-operator/pull/475)
* pkg/controller: get ControllerConfig directly [#572](https://github.com/openshift/machine-config-operator/pull/572)
* consistenly check MC fragments version [#490](https://github.com/openshift/machine-config-operator/pull/490)
* clusterrole: add featuregates to allowed resources [#574](https://github.com/openshift/machine-config-operator/pull/574)
* controller/node: Support master/worker combined, and 1 custom role [#575](https://github.com/openshift/machine-config-operator/pull/575)
* pkg/operator: fix log messages [#573](https://github.com/openshift/machine-config-operator/pull/573)
* [Bug 1677198](https://bugzilla.redhat.com/show_bug.cgi?id=1677198): pkg/daemon: allow empty new links section [#580](https://github.com/openshift/machine-config-operator/pull/580)
* controller/bootstrap: use files with multiple yaml documents [#577](https://github.com/openshift/machine-config-operator/pull/577)
* operator: remove the generated ssh key machineconfig [#356](https://github.com/openshift/machine-config-operator/pull/356)
* mcc: fix spelling of vsphere and update golden files [#589](https://github.com/openshift/machine-config-operator/pull/589)
* kubeletconfig_features: fix ignition object creation [#578](https://github.com/openshift/machine-config-operator/pull/578)
* pkg/operator: fix struct doc [#587](https://github.com/openshift/machine-config-operator/pull/587)
* BUG 1670700: *: add etcd-metric-ca flag [#590](https://github.com/openshift/machine-config-operator/pull/590)
* Remove amd64 from dockerfile copy path to support Mulit-Arch Builds. [#586](https://github.com/openshift/machine-config-operator/pull/586)
* ContainerRuntimeConfig: Enable "info" debug level [#592](https://github.com/openshift/machine-config-operator/pull/592)
* [Full changelog](https://github.com/openshift/machine-config-operator/compare/9da96326a5ff737869709f5fa2e6c716df4dbaf4...ee4e0643716780def1aadaaf7ef4f67b2f99cd7f)


### [multus-cni](https://github.com/openshift/multus-cni)

* add kubelet client for Pod resource info [#9](https://github.com/openshift/multus-cni/pull/9)
* fix missing deviceID in NetConfList [#10](https://github.com/openshift/multus-cni/pull/10)
* [Full changelog](https://github.com/openshift/multus-cni/compare/02913fb96f1d5ad1244a3970d1aa8fc33648fa0c...a805f63a045589e10f98cc1e2107b3a31e386f00)


### [oauth-proxy](https://github.com/openshift/oauth-proxy)

* Don't segfault on wrong upstream URL [#109](https://github.com/openshift/oauth-proxy/pull/109)
* [Full changelog](https://github.com/openshift/oauth-proxy/compare/40c12481bfdd3e87d133736351c907000d5759b2...2653eabc8010e979fca3cd88235ec68f39835449)


### [openstack-machine-controllers](https://github.com/openshift/cluster-api-provider-openstack)

* Rebase from upstream to bring in: Tags and minor upstream fixes [#28](https://github.com/openshift/cluster-api-provider-openstack/pull/28)
* [Full changelog](https://github.com/openshift/cluster-api-provider-openstack/compare/9e913e83ca639e7f6e10fdffa8445f504b101f3c...137fd3c3610e5fc616337f952b9b717a7f32f5e4)


### [operator-lifecycle-manager](https://github.com/operator-framework/operator-lifecycle-manager)

* chore(build): build binaries locally and load them into containers [#777](https://github.com/operator-framework/operator-lifecycle-manager/pull/777)
* chore(catalog): add ClusterOperator status for catalog operator [#776](https://github.com/operator-framework/operator-lifecycle-manager/pull/776)
* fix(olm): use hashes for provided api labels [#778](https://github.com/operator-framework/operator-lifecycle-manager/pull/778)
* fix(metrics): add service monitor config [#682](https://github.com/operator-framework/operator-lifecycle-manager/pull/682)
* chore(api): bump operatorgroup and packagemanifest to v1 [#774](https://github.com/operator-framework/operator-lifecycle-manager/pull/774)
* feat(packageserver): support apps.redhat.com and operators.coreos.com [#788](https://github.com/operator-framework/operator-lifecycle-manager/pull/788)
* fix(packageserver): add struct tags to PackageManifestList [#791](https://github.com/operator-framework/operator-lifecycle-manager/pull/791)
* fix(annotation): don't annotate deployments that aren't owned by a CSV [#792](https://github.com/operator-framework/operator-lifecycle-manager/pull/792)
* Update Kubernetes dependencies to 1.12.7 [#793](https://github.com/operator-framework/operator-lifecycle-manager/pull/793)
* chore(deploy): set priorityclass on olm pods [#775](https://github.com/operator-framework/operator-lifecycle-manager/pull/775)
* docs(subscriptions): add improved sub status proposal [#741](https://github.com/operator-framework/operator-lifecycle-manager/pull/741)
* [Full changelog](https://github.com/operator-framework/operator-lifecycle-manager/compare/9ee66d01eec216d7e1548e4fe2176a5f1f5161ac...9ba3512c5406b62179968e2432b284e9a30c321e)


### [operator-marketplace](https://github.com/operator-framework/operator-marketplace)

* Bug 1691241 [csc] Prevent getting stuck in current state when cache is invalid for csc [#141](https://github.com/operator-framework/operator-marketplace/pull/141)
* [Bug 1690592](https://bugzilla.redhat.com/show_bug.cgi?id=1690592): Update API group and version [#139](https://github.com/operator-framework/operator-marketplace/pull/139)
* [Full changelog](https://github.com/operator-framework/operator-marketplace/compare/213899cee5dcf37f78522f56196a860be10b1e94...e274d6b40505e977e12061becf27218f5eb717fb)


### [operator-registry](https://github.com/operator-framework/operator-registry)

* Remove OperatorSource dependency [#40](https://github.com/operator-framework/operator-registry/pull/40)
* Permit skipping operators during updates [#41](https://github.com/operator-framework/operator-registry/pull/41)
* [Full changelog](https://github.com/operator-framework/operator-registry/compare/a6882d4a64ecdb45df7ebc64fe55bf23d56b7713...1dfcfb9d4c62bf7a2997fe5db0f3935215710f32)


### [service-ca-operator](https://github.com/openshift/service-ca-operator)

* Fix ERROR: logging before flag.Parse [#39](https://github.com/openshift/service-ca-operator/pull/39)
* [Full changelog](https://github.com/openshift/service-ca-operator/compare/6465c13a1ace51c47ccda4a0332682db14e7120a...b97531fe8281de7975fbcb9f503d84631d69bf00)


### [telemeter](https://github.com/openshift/telemeter)

* jsonnet: add session affinity to prom svc [#145](https://github.com/openshift/telemeter/pull/145)
* pkg/metricfamily: add elide transformer [#146](https://github.com/openshift/telemeter/pull/146)
* jsonnet/telemeter: elide prometheus_replica label [#147](https://github.com/openshift/telemeter/pull/147)
* pkg/forwarder: fix match flag [#148](https://github.com/openshift/telemeter/pull/148)
* pkg/authorize/tollbooth: handle 404 - cluster not found [#150](https://github.com/openshift/telemeter/pull/150)
* [Full changelog](https://github.com/openshift/telemeter/compare/c8f63cf406225c6d6f38af0b04e44b266e751391...efc4bdd10d6fa8b790e8d948ea952c2931dc70bd)


