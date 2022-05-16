resource "ciscoevpn_subinterface" "leaf1_green" {
  depends_on = [ 
    ciscoevpn_vrf.green
  ]
  host         = var.iosxe_borders.0
  ethernet = "1/1/1.253"
  interface_speed = 10
  dot1q = 253
  vrf = ciscoevpn_vrf.green.name
  ipv4_address = "100.119.253.10"
  ipv4_mask = "255.255.255.252"
  ipv4_remote = "100.119.253.9"
}
resource "ciscoevpn_subinterface" "leaf1_blue" {
  depends_on = [ 
    ciscoevpn_vrf.blue
  ]
  host         = var.iosxe_borders.0
  ethernet = "1/1/1.254"
  interface_speed = 10
  dot1q = 254
  vrf = ciscoevpn_vrf.blue.name
  ipv4_address = "100.119.254.10"
  ipv4_mask = "255.255.255.252"
  ipv4_remote = "100.119.254.9"
}

resource "ciscoevpn_subinterface" "leaf2_green" {
  depends_on = [ 
    ciscoevpn_vrf.green
  ]
  host         = var.iosxe_borders.1
  ethernet = "1/1/1.253"
  interface_speed = 10
  dot1q = 253
  vrf = ciscoevpn_vrf.green.name
  ipv4_address = "100.119.253.14"
  ipv4_mask = "255.255.255.252"
  ipv4_remote = "100.119.253.13"
}
resource "ciscoevpn_subinterface" "leaf2_blue" {
  depends_on = [ 
    ciscoevpn_vrf.blue
  ]
  host         = var.iosxe_borders.1
  ethernet = "1/1/1.254"
  interface_speed = 10
  dot1q = 254
  vrf = ciscoevpn_vrf.blue.name
  ipv4_address = "100.119.254.14"
  ipv4_mask = "255.255.255.252"
  ipv4_remote = "100.119.254.13"
}