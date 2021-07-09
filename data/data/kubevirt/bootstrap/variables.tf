variable "storage" {
  type        = string
  description = "bootstrap VM disk size, of type Quantity (see: https://github.com/kubernetes/apimachinery/blob/master/pkg/api/resource/quantity.go)"
  default     = "35Gi"
}

variable "memory" {
  type        = string
  description = "bootstrap VM memory size, of type Quantity (see: https://github.com/kubernetes/apimachinery/blob/master/pkg/api/resource/quantity.go)"
  default     = "8G"
}

variable "cpu" {
  type        = string
  description = "bootstrap VM number of cores"
  default     = "4"
}
