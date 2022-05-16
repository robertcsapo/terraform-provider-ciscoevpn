package subinterface

type CiscoIOSXENativeEthernet struct {
	FourHundred []CiscoIOSXENativeEthernetInterface `json:"Cisco-IOS-XE-native:FourHundredGigE,omitempty"`
	Hundred     []CiscoIOSXENativeEthernetInterface `json:"Cisco-IOS-XE-native:HundredGigE,omitempty"`
	Forty       []CiscoIOSXENativeEthernetInterface `json:"Cisco-IOS-XE-native:FortyGigabitEthernet,omitempty"`
	TwentyFive  []CiscoIOSXENativeEthernetInterface `json:"Cisco-IOS-XE-native:TwentyFiveGigE,omitempty"`
	Ten         []CiscoIOSXENativeEthernetInterface `json:"Cisco-IOS-XE-native:TenGigabitEthernet,omitempty"`
	One         []CiscoIOSXENativeEthernetInterface `json:"Cisco-IOS-XE-native:GigabitEthernet,omitempty"`
}
type CiscoIOSXENativeEthernetInterfaceDot1Q struct {
	VlanID int `json:"vlan-id,omitempty"`
}
type CiscoIOSXENativeEthernetInterfaceEncapsulation struct {
	Dot1Q CiscoIOSXENativeEthernetInterfaceDot1Q `json:"dot1Q,omitempty"`
}
type CiscoIOSXENativeEthernetInterfaceVrf struct {
	Forwarding string `json:"forwarding,omitempty"`
}
type CiscoIOSXENativeEthernetInterfacePrimary struct {
	Address string `json:"address,omitempty"`
	Mask    string `json:"mask,omitempty"`
}
type CiscoIOSXENativeEthernetInterfaceAddress struct {
	Primary CiscoIOSXENativeEthernetInterfacePrimary `json:"primary,omitempty"`
}
type CiscoIOSXENativeEthernetInterfaceIP struct {
	Address CiscoIOSXENativeEthernetInterfaceAddress `json:"address,omitempty"`
}
type CiscoIOSXENativeEthernetInterface struct {
	Name          string                                         `json:"name,omitempty"`
	Description   string                                         `json:"description,omitempty"`
	Encapsulation CiscoIOSXENativeEthernetInterfaceEncapsulation `json:"encapsulation,omitempty"`
	Vrf           CiscoIOSXENativeEthernetInterfaceVrf           `json:"vrf,omitempty"`
	IP            CiscoIOSXENativeEthernetInterfaceIP            `json:"ip,omitempty"`
}
