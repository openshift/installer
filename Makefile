all: test terraform_fmt

#terraform_files=$(shell find modules steps -type f -name "*.tf")
#terraform_files+=("config.tf")

#src_files=$(shell find installer -type f -name "*.go")

test:
	CGO_ENABLED=0 go test ./installer/pkg/...

terraform_fmt:
	terraform fmt -list -check -write=false

#verify-gofmt:
#	./hack/gofmt-all.sh -v

#gofmt:
#	./hack/gofmt-all.sh

#verify: verify-gofmt test
