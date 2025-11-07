# conformance.mk
# Do Not Rename this file!
# Make configuration that effects conformance E2E behaviors should go in here,
# to allow us to maintain the core Makefile without having to execute
# long-running conformance E2E jobs every time that file changes
CONFORMANCE_FLAVOR ?=
CONFORMANCE_E2E_ARGS ?= -kubetest.config-file=$(KUBETEST_CONF_PATH)
CONFORMANCE_E2E_ARGS += $(E2E_ARGS)

##@ Conformance testing:

.PHONY: test-conformance
test-conformance: ## Run conformance test on workload cluster.
ifeq ($(MANAGER_IMAGE),)
	$(MAKE) test-e2e-skip-push GINKGO_FOCUS="Conformance" E2E_ARGS='$(CONFORMANCE_E2E_ARGS)' CONFORMANCE_FLAVOR='$(CONFORMANCE_FLAVOR)'
else ## If MANAGER_IMAGE is set, use it for the conformance test (test-e2e-skip-push overwrites it).
	$(MAKE) test-e2e-custom-image GINKGO_FOCUS="Conformance" E2E_ARGS='$(CONFORMANCE_E2E_ARGS)' CONFORMANCE_FLAVOR='$(CONFORMANCE_FLAVOR)'
endif

test-conformance-fast: ## Run conformance test on workload cluster using a subset of the conformance suite in parallel.
	$(MAKE) test-conformance CONFORMANCE_E2E_ARGS="-kubetest.config-file=$(KUBETEST_FAST_CONF_PATH) -kubetest.ginkgo-nodes=5 $(E2E_ARGS)"

.PHONY: test-windows-upstream
test-windows-upstream: ## Run windows upstream tests on workload cluster.
ifneq ($(WIN_REPO_URL), )
	curl --retry $(CURL_RETRIES) $(WIN_REPO_URL) -o $(KUBETEST_REPO_LIST_PATH)/custom-repo-list.yaml
endif
	$(MAKE) test-conformance CONFORMANCE_E2E_ARGS="-kubetest.config-file=$(KUBETEST_WINDOWS_CONF_PATH) -kubetest.repo-list-path=$(KUBETEST_REPO_LIST_PATH) $(E2E_ARGS)"
