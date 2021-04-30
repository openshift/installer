# NOTE: May not be required in OpenShift Installer
# Provider maybe imported in the installer Go code
terraform {
  required_providers {
    ibm = {
      source = "IBM-Cloud/ibm"
      version = "1.24.0"
    }
  }
}