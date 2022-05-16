
resource "ciscoevpn_vlan" "vlan_101" {
  depends_on    = [ciscoevpn_evpn_instance.instance_101]
  roles         = ["leafs"]
  vlan_id       = 101
  evpn_instance = ciscoevpn_evpn_instance.instance_101.instance_id
  vni           = 10101
}
resource "ciscoevpn_vlan" "vlan_102" {
  depends_on    = [ciscoevpn_evpn_instance.instance_102]
  roles         = ["leafs"]
  vlan_id       = 102
  evpn_instance = ciscoevpn_evpn_instance.instance_102.instance_id
  vni           = 10102
}
resource "ciscoevpn_vlan" "vlan_103" {
  roles   = ["leafs"]
  vlan_id = 103
  vni     = 10103
}
resource "ciscoevpn_vlan" "vlan_104" {
  roles   = ["leafs"]
  vlan_id = 104
  vni     = 10104
}
