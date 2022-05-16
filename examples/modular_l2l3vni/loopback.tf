
# Loopbacks for Spines
resource "ciscoevpn_loopback" "spine1_100" {
  host         = var.iosxe_spines.0
  loopback_id  = 100
  ipv4_address = "100.119.11.1"
  ipv4_mask    = "255.255.255.255"
  pim_sm       = true
}
resource "ciscoevpn_loopback" "spine1_200" {
  host         = var.iosxe_spines.0
  loopback_id  = 200
  ipv4_address = "100.119.12.1"
  ipv4_mask    = "255.255.255.255"
  pim_sm       = true
}
resource "ciscoevpn_loopback" "spine2_100" {
  host         = var.iosxe_spines.1
  loopback_id  = 100
  ipv4_address = "100.119.11.2"
  ipv4_mask    = "255.255.255.255"
  pim_sm       = true
}
resource "ciscoevpn_loopback" "spine2_200" {
  host         = var.iosxe_spines.1
  loopback_id  = 200
  ipv4_address = "100.119.12.2"
  ipv4_mask    = "255.255.255.255"
  pim_sm       = true
}

# Loopbacks for Leafs
resource "ciscoevpn_loopback" "leaf1_100" {
  host         = var.iosxe_leafs.0
  loopback_id  = 100
  ipv4_address = "100.119.11.11"
  ipv4_mask    = "255.255.255.255"
  pim_sm       = true
}
resource "ciscoevpn_loopback" "leaf1_200" {
  host         = var.iosxe_leafs.0
  loopback_id  = 200
  ipv4_address = "100.119.12.11"
  ipv4_mask    = "255.255.255.255"
  pim_sm       = true
}

resource "ciscoevpn_loopback" "leaf2_100" {
  host         = var.iosxe_leafs.1
  loopback_id  = 100
  ipv4_address = "100.119.11.12"
  ipv4_mask    = "255.255.255.255"
  pim_sm       = true
}
resource "ciscoevpn_loopback" "leaf2_200" {
  host         = var.iosxe_leafs.1
  loopback_id  = 200
  ipv4_address = "100.119.12.12"
  ipv4_mask    = "255.255.255.255"
  pim_sm       = true
}