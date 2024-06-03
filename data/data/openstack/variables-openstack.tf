variable "openstack_master_root_volume_size" {
  type        = number
  default     = null
  description = "The size of the volume in gigabytes for the root block device of master nodes."
}

variable "openstack_base_image_name" {
  type        = string
  description = "Name of the base image to use for the nodes."
}

variable "openstack_bootstrap_shim_ignition" {
  type        = string
  default     = ""
  description = "Generated pointer/shim ignition config with user ca bundle."
}

variable "openstack_credentials_auth_url" {
  type    = string
  default = ""

  description = <<EOF
required if cloud is not specified) The Identity authentication URL. If omitted, the OS_AUTH_URL environment variable is used.
EOF

}

variable "openstack_credentials_cert" {
  type = string
  default = ""

  description = <<EOF
Specify client certificate file for SSL client authentication. You can specify either a path to the file or the contents of the certificate. If omitted the OS_CERT environment variable is used.
EOF

}

variable "openstack_credentials_cloud" {
  type    = string
  default = ""

  description = <<EOF
required if auth_url is not specified) An entry in a clouds.yaml file. See the openstacksdk(https://docs.openstack.org/openstacksdk/latest/user/config/configuration.html#config-files) documentation for more information about clouds.yaml files. If omitted, the OS_CLOUD environment variable is used.
EOF

}

variable "openstack_credentials_domain_id" {
  type = string
  default = ""

  description = <<EOF
The ID of the Domain to scope to (Identity v3). If omitted, the OS_DOMAIN_ID environment variable is checked.
EOF

}

variable "openstack_credentials_domain_name" {
  type    = string
  default = ""

  description = <<EOF
The Name of the Domain to scope to (Identity v3). If omitted, the following environment variables are checked (in this order): OS_DOMAIN_NAME, OS_DEFAULT_DOMAIN.
EOF

}

variable "openstack_credentials_endpoint_type" {
  type = string
  default = "public"

  description = <<EOF
Specify which type of endpoint to use from the service catalog. It can be set using the OS_ENDPOINT_TYPE environment variable. If not set, public endpoints is used.
EOF

}

variable "openstack_credentials_insecure" {
  default = false

  description = <<EOF
Trust self-signed SSL certificates. If omitted, the OS_INSECURE environment variable is used.
EOF

}

variable "openstack_credentials_key" {
  type = string
  default = ""

  description = <<EOF
Specify client private key file for SSL client authentication. You can specify either a path to the file or the contents of the key. If omitted the OS_KEY environment variable is used.
EOF

}

variable "openstack_credentials_password" {
  type    = string
  default = ""

  description = <<EOF
The Password to login with. If omitted, the OS_PASSWORD environment variable is used.
EOF

}

variable "openstack_credentials_project_domain_id" {
  type = string
  default = ""

  description = <<EOF
The domain ID where the project is located If omitted, the OS_PROJECT_DOMAIN_ID environment variable is checked.
EOF

}

variable "openstack_credentials_project_domain_name" {
  type    = string
  default = ""

  description = <<EOF
The domain name where the project is located. If omitted, the OS_PROJECT_DOMAIN_NAME environment variable is checked.
EOF

}

variable "openstack_credentials_region" {
  type = string
  default = ""

  description = <<EOF
The region of the OpenStack cloud to use. If omitted, the OS_REGION_NAME environment variable is used. If OS_REGION_NAME is not set, then no region will be used. It should be possible to omit the region in single-region OpenStack environments, but this behavior may vary depending on the OpenStack environment being used.
EOF

}

variable "openstack_credentials_swauth" {
  default = false

  description = <<EOF
Set to true to authenticate against Swauth, a Swift-native authentication system. If omitted, the OS_SWAUTH environment variable is used. You must also set username to the Swauth/Swift username such as username:project. Set the password to the Swauth/Swift key. Finally, set auth_url as the location of the Swift service. Note that this will only work when used with the OpenStack Object Storage resources.
EOF

}

variable "openstack_credentials_tenant_id" {
  type = string
  default = ""

  description = <<EOF
The ID of the Tenant (Identity v2) or Project (Identity v3) to login with. If omitted, the OS_TENANT_ID or OS_PROJECT_ID environment variables are used.
EOF

}

variable "openstack_credentials_tenant_name" {
  type    = string
  default = ""

  description = <<EOF
The Name of the Tenant (Identity v2) or Project (Identity v3) to login with. If omitted, the OS_TENANT_NAME or OS_PROJECT_NAME environment variable are used.
EOF

}

variable "openstack_credentials_token" {
  type = string
  default = ""

  description = <<EOF
Required if not using user_name and password) A token is an expiring, temporary means of access issued via the Keystone service. By specifying a token, you do not have to specify a username/password combination, since the token was already created by a username/password out of band of Terraform. If omitted, the OS_TOKEN or OS_AUTH_TOKEN environment variables are used.
EOF

}

variable "openstack_credentials_use_octavia" {
  default = false

  description = <<EOF
If set to true, API requests will go the Load Balancer service (Octavia) instead of the Networking service (Neutron).
EOF

}

variable "openstack_credentials_user_domain_id" {
  type = string
  default = ""

  description = <<EOF
The domain ID where the user is located. If omitted, the OS_USER_DOMAIN_ID environment variable is checked.
EOF

}

variable "openstack_credentials_user_domain_name" {
  type    = string
  default = ""

  description = <<EOF
The domain name where the user is located. If omitted, the OS_USER_DOMAIN_NAME environment variable is checked.
EOF

}

variable "openstack_credentials_user_id" {
  type = string
  default = ""

  description = <<EOF
The User ID to login with. If omitted, the OS_USER_ID environment variable is used.
EOF

}

variable "openstack_credentials_user_name" {
  type    = string
  default = ""

  description = <<EOF
The Username to login with. If omitted, the OS_USERNAME environment variable is used.
EOF

}

variable "openstack_external_network" {
  type = string
  default = ""

  description = <<EOF
(optional) Name of the external network. The network is used to provide
Floating IP access to the deployed nodes. Optional, but either the Name
or UUID option must be specified.
EOF

}

variable "openstack_external_network_id" {
  type    = string
  default = ""

  description = <<EOF
(optional) UUID of the external network. The network is used to provide
Floating IP access to the deployed nodes. Optional, but either the Name
or UUID option must be specified.
EOF

}

variable "openstack_master_extra_sg_ids" {
  type = list(string)
  default = []

  description = <<EOF
(optional) List of additional security group IDs for master nodes.

Example: `["sg-51530134", "sg-b253d7cc"]`
EOF

}

variable "openstack_api_floating_ip" {
  type    = string
  default = ""

  description = <<EOF
(optional) Existing Floating IP to attach to the OpenShift API created by the installer.
EOF

}

variable "openstack_ingress_floating_ip" {
  type = string
  default = ""

  description = <<EOF
(optional) Existing Floating IP to attach to the ingress port created by the installer.
EOF

}

variable "openstack_api_int_ips" {
  type        = list(string)
  description = "IPs on the node subnets reserved for api-int VIP."
}

variable "openstack_ingress_ips" {
  type        = list(string)
  description = "IPs on the nodes subnets reserved for the ingress VIP."
}

variable "openstack_external_dns" {
  type        = list(string)
  description = "IP addresses of exernal dns servers to add to networks."
  default     = []
}

variable "openstack_additional_network_ids" {
  type        = list(string)
  description = "IDs of additional networks for master nodes."
  default     = []
}

variable "openstack_additional_ports" {
  type = list(list(object({
    network_id = string
    fixed_ips = list(object({
      subnet_id  = string
      ip_address = string
    }))
  })))
  description = "Additional ports for each master node."
  default     = [[], [], []]
}

variable "openstack_master_flavor_name" {
  type        = string
  description = "Instance size for the master node(s). Example: `m1.medium`."
}

variable "openstack_octavia_support" {
  type    = bool
  default = false

  description = <<EOF
False if the OpenStack Octavia endpoint is missing and True if it exists.
EOF

}

variable "openstack_master_server_group_name" {
  type = string
  description = "Name of the server group for the master nodes."
}

variable "openstack_master_server_group_policy" {
  type = string
  description = "Policy of the server group for the master nodes."
}

variable "openstack_default_machines_port" {
  type = object({
    network_id = string
    fixed_ips = list(object({
      subnet_id = string
      ip_address = string
    }))
  })
  default = null
  description = "The masters' default control-plane port. If empty, the installer will create a new network."
}

variable "openstack_machines_ports" {
  type = list(object({
    network_id = string
    fixed_ips = list(object({
      subnet_id = string
      ip_address = string
    }))
  }))
  description = "The control-plane port for each machine. If null, the default is used."
  default = [null, null, null]
}

variable "openstack_master_availability_zones" {
  type = list(string)
  default = [""]
  description = "List of availability Zones to Schedule the masters on."
}

variable "openstack_master_root_volume_availability_zones" {
  type = list(string)
  default = [""]
  description = "List of availability Zones to Schedule the masters root volumes on."
}

variable "openstack_master_root_volume_types" {
  type = list(string)
  default = [""]
  description = "List of volume types used by the masters root volumes."
}

variable "openstack_worker_server_group_names" {
  type = set(string)
  default = []
  description = "Names of the server groups for the worker nodes."
}

variable "openstack_worker_server_group_policy" {
  type = string
  description = "Policy of the server groups for the worker nodes."
}

variable "openstack_user_managed_load_balancer" {
  type = bool
  description = "True if the load balancer that is used for the control plane VIPs is managed by the user."
}
