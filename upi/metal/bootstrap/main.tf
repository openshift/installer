terraform {
  required_providers {
    metal = {
      source = "equinix/metal"
      version = "3.2.2"
    }

    matchbox = {
      source  = "poseidon/matchbox"
      version = "0.5.0"
    }
  }
}
