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
	"github.com/robertcsapo/terraform-provider-ciscoevpn/internal/model/loopback"
	"github.com/robertcsapo/terraform-provider-ciscoevpn/internal/provider/iosxe"
	"github.com/robertcsapo/terraform-provider-ciscoevpn/internal/provider/service"
)

func resourceCiscoNativeBgpNeighbor() *schema.Resource {
	return &schema.Resource{
		Description:   "Cisco BGP Neighbors",
		CreateContext: resourceCiscoNativeBgpNeighborCreate,
		ReadContext:   resourceCiscoNativeBgpNeighborRead,
		UpdateContext: resourceCiscoNativeBgpNeighborUpdate,
		DeleteContext: resourceCiscoNativeBgpNeighborDelete,
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
			"neighbors": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"remote_as": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"update_source": {
				Type:     schema.TypeString,
				Required: true,
			},
			"l2vpn_evpn": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"activate": {
				Type:     schema.TypeBool,
				Default:  true,
				Optional: true,
			},
			"send_community": {
				Type:         schema.TypeString,
				Default:      "both",
				Optional:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"route_reflector_client": {
				Type:     schema.TypeBool,
				Default:  false,
				Optional: true,
			},
			"response": { // TODO remove?
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The HTTP response from the HTTP "GET". The provider will set it to null-value at the time of HTTP "POST", "PATCH", "PUT", and "DELETE."`,
			},
		},
	}
}

func resourceCiscoNativeBgpNeighborCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Cisco BGP NEIGHBORS CREATE")
	var diags diag.Diagnostics
	var err error
	c, _ := meta.(*providerClient)
	svc := &service.Client{
		Provider: c.Provider,
	}

	id := loopbackId(d.Get("update_source").(string))
	d.Set("update_source", id)

	roles := d.Get("roles").([]interface{})
	for _, role := range roles {
		svc.Role = role.(string)
		devices := iosxe.HostRoles(c.Devices.List(), svc.Role) // TODO
		for _, device := range devices {
			svc.Device = device.(string)
			svc.Method = "GET"
			svc.Path = fmt.Sprintf("/data/Cisco-IOS-XE-native:native/interface/Loopback=%v", d.Get("update_source").(string))
			loopback, err := c.loopbackIP(svc)
			if err != nil {
				return diag.FromErr(err)
			}

			data := c.resourceCiscoIOSXEBgpNeighborData(d, loopback.CiscoIOSXENativeLoopback[0].IP.Address.Primary.Address, svc.Role)
			if data == nil {
				return diag.Errorf("No data in yang model")
			}

			if b, err := json.MarshalIndent(data, "", "\t"); err == nil {
				svc.Payload = string(b)
			}
			if svc.Provider.Get("debug").(bool) {
				debugJson(fmt.Sprintf("bgp_neighbors_%v_%v", svc.Role, loopback.CiscoIOSXENativeLoopback[0].IP.Address.Primary.Address), svc.Payload)
			}

			svc.Method = "PATCH"
			svc.Path = "/data/Cisco-IOS-XE-native:native/router/bgp"
			_, err = iosxe.SingleSession(svc)
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	err = d.Set("update_source", fmt.Sprintf("Loopback%v", id))
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(fmt.Sprintf("bgp_id_%v", d.Get("bgp_id").(int)))
	return diags
}

func resourceCiscoNativeBgpNeighborRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Println("[DEBUG] TODO")
	var diags diag.Diagnostics
	return diags
}

func resourceCiscoNativeBgpNeighborUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Cisco BGP NEIGHBORS UPDATE")
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
		Provider: c.Provider,
	}

	id := loopbackId(d.Get("update_source").(string))
	d.Set("update_source", id)

	roles := d.Get("roles").([]interface{})
	for _, role := range roles {
		svc.Role = role.(string)
		devices := iosxe.HostRoles(c.Devices.List(), svc.Role)
		for _, device := range devices {
			svc.Device = device.(string)
			svc.Method = "GET"
			svc.Path = fmt.Sprintf("/data/Cisco-IOS-XE-native:native/interface/Loopback=%v", d.Get("update_source").(string))
			loopback, err := c.loopbackIP(svc)
			if err != nil {
				return diag.FromErr(err)
			}

			data := c.resourceCiscoIOSXEBgpNeighborData(d, loopback.CiscoIOSXENativeLoopback[0].IP.Address.Primary.Address, svc.Role)
			if data == nil {
				return diag.Errorf("No data in yang model")
			}

			if b, err := json.MarshalIndent(data, "", "\t"); err == nil {
				svc.Payload = string(b)
			}
			if svc.Provider.Get("debug").(bool) {
				debugJson(fmt.Sprintf("bgp_neighbors_%v_%v", svc.Role, loopback.CiscoIOSXENativeLoopback[0].IP.Address.Primary.Address), svc.Payload)
			}

			svc.Method = "PATCH" // TODO Read Config and Patch
			svc.Path = "/data/Cisco-IOS-XE-native:native/router/bgp"
			_, err = iosxe.SingleSession(svc)
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	err = d.Set("update_source", fmt.Sprintf("Loopback%v", id))
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(fmt.Sprintf("bgp_id_%v_neighbors", d.Get("bgp_id").(int)))
	return diags
}

func resourceCiscoNativeBgpNeighborDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Cisco BGP NEIGHBORS DELETE")
	var diags diag.Diagnostics
	c, _ := meta.(*providerClient)
	svc := &service.Client{
		Provider: c.Provider,
	}

	id := loopbackId(d.Get("update_source").(string))
	d.Set("update_source", id)

	roles := d.Get("roles").([]interface{})
	for _, role := range roles {
		svc.Role = role.(string)
		devices := iosxe.HostRoles(c.Devices.List(), svc.Role)
		for _, device := range devices {
			svc.Device = device.(string)
			svc.Method = "GET"
			svc.Path = fmt.Sprintf("/data/Cisco-IOS-XE-native:native/interface/Loopback=%v", d.Get("update_source").(string))
			loopback, err := c.loopbackIP(svc)
			if err != nil {
				return diag.FromErr(err)
			}
			for _, id := range d.Get("neighbors").([]interface{}) {
				if id != loopback.CiscoIOSXENativeLoopback[0].IP.Address.Primary.Address {
					svc.Method = "DELETE"

					if d.Get("l2vpn_evpn").(bool) {
						svc.Path = fmt.Sprintf("/data/Cisco-IOS-XE-native:native/router/bgp=%v/address-family/no-vrf/l2vpn/evpn/l2vpn-evpn/neighbor=%v", d.Get("bgp_id").(int), id)
					}
					_, err = iosxe.SingleSession(svc)
					if err != nil {
						return diag.FromErr(err)
					}

					svc.Path = fmt.Sprintf("/data/Cisco-IOS-XE-native:native/router/bgp=%v/neighbor=%v", d.Get("bgp_id").(int), id)
					_, err = iosxe.SingleSession(svc)
					if err != nil {
						return diag.FromErr(err)
					}
				}
			}
		}
	}

	d.SetId("")
	return diags
}

func (*providerClient) loopbackIP(svc *service.Client) (*loopback.CiscoIOSXENativeLoopbackInterface, error) {
	var err error
	var payload string

	if payload, err = iosxe.SingleSession(svc); err != nil {
		return nil, err
	}

	lp := &loopback.CiscoIOSXENativeLoopbackInterface{}
	if err = json.Unmarshal([]byte(payload), &lp); err != nil {
		log.Panicln("[PANIC] Error with JSON data: ", err)
	}
	return lp, nil
}

func (*providerClient) resourceCiscoIOSXEBgpNeighborData(d *schema.ResourceData, localIP string, role string) *bgp.CiscoIOSXEBgpNeighbors {
	var n interface{}
	data := &bgp.CiscoIOSXEBgpNeighbors{}
	system := &bgp.CiscoIOSXEBgp{}
	system.ID = d.Get("bgp_id").(int)
	for _, id := range d.Get("neighbors").([]interface{}) {
		if id != localIP {
			systemNeighbor := &bgp.CiscoIOSXEBgpNeighborsNeighbor{}
			systemNeighbor.ID = id.(string)
			systemNeighbor.RemoteAs = d.Get("remote_as").(int)
			loopback, err := strconv.Atoi(d.Get("update_source").(string))
			if err != nil {
				log.Panicln("[PANIC] Can't find Loopback ID", err)
			}
			systemNeighbor.UpdateSource.Interface.Loopback = loopback
			system.Neighbor = append(system.Neighbor, *systemNeighbor)
		}
	}

	if d.Get("l2vpn_evpn").(bool) {
		EvpnAf := &bgp.CiscoIOSXEBgpNeighborsL2Vpn{}
		EvpnAf.AfName = "evpn"
		for _, id := range d.Get("neighbors").([]interface{}) {
			if id != localIP {
				EvpnNeighbor := &bgp.CiscoIOSXEBgpNeighborsEvpnNeighbor{}
				EvpnNeighbor.ID = id.(string)
				activate := d.Get("activate").(bool)
				if activate {
					EvpnNeighbor.Activate = append(EvpnNeighbor.Activate, n)
				}
				EvpnNeighbor.SendCommunity.SendCommunityWhere = d.Get("send_community").(string)

				if role == "spines" {
					if d.Get("route_reflector_client").(bool) {
						EvpnNeighbor.RouteReflectorClient = append(EvpnNeighbor.RouteReflectorClient, n)
					}
				}

				EvpnAf.L2VpnEvpn.Neighbor = append(EvpnAf.L2VpnEvpn.Neighbor, *EvpnNeighbor)

			}
		}
		system.AddressFamily.NoVrf.L2Vpn = append(system.AddressFamily.NoVrf.L2Vpn, *EvpnAf)
	}
	data.CiscoIOSXEBgpBgp = append(data.CiscoIOSXEBgpBgp, *system)
	return data

}
