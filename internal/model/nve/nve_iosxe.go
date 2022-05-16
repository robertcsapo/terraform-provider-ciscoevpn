package nve

type CiscoIOSXENativeNves struct {
	CiscoIOSXENativeNve []CiscoIOSXENativeNve `json:"Cisco-IOS-XE-native:nve"`
}
type CiscoIOSXENativeNveProtocol struct {
	//Bgp []interface{} `json:"bgp,omitempty"`
	Bgp []string `json:"bgp,omitempty"`
}
type CiscoIOSXENativeNveHostReachability struct {
	Protocol CiscoIOSXENativeNveProtocol `json:"protocol,omitempty"`
}
type CiscoIOSXENativeNveSourceInterface struct {
	Loopback int `json:"Loopback,omitempty"`
}
type CiscoIOSXENativeNveMcastGroup struct {
	MulticastGroupMin string `json:"multicast-group-min,omitempty"`
}
type CiscoIOSXENativeNveIrCpConfig struct {
	IngressReplication []string `json:"ingress-replication"`
}
type CiscoIOSXENativeNveVni struct {
	VniRange   string                         `json:"vni-range,omitempty"`
	Vrf        string                         `json:"vrf,omitempty"`
	McastGroup *CiscoIOSXENativeNveMcastGroup `json:"mcast-group,omitempty"`
	IrCpConfig *CiscoIOSXENativeNveIrCpConfig `json:"ir-cp-config,omitempty"`
}
type CiscoIOSXENativeNveMember struct {
	Vni []CiscoIOSXENativeNveVni `json:"vni,omitempty"`
}
type MemberInOneLine struct {
	Member CiscoIOSXENativeNveMember `json:"member"`
}
type CiscoIOSXENativeNve struct {
	Name             int                                 `json:"name,omitempty"`
	HostReachability CiscoIOSXENativeNveHostReachability `json:"host-reachability,omitempty"`
	SourceInterface  CiscoIOSXENativeNveSourceInterface  `json:"source-interface,omitempty"`
	MemberInOneLine  MemberInOneLine                     `json:"member-in-one-line"`
	Member           CiscoIOSXENativeNveMember           `json:"member,omitempty"`
	Description      string                              `json:"description,omitempty"`
}
