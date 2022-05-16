terraform {
  required_providers {
    ciscoevpn = {
      source = "robertcsapo/ciscoevpn"
      version = "1.0.0"
    }
  }
}

provider "ciscoevpn" {
  username = var.username
  password = var.password
  insecure = var.insecure
  timeout  = var.timeout
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