
resource "ciscoevpn_vrf" "green" {
  roles = ["spines", "leafs"]
  name  = "green"
  rd    = "1:1"
  ipv4  = true
  ipv6  = true
}
resource "ciscoevpn_vrf" "blue" {
  roles = ["spines", "leafs"]
  name  = "blue"
  rd    = "2:2"
  ipv4  = true
  ipv6  = true
}

resource "ciscoevpn_bgp_vrf" "green" {
  depends_on = [
    ciscoevpn_bgp_system.ibgp,
    ciscoevpn_bgp_neighbor.spines,
    ciscoevpn_bgp_neighbor.leafs,
    ciscoevpn_vrf.green
  ]
  roles                  = ["spines", "leafs"]
  bgp_id                 = ciscoevpn_bgp_system.ibgp.bgp_id
  vrf                    = ciscoevpn_vrf.green.name
  ipv4                   = true
  ipv6                   = true
  redistribute_static    = true
  redistribute_connected = true
}
resource "ciscoevpn_bgp_vrf" "blue" {
  depends_on = [
    ciscoevpn_bgp_system.ibgp,
    ciscoevpn_bgp_neighbor.spines,
    ciscoevpn_bgp_neighbor.leafs,
    ciscoevpn_vrf.blue
  ]
  roles                  = ["spines", "leafs"]
  bgp_id                 = ciscoevpn_bgp_system.ibgp.bgp_id
  vrf                    = ciscoevpn_vrf.blue.name
  ipv4                   = true
  ipv6                   = true
  redistribute_static    = true
  redistribute_connected = true
}
