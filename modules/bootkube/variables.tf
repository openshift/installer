// The assets directory.
variable assets_dir {
  type = "string"
}

variable trigger_ids {
  type = "list"
}

// The private key for the core user.
// This is used to copy assets to the initial master node.
variable core_private_key {
  type = "string"
}

variable hosts {
  type = "list"
}
