# L2VNI
resource "ciscoevpn_svi" "svi_101" {
  depends_on   = [ciscoevpn_vrf.green]
  roles        = ["leafs"]
  svi_id       = 101
  autostate    = false
  vrf          = ciscoevpn_vrf.green.name
  ipv4_address = "100.119.101.1"
  ipv4_mask    = "255.255.255.0"
}
resource "ciscoevpn_svi" "svi_102" {
  depends_on   = [ciscoevpn_vrf.blue]
  roles        = ["leafs"]
  svi_id       = 102
  autostate    = false
  vrf          = ciscoevpn_vrf.blue.name
  ipv4_address = "100.119.102.1"
  ipv4_mask    = "255.255.255.0"
}

# L3VNI
resource "ciscoevpn_svi" "svi_103" {
  depends_on = [ciscoevpn_vrf.green]
  roles      = ["leafs"]
  svi_id     = 103
  autostate  = false
  vrf        = ciscoevpn_vrf.green.name
  unnumbered = ciscoevpn_loopback.leaf1_200.interface_name
}
resource "ciscoevpn_svi" "svi_104" {
  depends_on = [ciscoevpn_vrf.blue]
  roles      = ["leafs"]
  svi_id     = 104
  autostate  = false
  vrf        = ciscoevpn_vrf.blue.name
  unnumbered = ciscoevpn_loopback.leaf1_200.interface_name
}