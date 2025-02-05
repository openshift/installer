TEST_NAMESPACE=ns-without-topology
TEST_CLUSTER_NAME=mycluster-without-topology

.PHONY: test-cluster-create
test-cluster-create: ## Create a workload cluster without topology
	clusterctl version
	clusterctl config repositories | grep nutanix
	clusterctl generate cluster ${TEST_CLUSTER_NAME} -i nutanix:${LOCAL_PROVIDER_VERSION} --list-variables -v 10
	clusterctl generate cluster ${TEST_CLUSTER_NAME} -i nutanix:${LOCAL_PROVIDER_VERSION} --target-namespace ${TEST_NAMESPACE}  -v 10 > ./${TEST_CLUSTER_NAME}.yaml
	kubectl create ns $(TEST_NAMESPACE) --dry-run=client -oyaml | kubectl apply --server-side -f -
	kubectl apply --server-side -f ./${TEST_CLUSTER_NAME}.yaml

.PHONY: test-cluster-delete
test-cluster-delete: ## Delete workload cluster without topology
	kubectl -n ${TEST_NAMESPACE} delete cluster ${TEST_CLUSTER_NAME} --ignore-not-found

.PHONY: generate-cluster-kubeconfig
generate-cluster-kubeconfig: ## generate kubeconfig of workload cluster without topology
	kubectl -n $(TEST_NAMESPACE) get secret
	kubectl -n ${TEST_NAMESPACE} get secret ${TEST_CLUSTER_NAME}-kubeconfig -o json | jq -r .data.value | base64 --decode > ${TEST_CLUSTER_NAME}.workload.kubeconfig

.PHONY: test-cluster-install-cni
test-cluster-install-cni: generate-cc-cluster-kubeconfig ## install cni on workload cluster without topology
	kubectl --kubeconfig ./${TEST_CLUSTER_NAME}.workload.kubeconfig apply -f https://raw.githubusercontent.com/nutanix-cloud-native/cluster-api-provider-nutanix/main/test/e2e/data/cni/calico/calico.yaml

.PHONY: list-bootstrap-resources
list-bootstrap-resources: ## List resources of bootstrap/management cluster
	kubectl get ns
	kubectl get all --all-namespaces
	kubectl -n capx-system get all
	kubectl -n $(TEST_NAMESPACE) get Cluster,NutanixCluster,Machine,NutanixMachine,KubeAdmControlPlane,MachineHealthCheck,nodes
	kubectl -n capx-system get pod

.PHONY: list-workload-resources
list-workload-resources: ## List resources of workload cluster
	kubectl --kubeconfig ./${TEST_CLUSTER_NAME}.workload.kubeconfig get nodes,ns
