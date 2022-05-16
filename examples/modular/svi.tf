resource "ciscoevpn_svi" "svi_101" {
  roles        = ["leafs"]
  svi_id       = 101
  autostate    = false
  vrf          = ciscoevpn_vrf.green.name
  ipv4_address = "100.119.101.1"
  ipv4_mask    = "255.255.255.0"
}
resource "ciscoevpn_svi" "svi_102" {
  roles        = ["leafs"]
  svi_id       = 102
  autostate    = true
  vrf          = ciscoevpn_vrf.blue.name
  ipv4_address = "100.119.102.1"
  ipv4_mask    = "255.255.255.0"
}