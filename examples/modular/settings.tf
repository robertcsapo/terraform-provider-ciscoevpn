variable "username" {
  type        = string
  default     = "developer"
  description = "Cisco IOS-XE Username"
}

variable "password" {
  type        = string
  default     = "C1sco12345"
  description = "Cisco IOS-XE Passwordd"
}

variable "insecure" {
  type        = bool
  default     = true
  description = "Cisco IOS-XE HTTPS insecure"
}

variable "timeout" {
  type        = number
  default     = 120
  description = "Cisco IOS-XE Client timeout"
}

variable "iosxe_spines" {
  type        = list(string)
  default     = []
  description = "Cisco IOS-XE Devices as Spines"
}

variable "iosxe_leafs" {
  type        = list(string)
  default     = []
  description = "Cisco IOS-XE Devices as Leafs"
}

variable "iosxe_borders" {
  type        = list(string)
  default     = []
  description = "Cisco IOS-XE Devices as Borders"
}