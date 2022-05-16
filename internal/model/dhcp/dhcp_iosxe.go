package dhcp

type CiscoIOSXENativeDhcps struct {
	CiscoIOSXENativeDhcp CiscoIOSXENativeDhcp `json:"Cisco-IOS-XE-native:dhcp"`
}
type CiscoIOSXENativeDhcpSuboption struct {
	LinkSelection  string `json:"link-selection,omitempty"`
	ServerOverride string `json:"server-override,omitempty"`
}
type CiscoIOSXENativeDhcpCompatibility struct {
	Suboption CiscoIOSXENativeDhcpSuboption `json:"suboption,omitempty"`
}
type CiscoIOSXEDhcpRelayInformationOption struct {
	OptionDefault []interface{} `json:"option-default,omitempty"`
	Vpn           []interface{} `json:"vpn,omitempty"`
}
type CiscoIOSXEDhcpRelayInformation struct {
	Option CiscoIOSXEDhcpRelayInformationOption `json:"option,omitempty"`
}
type CiscoIOSXEDhcpRelay struct {
	Information CiscoIOSXEDhcpRelayInformation `json:"information,omitempty"`
}
type CiscoIOSXEDhcpSnoopingConfSnoopingVlanList struct {
	ID string `json:"id,omitempty"`
}
type CiscoIOSXEDhcpSnoopingConfSnooping struct {
	VlanList []CiscoIOSXEDhcpSnoopingConfSnoopingVlanList `json:"vlan-list,omitempty"`
}
type CiscoIOSXEDhcpSnoopingConf struct {
	Snooping CiscoIOSXEDhcpSnoopingConfSnooping `json:"snooping,omitempty"`
}
type CiscoIOSXENativeDhcp struct {
	CiscoIOSXEDhcpCompatibility CiscoIOSXENativeDhcpCompatibility `json:"Cisco-IOS-XE-dhcp:compatibility"`
	CiscoIOSXEDhcpRelay         CiscoIOSXEDhcpRelay               `json:"Cisco-IOS-XE-dhcp:relay,omitempty"`
	CiscoIOSXEDhcpSnooping      []interface{}                     `json:"Cisco-IOS-XE-dhcp:snooping,omitempty"`
	CiscoIOSXEDhcpSnoopingConf  CiscoIOSXEDhcpSnoopingConf        `json:"Cisco-IOS-XE-dhcp:snooping-conf,omitempty"`
}
