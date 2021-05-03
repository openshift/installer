package powervs


// Platform stores all the global configuration that all machinesets
// use.
/// used by the installconfig, and filled in by the installconfig/platform/powervs::Platform() func
type Platform struct {
	// Region specifies the IBM Cloud region where the cluster will be created.
	Region string `json:"region"`

	// Subnets specifies existing subnets (by ID) where cluster
	// resources will be created.  Leave unset to have the installer
	// create subnets in a new VPC on your behalf.
	// @TODO: how will we handle networking?
	//
	// +optional ?
	Subnets []string `json:"subnets,omitempty"`

	// HostedZone is the ID of an existing hosted zone into which to add DNS
	// records for the cluster's internal API. An existing hosted zone can
	// only be used when also using existing subnets. The hosted zone must be
	// associated with the VPC containing the subnets.
	// Leave the hosted zone unset to have the installer create the hosted zone
	// on your behalf.
	// +optional
	HostedZone string `json:"hostedZone,omitempty"`

	// UserTags additional keys and values that the installer will add
	// as tags to all resources that it creates. Resources created by the
	// cluster itself may not include these tags.
	// +optional
	UserTags map[string]string `json:"userTags,omitempty"`

    // BootstrapOSImage is a URL to override the default OS image                                                                                                                                                                             
    // for the bootstrap node. The URL must contain a sha256 hash of the image                                                                                                                                                                
    // e.g https://mirror.example.com/images/image.ova.gz?sha256=a07bd...                                                                                                                                                                    
    //                                                                                                                                                                                                                                        
    // +optional                                                                                                                                                                                                                              
    BootstrapOSImage string `json:"bootstrapOSImage,omitempty" validate:"omitempty,osimageuri,urlexist"`

    // ClusterOSImage is a URL to override the default OS image                                                                                                                                                                               
    // for cluster nodes. The URL must contain a sha256 hash of the image                                                                                                                                                                     
    // e.g https://mirror.example.com/images/powervs.ova.gz?sha256=3b5a8...                                                                                                                                                                   
    //                                                                                                                                                                                                                                        
    // +optional                                                                                                                                                                                                                              
    ClusterOSImage string `json:"clusterOSImage,omitempty" validate:"omitempty,osimageuri,urlexist"`

}

