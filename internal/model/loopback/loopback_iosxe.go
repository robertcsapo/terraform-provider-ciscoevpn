package loopback

type CiscoIOSXENativeLoopbackInterface struct {
	CiscoIOSXENativeLoopback []CiscoIOSXENativeLoopback `json:"Cisco-IOS-XE-native:Loopback"`
}
type CiscoIOSXENativeLoopbackPrimary struct {
	Address string `json:"address,omitempty"`
	Mask    string `json:"mask,omitempty"`
}
type CiscoIOSXENativeLoopbackAddress struct {
	Primary CiscoIOSXENativeLoopbackPrimary `json:"primary,omitempty"`
}
type CiscoIOSXEMulticastPimModeChoiceCfg struct {
	SparseMode interface{} `json:"sparse-mode,omitempty"`
}
type CiscoIOSXENativeLoopbackPim struct {
	CiscoIOSXEMulticastPimModeChoiceCfg CiscoIOSXEMulticastPimModeChoiceCfg `json:"Cisco-IOS-XE-multicast:pim-mode-choice-cfg,omitempty"`
}
type CiscoIOSXENativeLoopbackIP struct {
	Address CiscoIOSXENativeLoopbackAddress `json:"address,omitempty"`
	Pim     CiscoIOSXENativeLoopbackPim     `json:"pim,omitempty"`
}
type CiscoIOSXENativeLoopback struct {
	Name        int                        `json:"name"`
	Description string                     `json:"description,omitempty"`
	IP          CiscoIOSXENativeLoopbackIP `json:"ip,omitempty"`
}
