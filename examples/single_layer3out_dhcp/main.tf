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
  depends_on = [
    ciscoevpn_bgp_system.ibgp,
    ciscoevpn_loopback.spine1_100,
    ciscoevpn_loopback.spine2_100
  ]
  roles  = ["spines"]
  bgp_id = ciscoevpn_bgp_system.ibgp.bgp_id
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
  depends_on = [
    ciscoevpn_bgp_system.ibgp,
    ciscoevpn_loopback.leaf1_100,
    ciscoevpn_loopback.leaf2_100
  ]
  roles  = ["leafs","borders"]
  bgp_id = ciscoevpn_bgp_system.ibgp.bgp_id
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

resource "ciscoevpn_vrf" "green" {
  roles = ["leafs", "borders"]
  name  = "green"
  rd    = "1:1"
  ipv4  = true
  ipv6  = true
}
resource "ciscoevpn_vrf" "blue" {
  roles = ["leafs", "borders"]
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
  roles                  = ["leafs", "borders"]
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
  roles                  = ["leafs,"borders"]
  bgp_id                 = ciscoevpn_bgp_system.ibgp.bgp_id
  vrf                    = ciscoevpn_vrf.blue.name
  ipv4                   = true
  ipv6                   = true
  redistribute_static    = true
  redistribute_connected = true
}

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
  rd                        = "101:101"
  rt                        = "101:101"
  rt_type                   = "both"
  ip_learning               = true
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
resource "ciscoevpn_vlan" "vlan_105" {
  roles   = ["leafs"]
  vlan_id = 105
  vni     = 10105
}

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

resource "ciscoevpn_dhcp" "leafs" {
  depends_on = [
    ciscoevpn_vlan.vlan_101,
    ciscoevpn_vlan.vlan_102,
  ]
  roles = ["leafs"]
  vlans = [
    ciscoevpn_vlan.vlan_101.vlan_id,
    ciscoevpn_vlan.vlan_102.vlan_id,
  ]
  relay_vpn = true
}
resource "ciscoevpn_dhcp_helper" "green" {
  depends_on = [
    ciscoevpn_svi.svi_103
  ]
  roles = ["leafs"]
  svi_id = ciscoevpn_svi.svi_103.svi_id
  ipv4_helper = [
    "100.127.255.2"
  ]
  vrf = "global"
  source_interface = ciscoevpn_loopback.leaf1_200.interface_name
}

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