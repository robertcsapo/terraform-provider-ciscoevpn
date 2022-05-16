variable "username" {
  type        = string
  default     = "developer"
  description = "IOS-XE Username"
}

variable "password" {
  type        = string
  default     = "C1sco12345"
  description = "IOS-XE Passwordd"
}

variable "insecure" {
  type        = bool
  default     = true
  description = "IOS-XE HTTPS insecure"
}

variable "timeout" {
  type        = number
  default     = 120
  description = "IOS-XE Client timeout"
}

variable "iosxe_spines" {
  type        = list(string)
  default     = []
  description = "IOS-XE Devices as Spines"
}

variable "iosxe_leafs" {
  type        = list(string)
  default     = []
  description = "IOS-XE Devices as Leafs"
}

variable "iosxe_borders" {
  type        = list(string)
  default     = []
  description = "IOS-XE Devices as Borders"
}