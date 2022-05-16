package bgp

type CiscoIOSXEBgpNeighbors struct {
	CiscoIOSXEBgpBgp []CiscoIOSXEBgp `json:"Cisco-IOS-XE-bgp:bgp,omitempty"`
}
type CiscoIOSXEBgpNeighborsInterface struct {
	Loopback int `json:"Loopback,omitempty"` // TODO move to common
}
type CiscoIOSXEBgpNeighborsUpdateSource struct {
	Interface CiscoIOSXEBgpNeighborsInterface `json:"interface,omitempty"`
}

type CiscoIOSXEBgpNeighborsNeighbor struct {
	ID           string                             `json:"id,omitempty"`
	RemoteAs     int                                `json:"remote-as,omitempty"`
	UpdateSource CiscoIOSXEBgpNeighborsUpdateSource `json:"update-source,omitempty"`
}
type CiscoIOSXEBgpNeighborsEvpnNeighbor struct {
	ID                   string                              `json:"id,omitempty"`
	Activate             []interface{}                       `json:"activate,omitempty"`
	RouteReflectorClient []interface{}                       `json:"route-reflector-client,omitempty"`
	SendCommunity        CiscoIOSXEBgpNeighborsSendCommunity `json:"send-community,omitempty"`
}

type CiscoIOSXEBgpNeighborsSendCommunity struct {
	SendCommunityWhere string `json:"send-community-where,omitempty"`
}
type CiscoIOSXEBgpNeighborsL2VpnEvpn struct {
	Neighbor []CiscoIOSXEBgpNeighborsEvpnNeighbor `json:"neighbor,omitempty"`
}
type CiscoIOSXEBgpNeighborsL2Vpn struct {
	AfName    string                          `json:"af-name,omitempty"`
	L2VpnEvpn CiscoIOSXEBgpNeighborsL2VpnEvpn `json:"l2vpn-evpn,omitempty"`
}
type CiscoIOSXEBgpNeighborsNoVrf struct {
	L2Vpn []CiscoIOSXEBgpNeighborsL2Vpn `json:"l2vpn,omitempty"`
}
type CiscoIOSXEBgpNeighborsAddressFamily struct {
	NoVrf CiscoIOSXEBgpNeighborsNoVrf `json:"no-vrf,omitempty"`
}

type CiscoIOSXEBgpBgpSystem struct {
	CiscoIOSXEBgpBgp []CiscoIOSXEBgp `json:"Cisco-IOS-XE-bgp:bgp,omitempty"`
}
type CiscoIOSXEBgpBgpSystemDefault struct {
	Ipv4Unicast bool `json:"ipv4-unicast"`
}
type CiscoIOSXEBgpBgpSystemRouterIDInterface struct {
	Loopback int `json:"Loopback,omitempty"`
}
type CiscoIOSXEBgpBgpSystemRouterID struct {
	IP        string                                  `json:"ip-id,omitempty"`
	Interface CiscoIOSXEBgpBgpSystemRouterIDInterface `json:"interface,omitempty"`
}
type CiscoIOSXEBgpBgp struct {
	Default            CiscoIOSXEBgpBgpSystemDefault  `json:"default"`
	LogNeighborChanges bool                           `json:"log-neighbor-changes,omitempty"`
	RouterID           CiscoIOSXEBgpBgpSystemRouterID `json:"router-id,omitempty"`
}

type CiscoIOSXEBgp struct {
	ID            int                                 `json:"id,omitempty"`
	Bgp           CiscoIOSXEBgpBgp                    `json:"bgp,omitempty"`
	Neighbor      []CiscoIOSXEBgpNeighborsNeighbor    `json:"neighbor,omitempty"`
	AddressFamily CiscoIOSXEBgpNeighborsAddressFamily `json:"address-family,omitempty"`
}

type CiscoIOSXEBgpWithVrfs struct {
	CiscoIOSXEBgpWithVrf CiscoIOSXEBgpWithVrf `json:"Cisco-IOS-XE-bgp:with-vrf,omitempty"`
}
type CiscoIOSXEBgpWithVrfL2Vpn struct {
	Evpn []string `json:"evpn,omitempty"`
}
type CiscoIOSXEBgpWithVrfAdvertise struct {
	L2Vpn CiscoIOSXEBgpWithVrfL2Vpn `json:"l2vpn,omitempty"`
}
type CiscoIOSXEBgpWithVrfNeighbor struct {
	ID       string        `json:"id,omitempty"`
	RemoteAs int           `json:"remote-as,omitempty"`
	Activate []interface{} `json:"activate,omitempty"`
}
type CiscoIOSXEBgpWithVrfRedistributeVrf struct {
	Connected interface{} `json:"connected,omitempty"`
	Static    interface{} `json:"static,omitempty"`
}
type CiscoIOSXEBgpWithVrfIpv4Unicast struct {
	Advertise       CiscoIOSXEBgpWithVrfAdvertise       `json:"advertise,omitempty"`
	Neighbor        []CiscoIOSXEBgpWithVrfNeighbor      `json:"neighbor,omitempty"`
	RedistributeVrf CiscoIOSXEBgpWithVrfRedistributeVrf `json:"redistribute-vrf,omitempty"`
}
type CiscoIOSXEBgpWithVrfIpv4 struct {
	AfName string                        `json:"af-name,omitempty"`
	Vrf    []CiscoIOSXEBgpWithVrfVrfIpv4 `json:"vrf,omitempty"`
}
type CiscoIOSXEBgpWithVrfRedistributeV6 struct {
	Connected interface{} `json:"connected,omitempty"`
	Static    interface{} `json:"static,omitempty"`
}
type CiscoIOSXEBgpWithVrfIpv6Unicast struct {
	Advertise      CiscoIOSXEBgpWithVrfAdvertise      `json:"advertise,omitempty"`
	RedistributeV6 CiscoIOSXEBgpWithVrfRedistributeV6 `json:"redistribute-v6,omitempty"`
}
type CiscoIOSXEBgpWithVrfVrfIpv6 struct {
	Name        string                          `json:"name"`
	Ipv6Unicast CiscoIOSXEBgpWithVrfIpv6Unicast `json:"ipv6-unicast,omitempty"`
}
type CiscoIOSXEBgpWithVrfVrfIpv4 struct {
	Name        string                          `json:"name"`
	Ipv4Unicast CiscoIOSXEBgpWithVrfIpv4Unicast `json:"ipv4-unicast,omitempty"`
}
type CiscoIOSXEBgpWithVrfIpv6 struct {
	AfName string                        `json:"af-name,omitempty"`
	Vrf    []CiscoIOSXEBgpWithVrfVrfIpv6 `json:"vrf,omitempty"`
}
type CiscoIOSXEBgpWithVrf struct {
	Ipv4 []CiscoIOSXEBgpWithVrfIpv4 `json:"ipv4,omitempty"`
	Ipv6 []CiscoIOSXEBgpWithVrfIpv6 `json:"ipv6,omitempty"`
}

type CiscoIOSXEBgpVrfIpv4Unicast struct {
	CiscoIOSXEBgpIpv4Unicast CiscoIOSXEBgpIpv4Unicast `json:"Cisco-IOS-XE-bgp:ipv4-unicast"`
}
type CiscoIOSXEBgpIpv4UnicastNeighbor struct {
	ID       string   `json:"id,omitempty"`
	RemoteAs int      `json:"remote-as,omitempty"`
	Activate []string `json:"activate,omitempty"`
}
type CiscoIOSXEBgpIpv4Unicast struct {
	Neighbor []CiscoIOSXEBgpIpv4UnicastNeighbor `json:"neighbor,omitempty"`
}
