module github.com/openshift/installer/terraform/ibm

go 1.16

require github.com/IBM-Cloud/terraform-provider-ibm v1.26.2

replace github.com/IBM-Cloud/terraform-provider-ibm => github.com/openshift/terraform-provider-ibm v1.26.2-openshift-2
