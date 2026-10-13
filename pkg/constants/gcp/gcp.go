package gcp

const (
	// ProjectNameFmt is the format string for GCP project resource name.
	ProjectNameFmt = "projects/%s"

	// ProjectParentPathFmt is the format string for parent path of a GCP project resource.
	ProjectParentPathFmt = "//cloudresourcemanager.googleapis.com/projects/%s"

	// ClusterIDLabelFmt is the format string for the default label
	// added to the OpenShift created GCP resources.
	ClusterIDLabelFmt = "kubernetes-io-cluster-%s"
)
