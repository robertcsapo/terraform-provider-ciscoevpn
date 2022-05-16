package vlan

type CiscoIOSXENativeVlans struct {
	CiscoIOSXENativeVlan CiscoIOSXENativeVlan `json:"Cisco-IOS-XE-native:vlan"`
}

type CiscoIOSXENativeVlan struct {
	CiscoIOSXEVlanConfigurationEntry []CiscoIOSXEVlanConfigurationEntry `json:"Cisco-IOS-XE-vlan:configuration-entry,omitempty"`
	CiscoIOSXEVlanVlanList           []CiscoIOSXEVlanVlanList           `json:"Cisco-IOS-XE-vlan:vlan-list,omitempty"`
}

type CiscoIOSXEVlanConfigurationEntryEvpnInstance struct {
	EvpnInstance int `json:"evpn-instance,omitempty"`
	Vni          int `json:"vni,omitempty"`
}
type CiscoIOSXEVlanConfigurationEntryEvpnMember struct {
	EvpnInstance CiscoIOSXEVlanConfigurationEntryEvpnInstance `json:"evpn-instance,omitempty"`
}
type CiscoIOSXEVlanConfigurationEntryMember struct {
	Vni          int                                           `json:"vni,omitempty"`
	EvpnInstance *CiscoIOSXEVlanConfigurationEntryEvpnInstance `json:"evpn-instance,omitempty"`
}
type CiscoIOSXEVlanConfigurationEntry struct {
	VlanID string                                 `json:"vlan-id,omitempty"`
	Member CiscoIOSXEVlanConfigurationEntryMember `json:"member,omitempty"`
}

type CiscoIOSXEVlanVlanList struct {
	ID   int    `json:"id"`
	Name string `json:"name,omitempty"`
}
