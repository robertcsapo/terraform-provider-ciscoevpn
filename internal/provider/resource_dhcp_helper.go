package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/robertcsapo/terraform-provider-ciscoevpn/internal/model/svi"
	"github.com/robertcsapo/terraform-provider-ciscoevpn/internal/provider/iosxe"
	"github.com/robertcsapo/terraform-provider-ciscoevpn/internal/provider/service"
)

func resourceCiscoNativeDhcpHelper() *schema.Resource {
	return &schema.Resource{
		Description:   "Cisco DHCP IP Helper",
		CreateContext: resourceCiscoNativeDhcpHelperCreate,
		ReadContext:   resourceCiscoNativeDhcpHelperRead,
		UpdateContext: resourceCiscoNativeDhcpHelperUpdate,
		DeleteContext: resourceCiscoNativeDhcpHelperDelete,
		Schema: map[string]*schema.Schema{
			"roles": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"svi_id": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(2, 4094),
			},
			"ipv4_helper": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.IsIPv4Address,
				},
			},
			"vrf": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "global",
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"source_interface": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceCiscoNativeDhcpHelperCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Cisco DHCP IP HELPER CREATE")
	var diags diag.Diagnostics
	var err error

	c, _ := meta.(*providerClient)
	svc := &service.Client{
		Method:   "PATCH",
		Path:     fmt.Sprintf("/data/Cisco-IOS-XE-native:native/interface/Vlan=%v", d.Get("svi_id").(int)),
		Provider: c.Provider,
		Devices:  c.Devices.List(),
	}

	roles := d.Get("roles").([]interface{})
	for _, role := range roles {
		svc.Role = role.(string)
		if _, ok := d.GetOk("ipv4_helper"); ok {
			data := c.resourceCiscoNativeDhcpHelperCreateIpv4Data(d)
			if data == nil {
				return diag.Errorf("No data in yang model")
			}

			if b, err := json.MarshalIndent(data, "", "\t"); err == nil {
				svc.Payload = string(b)
			}
			if svc.Provider.Get("debug").(bool) {
				debugJson(fmt.Sprintf("dhcp_helper_%v", svc.Role), svc.Payload)
			}

			_, err = iosxe.MultiSession(svc)
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	d.SetId(fmt.Sprintf("dhcp_helper_%v", svc.Role))
	return diags
}

func resourceCiscoNativeDhcpHelperRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Println("[DEBUG] TODO")
	var diags diag.Diagnostics
	return diags
}

func resourceCiscoNativeDhcpHelperUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Cisco DHCP IP HELPER UPDATE")
	var diags diag.Diagnostics
	var err error

	if d.HasChange("svi") {
		oldState, _ := d.GetChange("svi")
		d.Set("svi", oldState)
		return diag.Errorf("Not supported to change SVI")
	}
	if d.HasChange("roles") {
		oldState, _ := d.GetChange("roles")
		d.Set("roles", oldState)
		return diag.Errorf("Not supported to change Roles")
	}
	c, _ := meta.(*providerClient)
	svc := &service.Client{
		Method:   "PATCH",
		Path:     fmt.Sprintf("/data/Cisco-IOS-XE-native:native/interface/Vlan=%v", d.Get("svi_id").(int)),
		Provider: c.Provider,
		Devices:  c.Devices.List(),
	}

	roles := d.Get("roles").([]interface{})
	for _, role := range roles {
		svc.Role = role.(string)
		if _, ok := d.GetOk("ipv4_helper"); ok {
			data := c.resourceCiscoNativeDhcpHelperCreateIpv4Data(d)
			if data == nil {
				return diag.Errorf("No data in yang model")
			}

			if b, err := json.MarshalIndent(data, "", "\t"); err == nil {
				svc.Payload = string(b)
			}
			if svc.Provider.Get("debug").(bool) {
				debugJson(fmt.Sprintf("dhcp_helper_%v", svc.Role), svc.Payload)
			}

			_, err = iosxe.MultiSession(svc)
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	d.SetId(fmt.Sprintf("dhcp_helper_%v", svc.Role))
	return diags
}

func resourceCiscoNativeDhcpHelperDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Cisco DHCP IP HELPER DELETE")
	var err error
	var diags diag.Diagnostics

	c, _ := meta.(*providerClient)
	svc := &service.Client{
		Method:   "DELETE",
		Path:     fmt.Sprintf("/data/Cisco-IOS-XE-native:native/interface/Vlan=%v/ip?fields=helper-address", d.Get("svi_id").(int)),
		Provider: c.Provider,
		Devices:  c.Devices.List(),
	}

	roles := d.Get("roles").([]interface{})
	for _, role := range roles {
		svc.Role = role.(string)
		if _, ok := d.GetOk("ipv4_helper"); ok {
			helpers := d.Get("ipv4_helper").([]interface{})
			for _, helper := range helpers {
				svc.Path = fmt.Sprintf("/data/Cisco-IOS-XE-native:native/interface/Vlan=%v/ip/helper-address=%v", d.Get("svi_id").(int), helper.(string))
				_, err = iosxe.MultiSession(svc)
				if err != nil {
					return diag.FromErr(err)
				}
			}
		}
		if _, ok := d.GetOk("source_interface"); ok {
			svc.Path = fmt.Sprintf("/data/Cisco-IOS-XE-native:native/interface/Vlan=%v/ip/dhcp/relay", d.Get("svi_id").(int))
			_, err = iosxe.MultiSession(svc)
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	d.SetId("")
	return diags
}

func (*providerClient) resourceCiscoNativeDhcpHelperCreateIpv4Data(d *schema.ResourceData) *svi.CiscoIOSXENativeVlanDhcpHelper {
	data := &svi.CiscoIOSXENativeVlanDhcpHelper{}

	dhcpData := &svi.CiscoIOSXENativeVlanDhcpHelperSvi{}
	dhcpData.Name = d.Get("svi_id").(int)

	if v, ok := d.GetOk("source_interface"); ok {
		dhcpData.IP.Dhcp.CiscoIOSXEDhcpRelay.SourceInterface = v.(string)
	}
	helpers := d.Get("ipv4_helper").([]interface{})
	for _, helper := range helpers {
		h := &svi.CiscoIOSXENativeVlanDhcpHelperSviHelperAddress{
			Address: helper.(string),
		}
		if d.Get("vrf").(string) == "global" {
			h.Global = null()
		} else {
			log.Panicln("[PANIC] Only global supported for ip helper")
		}
		dhcpData.IP.HelperAddress = append(dhcpData.IP.HelperAddress, *h)
	}
	data.CiscoIOSXENativeVlan = append(data.CiscoIOSXENativeVlan, *dhcpData)
	return data

}
