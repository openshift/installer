module github.com/openshift/installer/terraform/terraform

go 1.16

require github.com/hashicorp/terraform v1.0.11

replace (
	github.com/golang/mock v1.5.0 => github.com/golang/mock v1.4.4
	google.golang.org/grpc v1.36.0 => google.golang.org/grpc v1.27.1
	k8s.io/client-go => k8s.io/client-go v0.0.0-20190620085101-78d2af792bab
)
