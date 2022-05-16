package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/robertcsapo/terraform-provider-ciscoevpn/internal/model/dhcp"
	"github.com/robertcsapo/terraform-provider-ciscoevpn/internal/provider/iosxe"
	"github.com/robertcsapo/terraform-provider-ciscoevpn/internal/provider/service"
)

func resourceCiscoNativeDhcp() *schema.Resource {
	return &schema.Resource{
		Description:   "Cisco DHCP Global",
		CreateContext: resourceCiscoNativeDhcpCreate,
		ReadContext:   resourceCiscoNativeDhcpRead,
		UpdateContext: resourceCiscoNativeDhcpUpdate,
		DeleteContext: resourceCiscoNativeDhcpDelete,
		Schema: map[string]*schema.Schema{
			"roles": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"vlans": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
			"relay_vpn": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
		},
	}
}

func resourceCiscoNativeDhcpCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Cisco DHCP CREATE")
	var diags diag.Diagnostics
	var err error

	c, _ := meta.(*providerClient)
	svc := &service.Client{
		Method:   "PATCH",
		Path:     "/data/Cisco-IOS-XE-native:native/ip/dhcp",
		Provider: c.Provider,
		Devices:  c.Devices.List(),
	}

	roles := d.Get("roles").([]interface{})
	for _, role := range roles {
		svc.Role = role.(string)
		data := c.resourceCiscoNativeDhcpData(d)
		if data == nil {
			return diag.Errorf("No data in yang model")
		}

		if b, err := json.MarshalIndent(data, "", "\t"); err == nil {
			svc.Payload = string(b)
		}
		if svc.Provider.Get("debug").(bool) {
			debugJson(fmt.Sprintf("dhcp_%v", svc.Role), svc.Payload)
		}

		_, err = iosxe.MultiSession(svc)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId(fmt.Sprintf("dhcp_global_%v", svc.Role))
	return diags
}

func resourceCiscoNativeDhcpRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Println("[DEBUG] TODO")
	var diags diag.Diagnostics
	return diags
}

func resourceCiscoNativeDhcpUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Cisco DHCP UPDATE")
	var diags diag.Diagnostics
	var err error

	if d.HasChange("roles") {
		oldState, _ := d.GetChange("roles")
		d.Set("roles", oldState)
		return diag.Errorf("Not supported to change Roles")
	}
	c, _ := meta.(*providerClient)
	svc := &service.Client{
		Method:   "PATCH",
		Path:     "/data/Cisco-IOS-XE-native:native/ip/dhcp",
		Provider: c.Provider,
		Devices:  c.Devices.List(),
	}

	roles := d.Get("roles").([]interface{})
	for _, role := range roles {
		svc.Role = role.(string)
		data := c.resourceCiscoNativeDhcpData(d)
		if data == nil {
			return diag.Errorf("No data in yang model")
		}

		if b, err := json.MarshalIndent(data, "", "\t"); err == nil {
			svc.Payload = string(b)
		}
		if svc.Provider.Get("debug").(bool) {
			debugJson(fmt.Sprintf("dhcp_%v", svc.Role), svc.Payload)
		}

		_, err = iosxe.MultiSession(svc)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId(fmt.Sprintf("dhcp_global_%v", svc.Role))
	return diags
}

func resourceCiscoNativeDhcpDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Cisco DHCP DELETE")
	var err error
	var diags diag.Diagnostics

	c, _ := meta.(*providerClient)
	svc := &service.Client{
		Method:   "DELETE",
		Path:     "/data/Cisco-IOS-XE-native:native/ip/dhcp",
		Provider: c.Provider,
		Devices:  c.Devices.List(),
	}

	roles := d.Get("roles").([]interface{})
	for _, role := range roles {
		svc.Role = role.(string)
		_, err = iosxe.MultiSession(svc)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId("")
	return diags
}

func (*providerClient) resourceCiscoNativeDhcpData(d *schema.ResourceData) *dhcp.CiscoIOSXENativeDhcps {
	var n interface{}
	data := &dhcp.CiscoIOSXENativeDhcps{}
	dhcpData := &dhcp.CiscoIOSXENativeDhcp{}

	dhcpData.CiscoIOSXEDhcpCompatibility.Suboption.LinkSelection = "standard"
	dhcpData.CiscoIOSXEDhcpCompatibility.Suboption.ServerOverride = "standard"
	dhcpData.CiscoIOSXEDhcpRelay.Information.Option.OptionDefault = append(dhcpData.CiscoIOSXEDhcpRelay.Information.Option.OptionDefault, n)
	dhcpData.CiscoIOSXEDhcpSnooping = append(dhcpData.CiscoIOSXEDhcpSnooping, n)
	if d.Get("relay_vpn").(bool) {
		dhcpData.CiscoIOSXEDhcpRelay.Information.Option.Vpn = append(dhcpData.CiscoIOSXEDhcpRelay.Information.Option.Vpn, n)
	}

	var vlans string
	v := d.Get("vlans").([]interface{})
	for _, vlan := range v {
		vlans = fmt.Sprintf("%v,%v", vlans, vlan.(int))
	}
	vlans = vlans[1:]
	vlanList := &dhcp.CiscoIOSXEDhcpSnoopingConfSnoopingVlanList{
		ID: vlans,
	}
	dhcpData.CiscoIOSXEDhcpSnoopingConf.Snooping.VlanList = append(dhcpData.CiscoIOSXEDhcpSnoopingConf.Snooping.VlanList, *vlanList)
	data.CiscoIOSXENativeDhcp = *dhcpData
	return data

}
