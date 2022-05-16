package provider

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/robertcsapo/terraform-provider-ciscoevpn/internal/model/loopback"
	"github.com/robertcsapo/terraform-provider-ciscoevpn/internal/provider/iosxe"
	"github.com/robertcsapo/terraform-provider-ciscoevpn/internal/provider/service"
)

func resourceCiscoNativeLoopbackInterface() *schema.Resource {
	return &schema.Resource{
		Description:   "Cisco EVPN Instance",
		CreateContext: resourceCiscoNativeLoopbackInterfaceCreate,
		ReadContext:   resourceCiscoNativeLoopbackInterfaceRead,
		UpdateContext: resourceCiscoNativeLoopbackInterfaceUpdate,
		DeleteContext: resourceCiscoNativeLoopbackInterfaceDelete,
		Schema: map[string]*schema.Schema{
			"host": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"loopback_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Default:  "Managed by Terraform (ciscoevpn)",
				Optional: true,
			},
			"ipv4_address": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.IsIPAddress,
			},
			"ipv4_mask": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.IsIPAddress,
			},
			"pim_sm": {
				Type:     schema.TypeBool,
				Default:  true,
				Optional: true,
			},
			"interface_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Name of the interface`,
			},
		},
	}
}

func resourceCiscoNativeLoopbackInterfaceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Cisco Loopback CREATE")
	var diags diag.Diagnostics
	var err error

	c, _ := meta.(*providerClient)
	svc := &service.Client{
		Method:   "PATCH",
		Path:     "/data/Cisco-IOS-XE-native:native/interface/Loopback",
		Provider: c.Provider,
		Device:   d.Get("host").(string),
	}

	data := c.resourceCiscoNativeLoopbackInterfaceData(d)
	if data == nil {
		err = errors.New("no data in yang model")
		return diag.FromErr(err)
	}

	if b, err := json.MarshalIndent(data, "", "\t"); err == nil {
		svc.Payload = string(b)
	}
	if svc.Provider.Get("debug").(bool) {
		debugJson(fmt.Sprintf("loopback_interface_%v", svc.Device), svc.Payload)
	}

	_, err = iosxe.SingleSession(svc)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("interface_name", fmt.Sprintf("Loopback%v", d.Get("loopback_id")))
	d.SetId(fmt.Sprintf("%v", d.Get("loopback_id").(int)))
	return diags
}

func resourceCiscoNativeLoopbackInterfaceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Cisco Loopback READ") // TODO
	var diags diag.Diagnostics
	return diags
}

func resourceCiscoNativeLoopbackInterfaceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Cisco Loopback UPDATE")
	var diags diag.Diagnostics
	var err error

	if d.HasChange("loopback_id") {
		oldState, _ := d.GetChange("loopback_id")
		d.Set("loopback_id", oldState)
		return diag.Errorf("Not supported to change Loopback number")
	}

	c, _ := meta.(*providerClient)
	svc := &service.Client{
		Method:   "PATCH",
		Path:     "/data/Cisco-IOS-XE-native:native/interface/Loopback",
		Provider: c.Provider,
		Device:   d.Get("host").(string),
	}

	data := c.resourceCiscoNativeLoopbackInterfaceData(d)
	if data == nil {
		err = errors.New("no data in yang model")
		return diag.FromErr(err)
	}

	if b, err := json.MarshalIndent(data, "", "\t"); err == nil {
		svc.Payload = string(b)
	}
	if svc.Provider.Get("debug").(bool) {
		debugJson(fmt.Sprintf("loopback_interface_%v", svc.Device), svc.Payload)
	}

	_, err = iosxe.SingleSession(svc)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("interface_name", fmt.Sprintf("Loopback%v", d.Get("loopback_id")))
	d.SetId(fmt.Sprintf("%v", d.Get("loopback_id").(int)))
	return diags
}

func resourceCiscoNativeLoopbackInterfaceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Cisco Loopback DELETE")
	var diags diag.Diagnostics
	var err error
	c, _ := meta.(*providerClient)
	svc := &service.Client{
		Method:   "DELETE",
		Path:     fmt.Sprintf("/data/Cisco-IOS-XE-native:native/interface/Loopback=%v", d.Get("loopback_id").(int)),
		Provider: c.Provider,
		Device:   d.Get("host").(string),
	}

	_, err = iosxe.SingleSession(svc)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return diags
}

func (*providerClient) resourceCiscoNativeLoopbackInterfaceData(d *schema.ResourceData) *loopback.CiscoIOSXENativeLoopbackInterface {
	data := &loopback.CiscoIOSXENativeLoopbackInterface{}
	lp := &loopback.CiscoIOSXENativeLoopback{
		Name: d.Get("loopback_id").(int),
	}
	lp.IP.Address.Primary.Address = d.Get("ipv4_address").(string)
	lp.IP.Address.Primary.Mask = d.Get("ipv4_mask").(string)
	if d.Get("pim_sm").(bool) {
		lp.IP.Pim.CiscoIOSXEMulticastPimModeChoiceCfg.SparseMode = map[string]string{}
	}
	lp.Description = d.Get("description").(string)
	data.CiscoIOSXENativeLoopback = append(data.CiscoIOSXENativeLoopback, *lp)
	return data
}
