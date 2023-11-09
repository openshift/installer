default: testacc

# Run acceptance tests
.PHONY: testacc
testacc:
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m

.PHONY: local
local:
	go build && \
	mkdir -p ~/.terraform.d/plugins/terraform.local/local/avgcp/1.0.0/linux_amd64/ && \
	mv terraform-provider-avgcp ~/.terraform.d/plugins/terraform.local/local/avgcp/1.0.0/linux_amd64/terraform-provider-avgcp_v1.0.0