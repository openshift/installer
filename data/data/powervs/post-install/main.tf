provider "ibm" {
  alias            = "powervs"
  ibmcloud_api_key = var.powervs_api_key
  region           = var.powervs_region
}
