variable "cluster_name" {
  type = string
}
variable "cluster_basedomain" {
  type = string
}
variable "masters_count" {
  type    = number
  default = 1
}
