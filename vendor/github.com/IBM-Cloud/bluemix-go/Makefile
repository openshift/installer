.PHONY : test
test: test_deps vet
	go test ./... -timeout 120m

.PHONY : test_deps
test_deps: 
	go get -t ./...

.PHONY : vet

vet:
	@echo 'go vet $$(go list ./... | grep -v vendor)'
	@go vet $$(go list ./... | grep -v vendor) ; if [ $$? -eq 1 ]; then \
		echo ""; \
		echo "Vet found suspicious constructs. Please check the reported constructs"; \
		echo "and fix them if necessary before submitting the code for review."; \
		exit 1; \
	fi
