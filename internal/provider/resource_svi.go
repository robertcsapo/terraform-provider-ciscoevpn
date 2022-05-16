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

func resourceCiscoNativeSvi() *schema.Resource {
	return &schema.Resource{
		Description:   "Cisco SVI",
		CreateContext: resourceCiscoNativeSviCreate,
		ReadContext:   resourceCiscoNativeSviRead,
		UpdateContext: resourceCiscoNativeSviUpdate,
		DeleteContext: resourceCiscoNativeSviDelete,
		Schema: map[string]*schema.Schema{
			"roles": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"svi_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"autostate": {
				Type:     schema.TypeBool,
				Default:  false,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Default:  "Managed by Terraform (ciscoevpn)",
				Optional: true,
			},
			"vrf": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ipv4_address": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.IsIPv4Address,
			},
			"ipv4_mask": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.IsIPv4Address,
			},
			"unnumbered": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
		},
	}
}

func resourceCiscoNativeSviCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Cisco SVI CREATE")
	var diags diag.Diagnostics
	var err error

	c, _ := meta.(*providerClient)
	svc := &service.Client{
		Method:   "PATCH",
		Path:     "/data/Cisco-IOS-XE-native:native/interface/Vlan",
		Provider: c.Provider,
		Devices:  c.Devices.List(),
	}

	roles := d.Get("roles").([]interface{})
	for _, role := range roles {
		svc.Role = role.(string)
		data := c.resourceCiscoNativeSviData(d)
		if data == nil {
			return diag.Errorf("No data in yang model")
		}

		if b, err := json.MarshalIndent(data, "", "\t"); err == nil {
			svc.Payload = string(b)
		}
		if svc.Provider.Get("debug").(bool) {
			debugJson(fmt.Sprintf("svi_%v_%v", svc.Role, d.Get("svi_id").(int)), svc.Payload)
		}

		_, err = iosxe.MultiSession(svc)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId(fmt.Sprintf("svi_%v", d.Get("svi_id").(int)))
	return diags
}

func resourceCiscoNativeSviRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Cisco SVI READ") // TODO
	var diags diag.Diagnostics
	return diags
}

func resourceCiscoNativeSviUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Cisco SVI UPDATE")
	var diags diag.Diagnostics
	var err error

	if d.HasChange("svi_id") {
		oldState, _ := d.GetChange("svi_id")
		d.Set("svi_id", oldState)
		return diag.Errorf("Not supported to change VLAN ID")
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
		data := c.resourceCiscoNativeSviData(d)
		if data == nil {
			return diag.Errorf("No data in yang model")
		}

		if b, err := json.MarshalIndent(data, "", "\t"); err == nil {
			svc.Payload = string(b)
		}
		if svc.Provider.Get("debug").(bool) {
			debugJson(fmt.Sprintf("svi_%v_%v", svc.Role, d.Get("svi_id").(int)), svc.Payload)
		}

		_, err = iosxe.MultiSession(svc)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId(fmt.Sprintf("svi_%v", d.Get("svi_id").(int)))
	return diags
}

func resourceCiscoNativeSviDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Cisco SVI DELETE")
	var diags diag.Diagnostics
	var err error

	c, _ := meta.(*providerClient)
	svc := &service.Client{
		Method:   "DELETE",
		Path:     fmt.Sprintf("/data/Cisco-IOS-XE-native:native/interface/Vlan=%v", d.Get("svi_id").(int)),
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

func (*providerClient) resourceCiscoNativeSviData(d *schema.ResourceData) *svi.CiscoIOSXENativeSvis {
	data := &svi.CiscoIOSXENativeSvis{}
	sviCfg := &svi.CiscoIOSXENativeSvi{}

	sviCfg.Name = d.Get("svi_id").(int)
	sviCfg.AutoState = d.Get("autostate").(bool)
	sviCfg.Description = d.Get("description").(string)
	if v, ok := d.GetOk("vrf"); ok {
		vrf := &svi.CiscoIOSXENativeVlanVrf{
			Forwarding: v.(string),
		}
		sviCfg.Vrf = vrf
	}

	if _, ok := d.GetOk("ipv4_address"); ok {
		sviIp := &svi.CiscoIOSXENativeVlanAddress{}
		sviIp.Primary.Address = d.Get("ipv4_address").(string)
		sviIp.Primary.Mask = d.Get("ipv4_mask").(string)
		sviCfg.IP.Address = sviIp
	} else if v, ok := d.GetOk("unnumbered"); ok {
		sviCfg.IP.Unnumbered = v.(string)
	} else {
		log.Panicln("[PANIC] Either IPv4 or Unnumbered has to be used")
	}

	data.CiscoIOSXENativeVlan = append(data.CiscoIOSXENativeVlan, *sviCfg)
	return data
}
