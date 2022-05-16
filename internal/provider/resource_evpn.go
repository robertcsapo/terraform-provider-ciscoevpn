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
	"github.com/robertcsapo/terraform-provider-ciscoevpn/internal/model/evpn"
	"github.com/robertcsapo/terraform-provider-ciscoevpn/internal/provider/iosxe"
	"github.com/robertcsapo/terraform-provider-ciscoevpn/internal/provider/service"
)

func resourceCiscoNativeL2VpnEvpn() *schema.Resource {
	return &schema.Resource{
		Description: "Cisco L2VPN EVPN",

		CreateContext: resourceCiscoNativeL2VpnEvpnCreate,
		ReadContext:   resourceCiscoNativeL2VpnEvpnRead,
		UpdateContext: resourceCiscoNativeL2VpnEvpnUpdate, //Todo Update
		DeleteContext: resourceCiscoNativeL2VpnEvpnDelete,

		Schema: map[string]*schema.Schema{
			"roles": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"replication_type": {
				Type:         schema.TypeString,
				Default:      "static",
				Optional:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"mac_duplication_limit": {
				Type:     schema.TypeInt,
				Default:  20,
				Optional: true,
			},
			"mac_duplication_time": {
				Type:     schema.TypeInt,
				Default:  10,
				Optional: true,
			},
			"ip_duplication_limit": {
				Type:     schema.TypeInt,
				Default:  20,
				Optional: true,
			},
			"ip_duplication_time": {
				Type:     schema.TypeInt,
				Default:  10,
				Optional: true,
			},
			"router_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"default_gateway": {
				Type:     schema.TypeString,
				Default:  "advertise",
				Optional: true,
			},
			"logging_peer_state": {
				Type:     schema.TypeBool,
				Default:  true,
				Optional: true,
			},
			"route_target_auto": {
				Type:     schema.TypeString,
				Default:  "vni",
				Optional: true,
			},
		},
	}
}

func resourceCiscoNativeL2VpnEvpnCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Cisco L2VPN EVPN CREATE")
	var diags diag.Diagnostics
	var err error
	c, _ := meta.(*providerClient)
	svc := &service.Client{
		Method:   "PATCH",
		Path:     "/data/Cisco-IOS-XE-native:native/l2vpn/evpn_cont/evpn",
		Provider: c.Provider,
		Devices:  c.Devices.List(),
	}

	id := loopbackId(d.Get("router_id").(string))
	d.Set("router_id", id)

	roles := d.Get("roles").([]interface{})
	for _, role := range roles {
		svc.Role = role.(string)
		data := c.resourceCiscoNativeL2VpnEvpnData(d)
		if data == nil {
			return diag.Errorf("No data in yang model") // TODO refactor
		}

		if b, err := json.MarshalIndent(data, "", "\t"); err == nil {
			svc.Payload = string(b)
		}
		if svc.Provider.Get("debug").(bool) {
			debugJson(fmt.Sprintf("evpn_loopback%v", d.Get("router_id").(string)), svc.Payload) // TODO
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
	d.SetId(fmt.Sprintf("l2vpn_evpn_%v", d.Get("router_id").(string)))
	return diags
}

func resourceCiscoNativeL2VpnEvpnRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Cisco EVPN READ TODO")
	var diags diag.Diagnostics
	return diags
}

func resourceCiscoNativeL2VpnEvpnUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Cisco L2VPN EVPN UPDATE")
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
		Path:     "/data/Cisco-IOS-XE-native:native/l2vpn/evpn_cont/evpn",
		Provider: c.Provider,
		Devices:  c.Devices.List(),
	}
	id := loopbackId(d.Get("router_id").(string))
	d.Set("router_id", id)

	roles := d.Get("roles").([]interface{})
	for _, role := range roles {
		svc.Role = fmt.Sprintf("%v", role)
		data := c.resourceCiscoNativeL2VpnEvpnData(d)
		if data == nil {
			return diag.Errorf("No data in yang model") // TODO refactor
		}

		if b, err := json.MarshalIndent(data, "", "\t"); err == nil {
			svc.Payload = string(b)
		}
		if svc.Provider.Get("debug").(bool) {
			debugJson(fmt.Sprintf("evpn_loopback%v", d.Get("router_id").(string)), svc.Payload) // TODO
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
	d.SetId(fmt.Sprintf("l2vpn_evpn_%v", d.Get("router_id").(string)))
	return diags
}

func resourceCiscoNativeL2VpnEvpnDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Cisco L2VPN EVPN DELETE")
	var diags diag.Diagnostics
	var err error
	c, _ := meta.(*providerClient)
	svc := &service.Client{
		Method:   "DELETE",
		Path:     "/data/Cisco-IOS-XE-native:native/l2vpn/evpn_cont/evpn",
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

func (*providerClient) resourceCiscoNativeL2VpnEvpnData(d *schema.ResourceData) *evpn.CiscoIOSXEL2Evpn {
	var err error
	empty := make([]int, 0) // TODO
	data := &evpn.CiscoIOSXEL2Evpn{}
	if d.Get("replication_type").(string) == "static" {
		data.CiscoIOSXEL2VpnEvpn.ReplicationType.Static = empty
	}
	data.CiscoIOSXEL2VpnEvpn.Mac.Duplication.Limit = d.Get("mac_duplication_limit").(int)
	data.CiscoIOSXEL2VpnEvpn.Mac.Duplication.Time = d.Get("mac_duplication_time").(int)
	data.CiscoIOSXEL2VpnEvpn.IP.Duplication.Limit = d.Get("ip_duplication_limit").(int)
	data.CiscoIOSXEL2VpnEvpn.IP.Duplication.Time = d.Get("ip_duplication_time").(int)

	if data.CiscoIOSXEL2VpnEvpn.RouterID.Interface.Loopback, err = strconv.Atoi(d.Get("router_id").(string)); err != nil {
		log.Panicln("[PANIC] Can't find Loopback ID", err)
	}
	if d.Get("default_gateway").(string) == "advertise" {
		data.CiscoIOSXEL2VpnEvpn.DefaultGateway.Advertise = empty
	}
	if d.Get("logging_peer_state").(bool) {
		data.CiscoIOSXEL2VpnEvpn.Logging.Peer.State = empty
	}
	if d.Get("route_target_auto").(string) == "vni" {
		data.CiscoIOSXEL2VpnEvpn.RouteTarget.Auto.Vni = empty
	}

	return data
}
