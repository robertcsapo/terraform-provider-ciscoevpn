
resource "ciscoevpn_evpn" "evpn" {
  depends_on            = [ciscoevpn_loopback.leaf1_200]
  roles                 = ["leafs"]
  replication_type      = "static"
  mac_duplication_limit = 20
  mac_duplication_time  = 10
  ip_duplication_limit  = 20
  ip_duplication_time   = 10
  router_id             = ciscoevpn_loopback.leaf1_200.interface_name
  default_gateway       = "advertise"
  logging_peer_state    = true
  route_target_auto     = "vni"
}

resource "ciscoevpn_evpn_instance" "instance_101" {
  roles                     = ["leafs"]
  instance_id               = 101
  vlan_based                = true
  encapsulation             = "vxlan"
  replication_type          = "ingress"
  rd                        = "65000:101"
  rt                        = "65000:101"
  rt_type                   = "both"
  ip_learning               = false
  default_gateway_advertise = false
  re_originate              = "route-type5"
}
resource "ciscoevpn_evpn_instance" "instance_102" {
  roles                     = ["leafs"]
  instance_id               = 102
  vlan_based                = true
  encapsulation             = "vxlan"
  ip_learning               = false
  default_gateway_advertise = false
}