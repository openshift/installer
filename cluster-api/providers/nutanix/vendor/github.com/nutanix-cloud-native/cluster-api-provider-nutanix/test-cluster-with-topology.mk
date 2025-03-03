TEST_NAMESPACE=ns-topology
TEST_CLUSTER_CLASS_NAME=my-clusterclass
TEST_TOPOLOGY_CLUSTER_NAME=mycluster-with-topology

.PHONY: test-cc-cluster-create
test-cc-cluster-create: cluster-templates ## Create a workload cluster with topology (which uses clusterclass)
	clusterctl generate cluster ${TEST_CLUSTER_CLASS_NAME} --from ./templates/cluster-template-clusterclass.yaml -n $(TEST_NAMESPACE) > ${TEST_CLUSTER_CLASS_NAME}.yaml
	clusterctl generate cluster ${TEST_TOPOLOGY_CLUSTER_NAME} --from ./templates/cluster-template-topology.yaml -n $(TEST_NAMESPACE) > ${TEST_TOPOLOGY_CLUSTER_NAME}.yaml
	kubectl create ns $(TEST_NAMESPACE) --dry-run=client -oyaml | kubectl apply --server-side -f -
	kubectl apply --server-side -f ./${TEST_CLUSTER_CLASS_NAME}.yaml
	kubectl apply --server-side -f ./${TEST_TOPOLOGY_CLUSTER_NAME}.yaml

.PHONY: test-cc-cluster-upgrade
test-cc-cluster-upgrade: ## Upgrade K8S version of cluster with topology (which uses clusterclass)
	clusterctl generate cluster ${TEST_TOPOLOGY_CLUSTER_NAME} --from ./templates/cluster-template-topology.yaml -n $(TEST_NAMESPACE) --kubernetes-version=${UPGRADE_K8S_VERSION_TO} > ${TEST_TOPOLOGY_CLUSTER_NAME}.yaml
	kubectl apply --server-side -f ./${TEST_TOPOLOGY_CLUSTER_NAME}.yaml

.PHONY: test-cc-cluster-delete
test-cc-cluster-delete: ## Delete workload cluster with topology
	kubectl -n $(TEST_NAMESPACE) delete cluster ${TEST_TOPOLOGY_CLUSTER_NAME} --ignore-not-found
	kubectl -n $(TEST_NAMESPACE) delete secret ${TEST_TOPOLOGY_CLUSTER_NAME} --ignore-not-found
	kubectl -n $(TEST_NAMESPACE) delete cm ${TEST_TOPOLOGY_CLUSTER_NAME}-pc-trusted-ca-bundle --ignore-not-found
	rm -f ${TEST_TOPOLOGY_CLUSTER_NAME}.yaml
	rm -f ${TEST_CLUSTER_CLASS_NAME}.yaml

.PHONY: generate-cc-cluster-kubeconfig
generate-cc-cluster-kubeconfig: ## generate kubeconfig of workload cluster with topology
	kubectl -n ${TEST_NAMESPACE} get secret ${TEST_TOPOLOGY_CLUSTER_NAME}-kubeconfig -o json | jq -r .data.value | base64 --decode > ${TEST_TOPOLOGY_CLUSTER_NAME}.workload.kubeconfig

.PHONY: test-cc-cluster-install-cni
test-cc-cluster-install-cni: generate-cc-cluster-kubeconfig ## install cni on workload cluster with topology
	kubectl --kubeconfig ./${TEST_TOPOLOGY_CLUSTER_NAME}.workload.kubeconfig apply -f https://raw.githubusercontent.com/nutanix-cloud-native/cluster-api-provider-nutanix/main/test/e2e/data/cni/calico/calico.yaml

.PHONY: list-cc-cluster-resources
list-cc-cluster-resources: generate-cc-cluster-kubeconfig ## list resources of workload cluster with topology
	kubectl -n capx-system get endpoints
	kubectl get crd | grep nutanix
	kubectl get cluster-api -A
	clusterctl describe cluster ${TEST_TOPOLOGY_CLUSTER_NAME} -n ${TEST_NAMESPACE}
	kubectl -n $(TEST_NAMESPACE) get Cluster,NutanixCluster,Machine,NutanixMachine,KubeAdmControlPlane,machinedeployments,MachineHealthCheck,nodes
	kubectl get ValidatingWebhookConfiguration,MutatingWebhookConfiguration -A
	kubectl --kubeconfig ./${TEST_TOPOLOGY_CLUSTER_NAME}.workload.kubeconfig get nodes,ns
	kubectl --kubeconfig ./${TEST_TOPOLOGY_CLUSTER_NAME}.workload.kubeconfig get pods -A