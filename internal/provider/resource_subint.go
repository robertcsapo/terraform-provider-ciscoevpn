package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/robertcsapo/terraform-provider-ciscoevpn/internal/model/subinterface"
	"github.com/robertcsapo/terraform-provider-ciscoevpn/internal/provider/iosxe"
	"github.com/robertcsapo/terraform-provider-ciscoevpn/internal/provider/service"
)

func resourceCiscoNativeSubInterface() *schema.Resource {
	return &schema.Resource{
		Description:   "Cisco Sub Interface",
		CreateContext: resourceCiscoNativeSubInterfaceCreate,
		ReadContext:   resourceCiscoNativeSubInterfaceRead,
		UpdateContext: resourceCiscoNativeSubInterfaceUpdate,
		DeleteContext: resourceCiscoNativeSubInterfaceDelete,
		Schema: map[string]*schema.Schema{
			"host": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"description": {
				Type:     schema.TypeString,
				Default:  "Managed by Terraform (ciscoevpn)",
				Optional: true,
			},
			"ethernet": {
				Type:     schema.TypeString,
				Required: true,
			},
			"interface_speed": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"dot1q": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"vrf": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ipv4_address": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.IsIPv4Address,
			},
			"ipv4_mask": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.IsIPv4Address,
			},
			"ipv4_remote": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.IsIPv4Address,
			},
		},
	}
}

func resourceCiscoNativeSubInterfaceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Cisco Sub Interface CREATE")
	var diags diag.Diagnostics
	var err error
	var uri string

	c, _ := meta.(*providerClient)
	svc := &service.Client{
		Method:   "PATCH",
		Provider: c.Provider,
		Device:   d.Get("host").(string),
	}

	data, uri := c.resourceCiscoNativeSubInterfaceData(d)
	svc.Path = fmt.Sprintf("/data/Cisco-IOS-XE-native:native/interface/%v", uri)
	if data == nil {
		return diag.Errorf("No data in yang model")
	}

	if b, err := json.MarshalIndent(data, "", "\t"); err == nil {
		svc.Payload = string(b)
	}

	if svc.Provider.Get("debug").(bool) {
		debugJson(fmt.Sprintf("subint_%v", d.Get("ipv4_address").(string)), svc.Payload)
	}

	_, err = iosxe.SingleSession(svc)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%v/%v", uri, d.Get("ethernet").(string)))
	return diags
}

func resourceCiscoNativeSubInterfaceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Println("[DEBUG] TODO")
	var diags diag.Diagnostics
	return diags
}

func resourceCiscoNativeSubInterfaceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Cisco Sub Interface UPDATE")
	var diags diag.Diagnostics
	var err error

	if d.HasChange("ethernet") {
		oldState, _ := d.GetChange("ethernet")
		d.Set("ethernet", oldState)
		return diag.Errorf("Not supported to change Ethernet slots of Sub Interface")
	}
	if d.HasChange("interface_speed") {
		oldState, _ := d.GetChange("interface_speed")
		d.Set("interface_speed", oldState)
		return diag.Errorf("Not supported to change Interface Speed of Sub Interface")
	}
	if d.HasChange("roles") {
		oldState, _ := d.GetChange("roles")
		d.Set("roles", oldState)
		return diag.Errorf("Not supported to change Roles")
	}

	c, _ := meta.(*providerClient)
	svc := &service.Client{
		Method:   "PATCH",
		Provider: c.Provider,
		Device:   d.Get("host").(string),
	}

	data, uri := c.resourceCiscoNativeSubInterfaceData(d)
	svc.Path = fmt.Sprintf("/data/Cisco-IOS-XE-native:native/interface/%v", uri)
	if data == nil {
		return diag.Errorf("No data in yang model")
	}

	if b, err := json.MarshalIndent(data, "", "\t"); err == nil {
		svc.Payload = string(b)
	}

	if svc.Provider.Get("debug").(bool) {
		debugJson(fmt.Sprintf("subint_%v", d.Get("ipv4_address").(string)), svc.Payload)
	}

	_, err = iosxe.SingleSession(svc)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%v/%v", uri, d.Get("ethernet").(string)))
	return diags
}

func resourceCiscoNativeSubInterfaceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Cisco Sub Interface DELETE")
	var err error
	var diags diag.Diagnostics
	var uri string

	c, _ := meta.(*providerClient)
	svc := &service.Client{
		Method:   "DELETE",
		Provider: c.Provider,
		Device:   d.Get("host").(string),
	}

	id := d.Id()
	ethernet := strings.Split(id, "/")
	slots := strings.Split(id, ethernet[0])
	slot := slots[1]
	slot = strings.ReplaceAll(slot[1:], "/", "%2F")
	uri = fmt.Sprintf("%v=%v", ethernet[0], slot)
	svc.Path = fmt.Sprintf("/data/Cisco-IOS-XE-native:native/interface/%v", uri)

	_, err = iosxe.SingleSession(svc)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return diags
}

func (*providerClient) resourceCiscoNativeSubInterfaceData(d *schema.ResourceData) (*subinterface.CiscoIOSXENativeEthernet, string) {
	var uri string
	data := &subinterface.CiscoIOSXENativeEthernet{}

	ethernet := &subinterface.CiscoIOSXENativeEthernetInterface{}
	ethernet.Name = d.Get("ethernet").(string)
	ethernet.Description = d.Get("description").(string)
	ethernet.IP.Address.Primary.Address = d.Get("ipv4_address").(string)
	ethernet.IP.Address.Primary.Mask = d.Get("ipv4_mask").(string)

	if v, ok := d.GetOk("vrf"); ok {
		ethernet.Vrf.Forwarding = v.(string)
	}
	if v, ok := d.GetOk("dot1q"); ok {
		ethernet.Encapsulation.Dot1Q.VlanID = v.(int)
	}

	// Handle IOS-XE Interface naming
	switch d.Get("interface_speed").(int) {
	case 10:
		data.Ten = append(data.Ten, *ethernet)
		uri = "TenGigabitEthernet"
	case 25:
		data.TwentyFive = append(data.TwentyFive, *ethernet)
		uri = "TwentyFiveGigE"
	case 40:
		data.Forty = append(data.Forty, *ethernet)
		uri = "FortyGigabitEthernet"
	case 100:
		data.Hundred = append(data.Hundred, *ethernet)
		uri = "HundredGigE"
	case 400:
		data.FourHundred = append(data.FourHundred, *ethernet)
		uri = "FourHundredGigE"
	default:
		data.One = append(data.One, *ethernet)
		uri = "GigabitEthernet"
	}

	return data, uri

}
