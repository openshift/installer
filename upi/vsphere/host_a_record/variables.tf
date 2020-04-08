variable "zone_id" {
  type        = string
  description = "The ID of the hosted zone to contain this record."
}

variable "records" {
  type        = map(string)
  description = "A records to be added to the zone_id"
}
