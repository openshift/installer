module github.com/openshift/installer/terraform/kubernetes

go 1.16

require github.com/hashicorp/terraform-provider-kubernetes v1.13.3

replace (
	github.com/go-openapi/spec => github.com/go-openapi/spec v0.19.5
	k8s.io/cli-runtime => k8s.io/cli-runtime v0.19.1
	k8s.io/client-go => k8s.io/client-go v0.19.1
	sigs.k8s.io/kustomize/pkg/transformers => ./vendor/k8s.io/cli-runtime/pkg/kustomize/k8sdeps/transformer
	sigs.k8s.io/kustomize/pkg/transformers/config => ./vendor/k8s.io/cli-runtime/pkg/kustomize/k8sdeps/transformer/config
)
