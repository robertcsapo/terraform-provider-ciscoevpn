package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/robertcsapo/terraform-provider-ciscoevpn/internal/model/bgp"
	"github.com/robertcsapo/terraform-provider-ciscoevpn/internal/provider/iosxe"
	"github.com/robertcsapo/terraform-provider-ciscoevpn/internal/provider/service"
)

func resourceCiscoNativeBgpVrf() *schema.Resource {
	return &schema.Resource{
		Description:   "Cisco BGP VRF",
		CreateContext: resourceCiscoNativeBgpVrfCreate,
		ReadContext:   resourceCiscoNativeBgpVrfRead,
		UpdateContext: resourceCiscoNativeBgpVrfUpdate,
		DeleteContext: resourceCiscoNativeBgpVrfDelete,
		Schema: map[string]*schema.Schema{
			"roles": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"bgp_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"vrf": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"ipv4": {
				Type:     schema.TypeBool,
				Default:  true,
				Optional: true,
			},
			"ipv6": {
				Type:     schema.TypeBool,
				Default:  true,
				Optional: true,
			},
			"redistribute_connected": {
				Type:     schema.TypeBool,
				Default:  true,
				Optional: true,
			},
			"redistribute_static": {
				Type:     schema.TypeBool,
				Default:  true,
				Optional: true,
			},
		},
	}
}

func resourceCiscoNativeBgpVrfCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Cisco BGP VRF CREATE")
	var diags diag.Diagnostics
	var err error

	c, _ := meta.(*providerClient)
	svc := &service.Client{
		Method:   "PATCH",
		Path:     fmt.Sprintf("/data/Cisco-IOS-XE-native:native/router/bgp=%v/address-family/with-vrf", d.Get("bgp_id").(int)), // TODO use GET for BGP id
		Provider: c.Provider,
		Devices:  c.Devices.List(),
	}

	roles := d.Get("roles").([]interface{})
	for _, role := range roles {
		svc.Role = role.(string)
		data := c.CiscoIOSXENativeVrfBgp(d)
		if data == nil {
			return diag.Errorf("No data in yang model")
		}

		if b, err := json.MarshalIndent(data, "", "\t"); err == nil {
			svc.Payload = string(b)
		}
		if svc.Provider.Get("debug").(bool) {
			debugJson(fmt.Sprintf("bgp_vrf_%v_%v", svc.Role, d.Get("vrf").(string)), svc.Payload)
		}

		_, err = iosxe.MultiSession(svc)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId(fmt.Sprintf("bgp_vrf_%v", d.Get("vrf").(string)))
	return diags
}

func resourceCiscoNativeBgpVrfRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Println("[DEBUG] TODO")
	return nil
}

func resourceCiscoNativeBgpVrfUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Cisco BGP VRF UPDATE")
	var diags diag.Diagnostics
	var err error

	if d.HasChange("bgp_id") {
		oldState, _ := d.GetChange("bgp_id")
		d.Set("bgp_id", oldState)
		return diag.Errorf("Not supported to change BGP ASN")
	}
	if d.HasChange("vrf") {
		oldState, _ := d.GetChange("vrf")
		d.Set("vrf", oldState)
		return diag.Errorf("Not supported to change VRF name")
	}
	if d.HasChange("roles") {
		oldState, _ := d.GetChange("roles")
		d.Set("roles", oldState)
		return diag.Errorf("Not supported to change Roles")
	}

	c, _ := meta.(*providerClient)
	svc := &service.Client{
		Provider: c.Provider,
		Devices:  c.Devices.List(),
	}

	roles := d.Get("roles").([]interface{})
	for _, role := range roles {
		svc.Role = role.(string)
		if d.HasChange("ipv4") {
			if !d.Get("ipv4").(bool) {
				svc.Method = "DELETE"
				svc.Path = fmt.Sprintf("/data/Cisco-IOS-XE-native:native/router/bgp=%v/address-family/with-vrf/ipv4/unicast/vrf=%v", d.Get("bgp_id").(int), d.Get("vrf").(string))
				_, err = iosxe.MultiSession(svc)
				if err != nil {
					return diag.FromErr(err)
				}
			}
		}
		if d.HasChange("ipv6") {
			if !d.Get("ipv6").(bool) {
				svc.Method = "DELETE"
				svc.Path = fmt.Sprintf("/data/Cisco-IOS-XE-native:native/router/bgp=%v/address-family/with-vrf/ipv6/unicast/vrf=%v", d.Get("bgp_id").(int), d.Get("vrf").(string))
				_, err = iosxe.MultiSession(svc)
				if err != nil {
					return diag.FromErr(err)
				}
			}
		}

		svc.Method = "PATCH"
		svc.Path = fmt.Sprintf("/data/Cisco-IOS-XE-native:native/router/bgp=%v/address-family/with-vrf", d.Get("bgp_id").(int))
		data := c.CiscoIOSXENativeVrfBgp(d)
		if data == nil {
			return diag.Errorf("No data in yang model")
		}

		if b, err := json.MarshalIndent(data, "", "\t"); err == nil {
			svc.Payload = string(b)
		}
		if svc.Provider.Get("debug").(bool) {
			debugJson(fmt.Sprintf("bgp_vrf_%v_%v", svc.Role, d.Get("vrf").(string)), svc.Payload)
		}

		_, err = iosxe.MultiSession(svc)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId(fmt.Sprintf("bgp_vrf_%v", d.Get("vrf").(string)))
	return diags
}

func resourceCiscoNativeBgpVrfDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Cisco BGP VRF DELETE")
	var diags diag.Diagnostics
	var err error

	c, _ := meta.(*providerClient)
	svc := &service.Client{
		Method:   "DELETE",
		Provider: c.Provider,
		Devices:  c.Devices.List(),
	}

	roles := d.Get("roles").([]interface{})
	for _, role := range roles {
		svc.Role = role.(string)
		if d.Get("ipv4").(bool) {
			svc.Path = fmt.Sprintf("/data/Cisco-IOS-XE-native:native/router/bgp=%v/address-family/with-vrf/ipv4/unicast/vrf=%v", d.Get("bgp_id").(int), d.Get("vrf").(string))
			_, err = iosxe.MultiSession(svc)
			if err != nil {
				return diag.FromErr(err)
			}
		}
		if d.Get("ipv6").(bool) {
			svc.Path = fmt.Sprintf("/data/Cisco-IOS-XE-native:native/router/bgp=%v/address-family/with-vrf/ipv6/unicast/vrf=%v", d.Get("bgp_id").(int), d.Get("vrf").(string))
			_, err = iosxe.MultiSession(svc)
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	d.SetId("")
	return diags
}

func (*providerClient) CiscoIOSXENativeVrfBgp(d *schema.ResourceData) *bgp.CiscoIOSXEBgpWithVrfs {
	data := &bgp.CiscoIOSXEBgpWithVrf{}

	if d.Get("ipv4").(bool) {
		ipv4 := &bgp.CiscoIOSXEBgpWithVrfIpv4{}
		ipv4.AfName = "unicast"
		ipv4Vrf := &bgp.CiscoIOSXEBgpWithVrfVrfIpv4{}
		ipv4Vrf.Name = d.Get("vrf").(string)
		ipv4Vrf.Ipv4Unicast.Advertise.L2Vpn.Evpn = null()
		if d.Get("redistribute_static").(bool) {
			ipv4Vrf.Ipv4Unicast.RedistributeVrf.Static = map[string]string{}
		}
		if d.Get("redistribute_connected").(bool) {
			ipv4Vrf.Ipv4Unicast.RedistributeVrf.Connected = map[string]string{}
		}
		ipv4.Vrf = append(ipv4.Vrf, *ipv4Vrf)
		data.Ipv4 = append(data.Ipv4, *ipv4)
	}
	if d.Get("ipv6").(bool) {
		ipv6 := &bgp.CiscoIOSXEBgpWithVrfIpv6{}
		ipv6.AfName = "unicast"
		ipv6Vrf := &bgp.CiscoIOSXEBgpWithVrfVrfIpv6{}
		ipv6Vrf.Name = d.Get("vrf").(string)
		ipv6Vrf.Ipv6Unicast.Advertise.L2Vpn.Evpn = null()
		if d.Get("redistribute_static").(bool) {
			ipv6Vrf.Ipv6Unicast.RedistributeV6.Static = map[string]string{}
		}
		if d.Get("redistribute_connected").(bool) {
			ipv6Vrf.Ipv6Unicast.RedistributeV6.Connected = map[string]string{}
		}
		ipv6.Vrf = append(ipv6.Vrf, *ipv6Vrf)
		data.Ipv6 = append(data.Ipv6, *ipv6)
	}
	return &bgp.CiscoIOSXEBgpWithVrfs{*data}
}
