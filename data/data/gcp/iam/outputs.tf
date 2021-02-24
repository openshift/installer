# This simply passes the variable value through this module,
# but enforces a dependency in the master module.
output "master_sa_email" {
  value = var.master_sa_email
}