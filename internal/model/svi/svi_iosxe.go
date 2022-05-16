package svi

type CiscoIOSXENativeSvis struct {
	CiscoIOSXENativeVlan []CiscoIOSXENativeSvi `json:"Cisco-IOS-XE-native:Vlan"`
}
type CiscoIOSXENativeVlanVrf struct {
	Forwarding string `json:"forwarding,omitempty"`
}
type CiscoIOSXENativeVlanPrimary struct {
	Address string `json:"address,omitempty"`
	Mask    string `json:"mask,omitempty"`
}
type CiscoIOSXENativeVlanAddress struct {
	Primary CiscoIOSXENativeVlanPrimary `json:"primary,omitempty"`
}
type CiscoIOSXENativeVlanIP struct {
	Address    *CiscoIOSXENativeVlanAddress `json:"address,omitempty"`
	Unnumbered string                       `json:"unnumbered,omitempty"`
}
type CiscoIOSXENativeSvi struct {
	Name        int                      `json:"name"`
	AutoState   bool                     `json:"autostate"`
	Description string                   `json:"description,omitempty"`
	Vrf         *CiscoIOSXENativeVlanVrf `json:"vrf,omitempty"`
	IP          CiscoIOSXENativeVlanIP   `json:"ip,omitempty"`
}

type CiscoIOSXENativeVlanDhcpHelper struct {
	CiscoIOSXENativeVlan []CiscoIOSXENativeVlanDhcpHelperSvi `json:"Cisco-IOS-XE-native:Vlan"`
}
type CiscoIOSXENativeVlanDhcpHelperSviHelperAddress struct {
	Address string   `json:"address,omitempty"`
	Global  []string `json:"global,omitempty"`
}
type CiscoIOSXENativeVlanDhcpHelperSviCiscoIOSXEDhcpRelay struct {
	SourceInterface string `json:"source-interface,omitempty"`
}
type CiscoIOSXENativeVlanDhcpHelperSviDhcp struct {
	CiscoIOSXEDhcpRelay CiscoIOSXENativeVlanDhcpHelperSviCiscoIOSXEDhcpRelay `json:"Cisco-IOS-XE-dhcp:relay,omitempty"`
}
type CiscoIOSXENativeVlanDhcpHelperSviIP struct {
	HelperAddress []CiscoIOSXENativeVlanDhcpHelperSviHelperAddress `json:"helper-address,omitempty"`
	Dhcp          CiscoIOSXENativeVlanDhcpHelperSviDhcp            `json:"dhcp,omitempty"`
}
type CiscoIOSXENativeVlanDhcpHelperSvi struct {
	Name int                                 `json:"name,omitempty"`
	IP   CiscoIOSXENativeVlanDhcpHelperSviIP `json:"ip,omitempty"`
}
