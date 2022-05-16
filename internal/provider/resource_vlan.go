package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/robertcsapo/terraform-provider-ciscoevpn/internal/model/vlan"
	"github.com/robertcsapo/terraform-provider-ciscoevpn/internal/provider/iosxe"
	"github.com/robertcsapo/terraform-provider-ciscoevpn/internal/provider/service"
)

func resourceCiscoNativeVlan() *schema.Resource {
	return &schema.Resource{
		Description:   "Cisco VLAN",
		CreateContext: resourceCiscoNativeVlanCreate,
		ReadContext:   resourceCiscoNativeVlanRead,
		UpdateContext: resourceCiscoNativeVlanUpdate,
		DeleteContext: resourceCiscoNativeVlanDelete,
		Schema: map[string]*schema.Schema{
			"roles": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"vlan_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"name": {
				Type:         schema.TypeString,
				Default:      "ManagedByTerraform",
				Optional:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"evpn_instance": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"vni": {
				Type:     schema.TypeInt,
				Required: true,
			},
		},
	}
}

func resourceCiscoNativeVlanCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Cisco VLAN CREATE")
	var diags diag.Diagnostics
	var err error

	c, _ := meta.(*providerClient)
	svc := &service.Client{
		Method:   "PATCH",
		Path:     "/data/Cisco-IOS-XE-native:native/vlan",
		Provider: c.Provider,
		Devices:  c.Devices.List(),
	}

	roles := d.Get("roles").([]interface{})
	for _, role := range roles {
		svc.Role = role.(string)
		data := c.resourceCiscoNativeVlanData(d)
		if data == nil {
			return diag.Errorf("No data in yang model")
		}

		if b, err := json.MarshalIndent(data, "", "\t"); err == nil {
			svc.Payload = string(b)
		}
		if svc.Provider.Get("debug").(bool) {
			debugJson(fmt.Sprintf("vlan_%v_%v", svc.Role, d.Get("vlan_id").(int)), svc.Payload)
		}

		_, err = iosxe.MultiSession(svc)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId(fmt.Sprintf("vlan_%v", d.Get("vlan_id").(int)))
	return diags
}

func resourceCiscoNativeVlanRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Cisco VLAN READ") // TODO
	var diags diag.Diagnostics
	return diags
}

func resourceCiscoNativeVlanUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Cisco VLAN UPDATE")
	var diags diag.Diagnostics
	var err error

	if d.HasChange("vlan_id") {
		oldState, _ := d.GetChange("vlan_id")
		d.Set("vlan_id", oldState)
		return diag.Errorf("Not supported to change VLAN ID")
	}
	if d.HasChange("vni") {
		oldState, _ := d.GetChange("vni")
		d.Set("vni", oldState)
		return diag.Errorf("Not supported to change VNI ID")
	}
	if d.HasChange("roles") {
		oldState, _ := d.GetChange("roles")
		d.Set("roles", oldState)
		return diag.Errorf("Not supported to change Roles")
	}

	c, _ := meta.(*providerClient)
	svc := &service.Client{
		Method:   "PATCH",
		Path:     "/data/Cisco-IOS-XE-native:native/vlan",
		Provider: c.Provider,
		Devices:  c.Devices.List(),
	}

	roles := d.Get("roles").([]interface{})
	for _, role := range roles {
		svc.Role = role.(string)
		data := c.resourceCiscoNativeVlanData(d)
		if data == nil {
			return diag.Errorf("No data in yang model")
		}

		if b, err := json.MarshalIndent(data, "", "\t"); err == nil {
			svc.Payload = string(b)
		}
		if svc.Provider.Get("debug").(bool) {
			debugJson(fmt.Sprintf("vlan_%v_%v", svc.Role, d.Get("vlan_id").(int)), svc.Payload)
		}

		_, err = iosxe.MultiSession(svc)
		if err != nil {
			return diag.FromErr(err)
		}
		d.SetId(fmt.Sprintf("vlan_%v", d.Get("vlan_id").(int)))
	}

	return diags
}

func resourceCiscoNativeVlanDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Cisco VLAN DELETE")
	var diags diag.Diagnostics
	var err error

	c, _ := meta.(*providerClient)
	svc := &service.Client{
		Method:   "DELETE",
		Path:     fmt.Sprintf("/data/Cisco-IOS-XE-native:native/vlan/configuration-entry=%v", d.Get("vlan_id").(int)),
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

func (*providerClient) resourceCiscoNativeVlanData(d *schema.ResourceData) *vlan.CiscoIOSXENativeVlans {
	data := &vlan.CiscoIOSXENativeVlan{}
	vlanCfg := &vlan.CiscoIOSXEVlanConfigurationEntry{}

	vlanCfg.VlanID = fmt.Sprintf("%v", d.Get("vlan_id").(int))

	if v, ok := d.GetOk("evpn_instance"); ok {
		member := &vlan.CiscoIOSXEVlanConfigurationEntryMember{}
		evpn := &vlan.CiscoIOSXEVlanConfigurationEntryEvpnInstance{}
		evpn.EvpnInstance = v.(int)
		evpn.Vni = d.Get("vni").(int)
		member.EvpnInstance = evpn
		vlanCfg.Member = *member
	} else {
		member := &vlan.CiscoIOSXEVlanConfigurationEntryMember{}
		member.Vni = d.Get("vni").(int)
		vlanCfg.Member = *member
	}

	vlanList := &vlan.CiscoIOSXEVlanVlanList{}
	vlanList.ID = d.Get("vlan_id").(int)
	if d.Get("name").(string) == "ManagedByTerraform" {
		vlanList.Name = fmt.Sprintf("%v_%v", d.Get("name").(string), d.Get("vlan_id").(int))
	} else {
		vlanList.Name = d.Get("name").(string)
	}

	data.CiscoIOSXEVlanConfigurationEntry = append(data.CiscoIOSXEVlanConfigurationEntry, *vlanCfg)
	data.CiscoIOSXEVlanVlanList = append(data.CiscoIOSXEVlanVlanList, *vlanList)
	return &vlan.CiscoIOSXENativeVlans{*data}
}
