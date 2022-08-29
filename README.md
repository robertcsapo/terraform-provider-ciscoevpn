# Terraform Provider Cisco EVPN
Tech Preview (Early field trial)

terraform-provider-ciscoevpn is a Terraform Provider for Cisco Catalyst 9000 Switches.

## Requirements for Development

- [Terraform](https://www.terraform.io/downloads.html) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.18

Then commit the changes to `go.mod` and `go.sum`.

## Using the provider

Use ```terraform init``` to download the plugin from Terrafrom Registry.

Configure the provider to connect towards your Cisco Catalyst 9000 Switches
```
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
  debug = false
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
```

Examples can be found in [examples/](./examples/).