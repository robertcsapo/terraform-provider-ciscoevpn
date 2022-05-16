
# Network Virtual Interface
resource "ciscoevpn_nve" "leafs" {
  depends_on = [
    ciscoevpn_vlan.vlan_101,
    ciscoevpn_vlan.vlan_102,
    ciscoevpn_vlan.vlan_103,
    ciscoevpn_vlan.vlan_104,
  ]
  roles            = ["leafs"]
  source_interface = ciscoevpn_loopback.leaf1_200.interface_name
  vni = {
    "${ciscoevpn_vrf.green.name}" = "${ciscoevpn_vlan.vlan_103.vni}"
    "${ciscoevpn_vrf.blue.name}"  = "${ciscoevpn_vlan.vlan_104.vni}"

  }
  vni_ipv4_multicast_group = {
    "225.0.0.101" = "${ciscoevpn_vlan.vlan_101.vni}"
  }
  vni_ingress_replication = ["${ciscoevpn_vlan.vlan_102.vni}"]
}