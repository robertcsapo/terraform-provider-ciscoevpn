package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/robertcsapo/terraform-provider-ciscoevpn/internal/model/evpn_instance"
	"github.com/robertcsapo/terraform-provider-ciscoevpn/internal/provider/iosxe"
	"github.com/robertcsapo/terraform-provider-ciscoevpn/internal/provider/service"
)

func resourceCiscoNativeEvpnInstance() *schema.Resource {
	return &schema.Resource{
		Description:   "Cisco EVPN Instance",
		CreateContext: resourceCiscoNativeEvpnInstanceCreate,
		ReadContext:   resourceCiscoNativeEvpnInstanceRead,
		UpdateContext: resourceCiscoNativeEvpnInstanceUpdate,
		DeleteContext: resourceCiscoNativeEvpnInstanceDelete,
		Schema: map[string]*schema.Schema{
			"roles": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"instance_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"vlan_based": {
				Type:     schema.TypeBool,
				Default:  true,
				Optional: true,
			},
			"encapsulation": {
				Type:         schema.TypeString,
				Default:      "vxlan",
				Optional:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"replication_type": {
				Type:         schema.TypeString,
				Default:      "static",
				Optional:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"rd": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"rt": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"rt_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"ip_learning": {
				Type:     schema.TypeBool,
				Default:  false,
				Optional: true,
			},
			"default_gateway_advertise": {
				Type:     schema.TypeBool,
				Default:  false,
				Optional: true,
			},
			"re_originate": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
		},
	}
}

func resourceCiscoNativeEvpnInstanceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Cisco EVPN INSTANCE CREATE")
	var diags diag.Diagnostics
	var err error

	c, _ := meta.(*providerClient)
	svc := &service.Client{
		Method:   "PATCH",
		Path:     "/data/Cisco-IOS-XE-native:native/l2vpn/evpn_cont/evpn-instance/evpn/instance",
		Provider: c.Provider,
		Devices:  c.Devices.List(),
	}

	roles := d.Get("roles").([]interface{})
	for _, role := range roles {
		svc.Role = role.(string)
		data := c.CiscoIOSXENativeEvpnInstanceData(d)
		if data == nil {
			return diag.Errorf("No data in yang model")
		}

		if b, err := json.MarshalIndent(data, "", "\t"); err == nil {
			svc.Payload = string(b)
		}
		if svc.Provider.Get("debug").(bool) {
			debugJson(fmt.Sprintf("evpn_instance_%v", d.Get("instance_id").(int)), svc.Payload)
		}

		_, err = iosxe.MultiSession(svc)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId(fmt.Sprintf("evpn_instance_%v", d.Get("instance_id").(int)))
	return diags
}

func resourceCiscoNativeEvpnInstanceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Cisco EVPN INSTANCE READ") // TODO
	var diags diag.Diagnostics
	return diags
}

func resourceCiscoNativeEvpnInstanceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Cisco EVPN INSTANCE UPDATE")
	var diags diag.Diagnostics
	var err error

	if d.HasChange("instance_id") {
		oldState, _ := d.GetChange("instance_id")
		d.Set("instance_id", oldState)
		return diag.Errorf("Not supported to change EVPN Instance ID number")
	}
	if d.HasChange("roles") {
		oldState, _ := d.GetChange("roles")
		d.Set("roles", oldState)
		return diag.Errorf("Not supported to change Roles")
	}

	c, _ := meta.(*providerClient)
	svc := &service.Client{
		Method:   "PATCH",
		Path:     "/data/Cisco-IOS-XE-native:native/l2vpn/evpn_cont/evpn-instance/evpn/instance",
		Provider: c.Provider,
		Devices:  c.Devices.List(),
	}

	roles := d.Get("roles").([]interface{})
	for _, role := range roles {
		svc.Role = role.(string)
		data := c.CiscoIOSXENativeEvpnInstanceData(d)
		if data == nil {
			return diag.Errorf("No data in yang model")
		}

		if b, err := json.MarshalIndent(data, "", "\t"); err == nil {
			svc.Payload = string(b)
		}
		if svc.Provider.Get("debug").(bool) {
			debugJson(fmt.Sprintf("evpn_instance_%v", d.Get("instance_id").(int)), svc.Payload)
		}

		_, err = iosxe.MultiSession(svc)

		if err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId(fmt.Sprintf("evpn_instance_%v", d.Get("instance_id").(int)))
	return diags
}

func resourceCiscoNativeEvpnInstanceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Cisco EVPN INSTANCE DELETE")
	var diags diag.Diagnostics
	c, _ := meta.(*providerClient)
	svc := &service.Client{
		Method:   "DELETE",
		Path:     fmt.Sprintf("/data/Cisco-IOS-XE-native:native/l2vpn/evpn_cont/evpn-instance/evpn/instance/instance=%v", d.Get("instance_id").(int)),
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

func (*providerClient) CiscoIOSXENativeEvpnInstanceData(d *schema.ResourceData) *evpn_instance.CiscoIOSXEL2VpnInstanceCiscoIOSXEL2VpnInstanceEvpn {
	data := &evpn_instance.CiscoIOSXEL2VpnInstanceCiscoIOSXEL2VpnInstanceEvpn{}
	ei := &evpn_instance.CiscoIOSXEL2VpnInstanceInstance{}
	ei.EvpnInstanceNum = d.Get("instance_id").(int)
	if d.Get("vlan_based").(bool) {
		ei.VlanBased.Encapsulation = d.Get("encapsulation").(string)
		if v, ok := d.GetOk("replication_type"); ok {
			switch v.(string) {
			case "static":
				ei.VlanBased.ReplicationType.Static = null()
			case "ingress":
				ei.VlanBased.ReplicationType.Ingress = null()
			default:
				log.Panicf("[PANIC] Replication Type (%v) not supported", v.(string))
			}
			/*
				if v.(string) == "ingress" {
					ei.VlanBased.ReplicationType.Ingress = null()
				} else {
				if v.(string) == "ingress" {
					ei.VlanBased.ReplicationType.Ingress = null()
				} else {
					log.Panicf("[PANIC] Replication Type (%v) not supported", v.(string))
				}
			*/
		}
		if v, ok := d.GetOk("rd"); ok {
			ei.VlanBased.Rd.RdValue = v.(string)
		}
		if v, ok := d.GetOk("rt_type"); ok {
			if v.(string) == "both" {
				ei.VlanBased.RouteTarget.Both.RtValue = d.Get("rd").(string)
			} else {
				log.Panicf("[PANIC] RT Type (%v) not supported", v.(string))
			}
		}
		if v, ok := d.GetOk("re_originate"); ok {
			if v.(string) == "route-type5" {
				ei.VlanBased.ReOriginate.RouteType5 = null()
			} else {
				log.Panicf("[PANIC] Reoriginates Type (%v) not supported", v.(string))
			}
		}

		if !d.Get("ip_learning").(bool) {
			ei.VlanBased.IP.LocalLearning.Disable = null()
		}
		if d.Get("default_gateway_advertise").(bool) {
			ei.VlanBased.DefaultGateway.Advertise = "enable"
		} else {
			ei.VlanBased.DefaultGateway.Advertise = "disable"
		}
	}
	data.CiscoIOSXEL2VpnInstance.Instance = append(data.CiscoIOSXEL2VpnInstance.Instance, *ei)
	return data
}
