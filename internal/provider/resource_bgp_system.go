package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/robertcsapo/terraform-provider-ciscoevpn/internal/model/bgp"
	"github.com/robertcsapo/terraform-provider-ciscoevpn/internal/provider/iosxe"
	"github.com/robertcsapo/terraform-provider-ciscoevpn/internal/provider/service"
)

func resourceCiscoNativeBgpSystem() *schema.Resource {
	return &schema.Resource{
		Description:   "Cisco BGP System",
		CreateContext: resourceCiscoNativeBgpSystemCreate,
		ReadContext:   resourceCiscoNativeBgpSystemRead,
		UpdateContext: resourceCiscoNativeBgpSystemUpdate,
		DeleteContext: resourceCiscoNativeBgpSystemDelete,
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
			"router_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"log_neighbor_changes": {
				Type:     schema.TypeBool,
				Default:  true,
				Optional: true,
			},
			"default_ipv4_unicast": {
				Type:     schema.TypeBool,
				Default:  false,
				Optional: true,
			},
		},
	}
}

func resourceCiscoNativeBgpSystemCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Cisco BGP SYSTEM CREATE")
	var diags diag.Diagnostics
	var err error
	c, _ := meta.(*providerClient)
	svc := &service.Client{
		Method:   "PATCH",
		Path:     "/data/Cisco-IOS-XE-native:native/router/bgp",
		Provider: c.Provider,
		Devices:  c.Devices.List(),
	}

	id := loopbackId(d.Get("router_id").(string))
	d.Set("router_id", id)

	roles := d.Get("roles").([]interface{})
	for _, role := range roles {
		svc.Role = role.(string)
		data := c.resourceCiscoIOSXEBgpSystemData(d)
		if data == nil {
			return diag.Errorf("No data in yang model")
		}

		if b, err := json.MarshalIndent(data, "", "\t"); err == nil {
			svc.Payload = string(b)
		}
		if svc.Provider.Get("debug").(bool) {
			debugJson(fmt.Sprintf("bgp_system_%v_%v", svc.Role, d.Get("bgp_id")), svc.Payload)
		}

		_, err = iosxe.MultiSession(svc)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	err = d.Set("router_id", fmt.Sprintf("Loopback%v", id))
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(fmt.Sprintf("bgp_systems_%v", d.Get("bgp_id").(int)))
	return diags
}

func resourceCiscoNativeBgpSystemRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Println("[DEBUG] TODO")
	return nil
}

func resourceCiscoNativeBgpSystemUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Cisco BGP SYSTEM UPDATE")
	var diags diag.Diagnostics
	var err error
	if d.HasChange("bgp_id") {
		oldState, _ := d.GetChange("bgp_id")
		d.Set("bgp_id", oldState)
		return diag.Errorf("Not supported to change BGP AS number")
	}
	if d.HasChange("roles") {
		oldState, _ := d.GetChange("roles")
		d.Set("roles", oldState)
		return diag.Errorf("Not supported to change Roles")
	}
	c, _ := meta.(*providerClient)
	svc := &service.Client{
		Method:   "PATCH",
		Path:     fmt.Sprintf("/data/Cisco-IOS-XE-native:native/router/bgp=%v", d.Get("bgp_id").(int)),
		Provider: c.Provider,
		Devices:  c.Devices.List(),
	}

	// TODO support interface vs ip_address
	id := loopbackId(d.Get("router_id").(string))
	d.Set("router_id", id)

	roles := d.Get("roles").([]interface{})
	for _, role := range roles {
		svc.Role = role.(string)
		data := c.resourceCiscoIOSXEBgpSystemData(d)
		if data == nil {
			return diag.Errorf("No data in yang model")
		}

		if b, err := json.MarshalIndent(data, "", "\t"); err == nil {
			svc.Payload = string(b)
		}
		if svc.Provider.Get("debug").(bool) {
			debugJson(fmt.Sprintf("bgp_system_%v_%v", svc.Role, d.Get("bgp_id")), svc.Payload)
		}

		_, err = iosxe.MultiSession(svc)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	err = d.Set("router_id", fmt.Sprintf("Loopback%v", id))
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(fmt.Sprintf("bgp_systems_%v", d.Get("bgp_id").(int)))
	return diags
}

func resourceCiscoNativeBgpSystemDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Cisco BGP System DELETE")
	var diags diag.Diagnostics
	c, _ := meta.(*providerClient)
	svc := &service.Client{
		Method:   "DELETE",
		Path:     fmt.Sprintf("/data/Cisco-IOS-XE-native:native/router/bgp=%v", d.Get("bgp_id").(int)),
		Provider: c.Provider,
		Devices:  c.Devices.List(),
	}

	roles := d.Get("roles").([]interface{})
	for _, role := range roles {
		svc.Role = role.(string)
		_, err := iosxe.MultiSession(svc)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	d.SetId("")
	return diags
}

func (*providerClient) resourceCiscoIOSXEBgpSystemData(d *schema.ResourceData) *bgp.CiscoIOSXEBgpBgpSystem {
	data := &bgp.CiscoIOSXEBgpBgpSystem{}
	system := &bgp.CiscoIOSXEBgp{}
	system.ID = d.Get("bgp_id").(int)
	system.Bgp.LogNeighborChanges = d.Get("log_neighbor_changes").(bool)
	id, err := strconv.Atoi(d.Get("router_id").(string))
	if err != nil {
		log.Panicln("[PANIC] Can't find Loopback ID", err)
	}
	system.Bgp.RouterID.Interface.Loopback = id
	if d.Get("default_ipv4_unicast").(bool) {
		system.Bgp.Default.Ipv4Unicast = true
	} else {
		system.Bgp.Default.Ipv4Unicast = false
	}
	data.CiscoIOSXEBgpBgp = append(data.CiscoIOSXEBgpBgp, *system)
	return data
}
