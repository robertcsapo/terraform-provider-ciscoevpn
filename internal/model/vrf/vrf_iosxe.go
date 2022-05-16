package vrf

type CiscoIOSXENativeVrf struct {
	CiscoIOSXENativeDefinition []CiscoIOSXENativeDefinition `json:"Cisco-IOS-XE-native:definition"`
}
type CiscoIOSXENativeDefinitionWithoutStitching struct {
	AsnIP string `json:"asn-ip,omitempty"`
}
type CiscoIOSXENativeDefinitionWithStitching struct {
	AsnIP     string   `json:"asn-ip,omitempty"`
	Stitching []string `json:"stitching,omitempty"`
}
type CiscoIOSXENativeDefinitionExportRouteTarget struct {
	WithoutStitching []CiscoIOSXENativeDefinitionWithoutStitching `json:"without-stitching,omitempty"`
	WithStitching    []CiscoIOSXENativeDefinitionWithStitching    `json:"with-stitching,omitempty"`
}
type CiscoIOSXENativeDefinitionImportRouteTarget struct {
	WithoutStitching []CiscoIOSXENativeDefinitionWithoutStitching `json:"without-stitching,omitempty"`
	WithStitching    []CiscoIOSXENativeDefinitionWithStitching    `json:"with-stitching,omitempty"`
}
type CiscoIOSXENativeDefinitionRouteTarget struct {
	ExportRouteTarget CiscoIOSXENativeDefinitionExportRouteTarget `json:"export-route-target,omitempty"`
	ImportRouteTarget CiscoIOSXENativeDefinitionImportRouteTarget `json:"import-route-target,omitempty"`
}
type CiscoIOSXENativeDefinitionIpv4 struct {
	RouteTarget CiscoIOSXENativeDefinitionRouteTarget `json:"route-target,omitempty"`
}
type CiscoIOSXENativeDefinitionIpv6 struct {
	RouteTarget CiscoIOSXENativeDefinitionRouteTarget `json:"route-target,omitempty"`
}
type CiscoIOSXENativeDefinitionAddressFamily struct {
	Ipv4 CiscoIOSXENativeDefinitionIpv4 `json:"ipv4,omitempty"`
	Ipv6 CiscoIOSXENativeDefinitionIpv6 `json:"ipv6,omitempty"`
}
type CiscoIOSXENativeDefinition struct {
	Name          string                                  `json:"name"`
	Rd            string                                  `json:"rd"`
	AddressFamily CiscoIOSXENativeDefinitionAddressFamily `json:"address-family,omitempty"`
}
