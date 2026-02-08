# e2e.mk
# Make configuration that effects E2E behaviors should go in here,
# to allow us to maintain the core Makefile without having to execute
# long-running E2E jobs every time that file changes

##@ E2E Testing:

.PHONY: test-e2e-run
test-e2e-run: generate-e2e-templates install-tools create-bootstrap ## Run e2e tests.
	if [ "$(MGMT_CLUSTER_TYPE)" == "aks" ]; then \
		source ./scripts/peer-vnets.sh && source_tilt_settings tilt-settings.yaml; \
	fi; \
	$(ENVSUBST) < $(E2E_CONF_FILE) > $(E2E_CONF_FILE_ENVSUBST) && \
	if [ -z "${AZURE_CLIENT_ID_USER_ASSIGNED_IDENTITY}" ]; then \
		export AZURE_CLIENT_ID_USER_ASSIGNED_IDENTITY=$(shell cat $(AZURE_IDENTITY_ID_FILEPATH)); \
	fi; \
	$(GINKGO) -v --trace --timeout=4h --tags=e2e --focus="$(GINKGO_FOCUS)" --skip="$(GINKGO_SKIP)" --nodes=$(GINKGO_NODES) --no-color=$(GINKGO_NOCOLOR) --output-dir="$(ARTIFACTS)" --junit-report="junit.e2e_suite.1.xml" $(GINKGO_ARGS) ./test/e2e -- \
		-e2e.artifacts-folder="$(ARTIFACTS)" \
		-e2e.config="$(E2E_CONF_FILE_ENVSUBST)" \
		-e2e.skip-log-collection="$(SKIP_LOG_COLLECTION)" \
		-e2e.skip-resource-cleanup=$(SKIP_CLEANUP) -e2e.use-existing-cluster=$(SKIP_CREATE_MGMT_CLUSTER) $(E2E_ARGS)

.PHONY: test-e2e-run-cleanup
test-e2e-run-cleanup: ## Run e2e cleanup tasks.
	$(MAKE) cleanup-workload-identity || true
	$(MAKE) clean-release-git || true
	if [ "$(MGMT_CLUSTER_TYPE)" == "aks" ] && [ "$(SKIP_CLEANUP)" != "true" ]; then \
		echo "Cleaning up AKS management cluster..."; \
		$(MAKE) aks-delete || true; \
	fi

.PHONY: test-e2e
test-e2e: ## Run "docker-build" and "docker-push" rules then run e2e tests.
	PULL_POLICY=IfNotPresent MANAGER_IMAGE=$(CONTROLLER_IMG)-$(ARCH):$(TAG) \
	$(MAKE) docker-build docker-push \
	test-e2e-run;

.PHONY: test-e2e-skip-push
test-e2e-skip-push: ## Run "docker-build" rule then run e2e tests.
	PULL_POLICY=IfNotPresent MANAGER_IMAGE=$(CONTROLLER_IMG)-$(ARCH):$(TAG) \
	$(MAKE) docker-build \
	test-e2e-run;

.PHONY: test-e2e-skip-build-and-push
test-e2e-skip-build-and-push:
	$(MAKE) set-manifest-image MANIFEST_IMG=$(CONTROLLER_IMG)-$(ARCH) MANIFEST_TAG=$(TAG) TARGET_RESOURCE="./config/capz/manager_image_patch.yaml"
	$(MAKE) set-manifest-pull-policy TARGET_RESOURCE="./config/capz/manager_pull_policy.yaml" PULL_POLICY=IfNotPresent
	MANAGER_IMAGE=$(CONTROLLER_IMG)-$(ARCH):$(TAG) \
	$(MAKE) test-e2e-run;

.PHONY: test-e2e-custom-image
test-e2e-custom-image: ## Run e2e tests with a custom image format (use MANAGER_IMAGE env var).
	@if [ -z "$(MANAGER_IMAGE)" ]; then \
		echo "MANAGER_IMAGE must be set"; \
		exit 1; \
	fi
	$(MAKE) set-manifest-image MANIFEST_IMG=$(shell echo $(MANAGER_IMAGE) | cut -d: -f1) MANIFEST_TAG=$(shell echo $(MANAGER_IMAGE) | cut -d: -f2) TARGET_RESOURCE="./config/capz/manager_image_patch.yaml"
	$(MAKE) set-manifest-pull-policy TARGET_RESOURCE="./config/capz/manager_pull_policy.yaml" PULL_POLICY=IfNotPresent
	$(MAKE) test-e2e-run;
