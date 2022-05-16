terraform {
  required_providers {
    ciscoevpn = {
      source = "robertcsapo/ciscoevpn"
      version = "1.0.0"
    }
  }
}
# ciscoevpn provider settings
provider "ciscoevpn" {
  username = var.username
  password = var.password
  insecure = var.insecure
  timeout  = var.timeout
  debug    = true
  roles {
    spines {
      iosxe = var.iosxe_spines
    }
    borders {
      iosxe = var.iosxe_borders
    }
    leafs {
      iosxe = var.iosxe_leafs
    }
  }
}