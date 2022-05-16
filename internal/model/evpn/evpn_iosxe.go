package evpn

type CiscoIOSXEL2Evpn struct {
	CiscoIOSXEL2VpnEvpn CiscoIOSXEL2EvpnEvpn `json:"Cisco-IOS-XE-l2vpn:evpn"`
}
type CiscoIOSXEL2VpnEvpnReplicationType struct {
	Static []int `json:"static"`
}
type CiscoIOSXEL2VpnEvpnDuplication struct {
	Limit int `json:"limit"`
	Time  int `json:"time"`
}
type CiscoIOSXEL2VpnEvpnMac struct {
	Duplication CiscoIOSXEL2VpnEvpnDuplication `json:"duplication"`
}
type CiscoIOSXEL2VpnEvpnIP struct {
	Duplication CiscoIOSXEL2VpnEvpnDuplication `json:"duplication"`
}
type CiscoIOSXEL2VpnEvpnInterface struct {
	Loopback int `json:"Loopback"`
}
type CiscoIOSXEL2VpnEvpnRouterID struct {
	Interface CiscoIOSXEL2VpnEvpnInterface `json:"interface"`
}
type CiscoIOSXEL2VpnEvpnDefaultGateway struct {
	Advertise []int `json:"advertise"`
}
type CiscoIOSXEL2VpnEvpnPeer struct {
	State []int `json:"state"`
}
type CiscoIOSXEL2VpnEvpnLogging struct {
	Peer CiscoIOSXEL2VpnEvpnPeer `json:"peer"`
}
type CiscoIOSXEL2VpnEvpnAuto struct {
	Vni []int `json:"vni"`
}
type CiscoIOSXEL2VpnEvpnRouteTarget struct {
	Auto CiscoIOSXEL2VpnEvpnAuto `json:"auto"`
}
type CiscoIOSXEL2EvpnEvpn struct {
	ReplicationType CiscoIOSXEL2VpnEvpnReplicationType `json:"replication-type"`
	Mac             CiscoIOSXEL2VpnEvpnMac             `json:"mac"`
	IP              CiscoIOSXEL2VpnEvpnIP              `json:"ip"`
	RouterID        CiscoIOSXEL2VpnEvpnRouterID        `json:"router-id"`
	DefaultGateway  CiscoIOSXEL2VpnEvpnDefaultGateway  `json:"default-gateway"`
	Logging         CiscoIOSXEL2VpnEvpnLogging         `json:"logging"`
	RouteTarget     CiscoIOSXEL2VpnEvpnRouteTarget     `json:"route-target"`
}
