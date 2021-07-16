variable "storage" {
  type        = string
  description = "persistant data volume disk size, of type Quantity (see: https://github.com/kubernetes/apimachinery/blob/master/pkg/api/resource/quantity.go)"
  default     = "20Gi"
}
