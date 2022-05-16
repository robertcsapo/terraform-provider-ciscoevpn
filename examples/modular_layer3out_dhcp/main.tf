terraform {
  required_providers {
    ciscoevpn = {
      source = "terraform.local/local/ciscoevpn"
      //version = "0.0.1"
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