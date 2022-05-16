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

func resourceCiscoNativeBgpNeighborVrfUnicast() *schema.Resource {
	return &schema.Resource{
		Description:   "Cisco BGP Neighbor VRF Unicast",
		CreateContext: resourceCiscoNativeBgpNeighborVrfUnicastCreate,
		ReadContext:   resourceCiscoNativeBgpNeighborVrfUnicastRead,
		UpdateContext: resourceCiscoNativeBgpNeighborVrfUnicastUpdate,
		DeleteContext: resourceCiscoNativeBgpNeighborVrfUnicastDelete,
		Schema: map[string]*schema.Schema{
			"host": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"bgp_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"remote_as": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"activate": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"vrf": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"ipv4_neighbors": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceCiscoNativeBgpNeighborVrfUnicastCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Cisco NEIGHBORS VRF UNICAST CREATE")
	var diags diag.Diagnostics
	var err error

	c, _ := meta.(*providerClient)
	svc := &service.Client{
		Method:   "PATCH",
		Provider: c.Provider,
		Device:   d.Get("host").(string),
	}
	if _, ok := d.GetOk("ipv4_neighbors"); ok {
		data := c.resourceCiscoNativeBgpNeighborVrfUnicastIpv4Data(d)
		svc.Path = fmt.Sprintf("/data/Cisco-IOS-XE-native:native/router/bgp=%v/address-family/with-vrf/ipv4/unicast/vrf=%v/ipv4-unicast/", d.Get("bgp_id").(int), d.Get("vrf").(string))
		if data == nil {
			return diag.Errorf("No data in yang model")
		}

		if b, err := json.MarshalIndent(data, "", "\t"); err == nil {
			svc.Payload = string(b)
		}
		if svc.Provider.Get("debug").(bool) {
			debugJson(fmt.Sprintf("bgp_neighbor_vrf_unicast_ipv4_%v_%v", svc.Device, d.Get("vrf").(string)), svc.Payload)
		}

		_, err = iosxe.SingleSession(svc)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	d.SetId(fmt.Sprintf("bgp_neighbor_vrf_unicast_%v", d.Get("vrf").(string)))
	return diags
}

func resourceCiscoNativeBgpNeighborVrfUnicastRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Println("[DEBUG] TODO")
	return nil
}

func resourceCiscoNativeBgpNeighborVrfUnicastUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Cisco NEIGHBORS VRF UNICAST UPDATE")
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

	c, _ := meta.(*providerClient)
	svc := &service.Client{
		Method:   "PATCH",
		Provider: c.Provider,
		Device:   d.Get("host").(string),
	}

	if _, ok := d.GetOk("ipv4_neighbors"); ok {
		data := c.resourceCiscoNativeBgpNeighborVrfUnicastIpv4Data(d)
		svc.Path = fmt.Sprintf("/data/Cisco-IOS-XE-native:native/router/bgp=%v/address-family/with-vrf/ipv4/unicast/vrf=%v/ipv4-unicast/", d.Get("bgp_id").(int), d.Get("vrf").(string))
		if data == nil {
			return diag.Errorf("No data in yang model")
		}

		if b, err := json.MarshalIndent(data, "", "\t"); err == nil {
			svc.Payload = string(b)
		}
		if svc.Provider.Get("debug").(bool) {
			debugJson(fmt.Sprintf("bgp_neighbor_vrf_unicast_ipv4_%v_%v", svc.Device, d.Get("vrf").(string)), svc.Payload)
		}

		_, err = iosxe.SingleSession(svc)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	d.SetId(fmt.Sprintf("bgp_neighbor_vrf_unicast_%v", d.Get("vrf").(string)))
	return diags
}

func resourceCiscoNativeBgpNeighborVrfUnicastDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Cisco NEIGHBORS VRF UNICAST DELETE")
	var diags diag.Diagnostics
	var err error
	c, _ := meta.(*providerClient)
	svc := &service.Client{
		Method:   "DELETE",
		Provider: c.Provider,
		Device:   d.Get("host").(string),
	}

	if _, ok := d.GetOk("ipv4_neighbors"); ok {
		neighbors := d.Get("ipv4_neighbors").([]interface{})
		for _, neighbor := range neighbors {
			svc.Path = fmt.Sprintf("/data/Cisco-IOS-XE-native:native/router/bgp=%v/address-family/with-vrf/ipv4/unicast/vrf=%v/ipv4-unicast/neighbor=%v", d.Get("bgp_id").(int), d.Get("vrf").(string), neighbor.(string))
			_, err = iosxe.SingleSession(svc)
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	d.SetId("")
	return diags
}

func (*providerClient) resourceCiscoNativeBgpNeighborVrfUnicastIpv4Data(d *schema.ResourceData) *bgp.CiscoIOSXEBgpVrfIpv4Unicast {
	data := &bgp.CiscoIOSXEBgpVrfIpv4Unicast{}

	neighbors := d.Get("ipv4_neighbors").([]interface{})
	for _, neighbor := range neighbors {
		n := &bgp.CiscoIOSXEBgpIpv4UnicastNeighbor{
			ID:       neighbor.(string),
			RemoteAs: d.Get("remote_as").(int),
		}
		if d.Get("activate").(bool) {
			n.Activate = null()
		}
		data.CiscoIOSXEBgpIpv4Unicast.Neighbor = append(data.CiscoIOSXEBgpIpv4Unicast.Neighbor, *n)
	}
	return data
}
