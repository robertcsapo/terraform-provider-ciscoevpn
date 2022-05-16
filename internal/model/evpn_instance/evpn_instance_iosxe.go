package evpn_instance

type CiscoIOSXEL2VpnInstanceCiscoIOSXEL2VpnInstanceEvpn struct {
	CiscoIOSXEL2VpnInstance CiscoIOSXEL2VpnInstance `json:"Cisco-IOS-XE-l2vpn:instance"`
}
type CiscoIOSXEL2VpnInstanceReplicationType struct {
	Static  []string `json:"static,omitempty"`
	Ingress []string `json:"ingress,omitempty"`
}
type CiscoIOSXEL2VpnInstanceRd struct {
	RdValue string `json:"rd-value,omitempty"`
}
type CiscoIOSXEL2VpnInstanceBoth struct {
	RtValue string `json:"rt-value,omitempty"`
}
type CiscoIOSXEL2VpnInstanceRouteTarget struct {
	Both CiscoIOSXEL2VpnInstanceBoth `json:"both,omitempty"`
}
type CiscoIOSXEL2VpnInstanceLocalLearning struct {
	Disable []string `json:"disable,omitempty"`
}
type CiscoIOSXEL2VpnInstanceIP struct {
	LocalLearning CiscoIOSXEL2VpnInstanceLocalLearning `json:"local-learning,omitempty"`
}
type CiscoIOSXEL2VpnInstanceDefaultGateway struct {
	Advertise string `json:"advertise,omitempty"`
}
type CiscoIOSXEL2VpnInstanceReOriginate struct {
	RouteType5 []string `json:"route-type5,omitempty"`
}
type CiscoIOSXEL2VpnInstanceVlanBased struct {
	ReplicationType CiscoIOSXEL2VpnInstanceReplicationType `json:"replication-type,omitempty"`
	Encapsulation   string                                 `json:"encapsulation,omitempty"`
	Rd              CiscoIOSXEL2VpnInstanceRd              `json:"rd,omitempty"`
	RouteTarget     CiscoIOSXEL2VpnInstanceRouteTarget     `json:"route-target,omitempty"`
	IP              CiscoIOSXEL2VpnInstanceIP              `json:"ip,omitempty"`
	DefaultGateway  CiscoIOSXEL2VpnInstanceDefaultGateway  `json:"default-gateway,omitempty"`
	ReOriginate     CiscoIOSXEL2VpnInstanceReOriginate     `json:"re-originate,omitempty"`
}
type CiscoIOSXEL2VpnInstanceInstance struct {
	EvpnInstanceNum int                              `json:"evpn-instance-num,omitempty"`
	VlanBased       CiscoIOSXEL2VpnInstanceVlanBased `json:"vlan-based,omitempty"`
}
type CiscoIOSXEL2VpnInstance struct {
	Instance []CiscoIOSXEL2VpnInstanceInstance `json:"instance,omitempty"`
}
