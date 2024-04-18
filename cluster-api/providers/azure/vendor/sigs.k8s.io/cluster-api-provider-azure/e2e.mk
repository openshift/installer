# e2e.mk
# Make configuration that effects E2E behaviors should go in here,
# to allow us to maintain the core Makefile without having to execute
# long-running E2E jobs every time that file changes

##@ E2E Testing:

.PHONY: test-e2e-run
test-e2e-run: generate-e2e-templates install-tools kind-create-bootstrap ## Run e2e tests.
	$(ENVSUBST) < $(E2E_CONF_FILE) > $(E2E_CONF_FILE_ENVSUBST) && \
    $(GINKGO) -v --trace --timeout=4h --tags=e2e --focus="$(GINKGO_FOCUS)" --skip="$(GINKGO_SKIP)" --nodes=$(GINKGO_NODES) --no-color=$(GINKGO_NOCOLOR) --output-dir="$(ARTIFACTS)" --junit-report="junit.e2e_suite.1.xml" $(GINKGO_ARGS) ./test/e2e -- \
    	-e2e.artifacts-folder="$(ARTIFACTS)" \
    	-e2e.config="$(E2E_CONF_FILE_ENVSUBST)" \
    	-e2e.skip-log-collection="$(SKIP_LOG_COLLECTION)" \
    	-e2e.skip-resource-cleanup=$(SKIP_CLEANUP) -e2e.use-existing-cluster=$(SKIP_CREATE_MGMT_CLUSTER) $(E2E_ARGS)
	$(MAKE) clean-release-git

.PHONY: test-e2e
test-e2e: ## Run "docker-build" and "docker-push" rules then run e2e tests.
	PULL_POLICY=IfNotPresent MANAGER_IMAGE=$(CONTROLLER_IMG)-$(ARCH):$(TAG) \
	$(MAKE) docker-build docker-push \
	test-e2e-run

.PHONY: test-e2e-skip-push
test-e2e-skip-push: ## Run "docker-build" rule then run e2e tests.
	PULL_POLICY=IfNotPresent MANAGER_IMAGE=$(CONTROLLER_IMG)-$(ARCH):$(TAG) \
	$(MAKE) docker-build \
	test-e2e-run

.PHONY: test-e2e-skip-build-and-push
test-e2e-skip-build-and-push:
	$(MAKE) set-manifest-image MANIFEST_IMG=$(CONTROLLER_IMG)-$(ARCH) MANIFEST_TAG=$(TAG) TARGET_RESOURCE="./config/capz/manager_image_patch.yaml"
	$(MAKE) set-manifest-pull-policy TARGET_RESOURCE="./config/capz/manager_pull_policy.yaml" PULL_POLICY=IfNotPresent
	MANAGER_IMAGE=$(CONTROLLER_IMG)-$(ARCH):$(TAG) \
	$(MAKE) test-e2e-run
