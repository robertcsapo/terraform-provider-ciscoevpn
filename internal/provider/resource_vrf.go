package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/robertcsapo/terraform-provider-ciscoevpn/internal/model/vrf"
	"github.com/robertcsapo/terraform-provider-ciscoevpn/internal/provider/iosxe"
	"github.com/robertcsapo/terraform-provider-ciscoevpn/internal/provider/service"
)

func resourceCiscoNativeVrf() *schema.Resource {
	return &schema.Resource{
		Description:   "Cisco VRF",
		CreateContext: resourceCiscoNativeVrfCreate,
		ReadContext:   resourceCiscoNativeVrfRead,
		UpdateContext: resourceCiscoNativeVrfUpdate,
		DeleteContext: resourceCiscoNativeVrfDelete,
		Schema: map[string]*schema.Schema{
			"roles": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"rd": {
				Type:     schema.TypeString,
				Required: true,
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
		},
	}
}

func resourceCiscoNativeVrfCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Cisco VRF CREATE")
	var diags diag.Diagnostics
	var err error

	c, _ := meta.(*providerClient)
	svc := &service.Client{
		Method:   "PATCH",
		Path:     "/data/Cisco-IOS-XE-native:native/vrf/definition",
		Provider: c.Provider,
		Devices:  c.Devices.List(),
	}

	roles := d.Get("roles").([]interface{})
	for _, role := range roles {
		svc.Role = role.(string)
		data := c.CiscoIOSXENativeVrfData(d)
		if data == nil {
			return diag.Errorf("No data in yang model")
		}

		if b, err := json.MarshalIndent(data, "", "\t"); err == nil {
			svc.Payload = string(b)
		}
		if svc.Provider.Get("debug").(bool) {
			debugJson(fmt.Sprintf("vrf_%v_%v", svc.Role, d.Get("name").(string)), svc.Payload)
		}

		_, err = iosxe.MultiSession(svc)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId(fmt.Sprintf("%v", d.Get("name").(string)))
	return diags
}

func resourceCiscoNativeVrfRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Println("[DEBUG] TODO")
	var diags diag.Diagnostics
	return diags
}

func resourceCiscoNativeVrfUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Cisco VRF UPDATE")
	var diags diag.Diagnostics
	var err error

	if d.HasChange("name") {
		oldState, _ := d.GetChange("name")
		d.Set("name", oldState)
		return diag.Errorf("Not supported to change Name of VRF")
	}
	if d.HasChange("roles") {
		oldState, _ := d.GetChange("roles")
		d.Set("roles", oldState)
		return diag.Errorf("Not supported to change Roles")
	}
	c, _ := meta.(*providerClient)
	svc := &service.Client{
		Method:   "PATCH",
		Path:     fmt.Sprintf("/data/Cisco-IOS-XE-native:native/vrf/definition=%v", d.Get("name").(string)),
		Provider: c.Provider,
		Devices:  c.Devices.List(),
	}

	roles := d.Get("roles").([]interface{})
	for _, role := range roles {
		svc.Role = role.(string)
		data := c.CiscoIOSXENativeVrfData(d)
		if data == nil {
			return diag.Errorf("No data in yang model")
		}

		if b, err := json.MarshalIndent(data, "", "\t"); err == nil {
			svc.Payload = string(b)
		}
		if svc.Provider.Get("debug").(bool) {
			debugJson(fmt.Sprintf("vrf_%v_%v", svc.Role, d.Get("name").(string)), svc.Payload)
		}

		_, err = iosxe.MultiSession(svc)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId(fmt.Sprintf("%v", d.Get("name").(string)))
	return diags
}

func resourceCiscoNativeVrfDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Cisco VRF DELETE")
	var err error
	var diags diag.Diagnostics

	c, _ := meta.(*providerClient)
	svc := &service.Client{
		Method:   "DELETE",
		Path:     fmt.Sprintf("/data/Cisco-IOS-XE-native:native/vrf/definition=%v", d.Get("name").(string)),
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

func (*providerClient) CiscoIOSXENativeVrfData(d *schema.ResourceData) *vrf.CiscoIOSXENativeVrf {
	data := &vrf.CiscoIOSXENativeVrf{}
	vrfData := &vrf.CiscoIOSXENativeDefinition{}

	vrfData.Name = d.Get("name").(string)
	vrfData.Rd = d.Get("rd").(string)

	withoutStiching := &vrf.CiscoIOSXENativeDefinitionWithoutStitching{
		AsnIP: vrfData.Rd,
	}
	withStiching := &vrf.CiscoIOSXENativeDefinitionWithStitching{
		AsnIP:     vrfData.Rd,
		Stitching: null(),
	}
	if d.Get("ipv4").(bool) {
		vrfData.AddressFamily.Ipv4.RouteTarget.ImportRouteTarget.WithoutStitching = append(vrfData.AddressFamily.Ipv4.RouteTarget.ImportRouteTarget.WithoutStitching, *withoutStiching)
		vrfData.AddressFamily.Ipv4.RouteTarget.ImportRouteTarget.WithStitching = append(vrfData.AddressFamily.Ipv4.RouteTarget.ImportRouteTarget.WithStitching, *withStiching)
		vrfData.AddressFamily.Ipv4.RouteTarget.ExportRouteTarget.WithoutStitching = append(vrfData.AddressFamily.Ipv4.RouteTarget.ExportRouteTarget.WithoutStitching, *withoutStiching)
		vrfData.AddressFamily.Ipv4.RouteTarget.ExportRouteTarget.WithStitching = append(vrfData.AddressFamily.Ipv4.RouteTarget.ExportRouteTarget.WithStitching, *withStiching)
	}
	if d.Get("ipv6").(bool) {
		vrfData.AddressFamily.Ipv6.RouteTarget.ImportRouteTarget.WithoutStitching = append(vrfData.AddressFamily.Ipv6.RouteTarget.ImportRouteTarget.WithoutStitching, *withoutStiching)
		vrfData.AddressFamily.Ipv6.RouteTarget.ImportRouteTarget.WithStitching = append(vrfData.AddressFamily.Ipv6.RouteTarget.ImportRouteTarget.WithStitching, *withStiching)
		vrfData.AddressFamily.Ipv6.RouteTarget.ExportRouteTarget.WithoutStitching = append(vrfData.AddressFamily.Ipv6.RouteTarget.ExportRouteTarget.WithoutStitching, *withoutStiching)
		vrfData.AddressFamily.Ipv6.RouteTarget.ExportRouteTarget.WithStitching = append(vrfData.AddressFamily.Ipv6.RouteTarget.ExportRouteTarget.WithStitching, *withStiching)
	}
	data.CiscoIOSXENativeDefinition = append(data.CiscoIOSXENativeDefinition, *vrfData)
	return data

}
