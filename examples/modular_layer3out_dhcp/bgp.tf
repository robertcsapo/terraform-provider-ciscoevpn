resource "ciscoevpn_bgp_system" "ibgp" {
  depends_on = [
    ciscoevpn_loopback.spine1_100,
    ciscoevpn_loopback.spine2_100,
    ciscoevpn_loopback.leaf1_100,
    ciscoevpn_loopback.leaf2_100,
  ]
  roles                = ["spines", "leafs"]
  router_id            = ciscoevpn_loopback.spine1_100.interface_name
  bgp_id               = 65534
  log_neighbor_changes = true
  default_ipv4_unicast = false
}

resource "ciscoevpn_bgp_neighbor" "spines" {
  depends_on = [ciscoevpn_bgp_system.ibgp, ciscoevpn_loopback.spine1_100]
  roles      = ["spines"]
  bgp_id     = ciscoevpn_bgp_system.ibgp.bgp_id
  neighbors = [
    "${ciscoevpn_loopback.spine1_100.ipv4_address}",
    "${ciscoevpn_loopback.spine2_100.ipv4_address}",
    "${ciscoevpn_loopback.leaf1_100.ipv4_address}",
    "${ciscoevpn_loopback.leaf2_100.ipv4_address}",
  ]
  update_source          = ciscoevpn_loopback.spine1_100.interface_name
  remote_as              = ciscoevpn_bgp_system.ibgp.bgp_id
  send_community         = "both"
  activate               = true
  l2vpn_evpn             = true
  route_reflector_client = true
}

resource "ciscoevpn_bgp_neighbor" "leafs" {
  depends_on = [ciscoevpn_bgp_system.ibgp, ciscoevpn_loopback.leaf1_100]
  roles      = ["leafs"]
  bgp_id     = ciscoevpn_bgp_system.ibgp.bgp_id
  neighbors = [
    "${ciscoevpn_loopback.spine1_100.ipv4_address}",
    "${ciscoevpn_loopback.spine2_100.ipv4_address}",
  ]
  update_source  = ciscoevpn_loopback.leaf1_100.interface_name
  remote_as      = ciscoevpn_bgp_system.ibgp.bgp_id
  send_community = "both"
  activate       = true
  l2vpn_evpn     = true
}