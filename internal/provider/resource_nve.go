package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/robertcsapo/terraform-provider-ciscoevpn/internal/model/nve"
	"github.com/robertcsapo/terraform-provider-ciscoevpn/internal/provider/iosxe"
	"github.com/robertcsapo/terraform-provider-ciscoevpn/internal/provider/service"
)

func resourceCiscoNativeNve() *schema.Resource {
	return &schema.Resource{
		Description:   "Cisco NVE",
		CreateContext: resourceCiscoNativeNveCreate,
		ReadContext:   resourceCiscoNativeNveRead,
		UpdateContext: resourceCiscoNativeNveUpdate,
		DeleteContext: resourceCiscoNativeNveDelete,
		Schema: map[string]*schema.Schema{
			"roles": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"description": {
				Type:     schema.TypeString,
				Default:  "Managed by Terraform (ciscoevpn)",
				Optional: true,
			},
			"source_interface": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"vni": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"vni_ipv4_multicast_group": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"vni_ingress_replication": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
		},
	}
}

func resourceCiscoNativeNveCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Cisco NVE CREATE")
	var diags diag.Diagnostics
	var err error

	c, _ := meta.(*providerClient)
	svc := &service.Client{
		Method:   "PATCH",
		Path:     "/data/Cisco-IOS-XE-native:native/interface/nve",
		Provider: c.Provider,
		Devices:  c.Devices.List(),
	}

	id := loopbackId(d.Get("source_interface").(string))
	d.Set("source_interface", id)

	roles := d.Get("roles").([]interface{})
	for _, role := range roles {
		svc.Role = role.(string)
		data := c.resourceCiscoNativeNveData(d)
		if data == nil {
			return diag.Errorf("No data in yang model")
		}

		if b, err := json.MarshalIndent(data, "", "\t"); err == nil {
			svc.Payload = string(b)
		}
		if svc.Provider.Get("debug").(bool) {
			debugJson(fmt.Sprintf("nve_%v", svc.Role), svc.Payload)
		}

		_, err = iosxe.MultiSession(svc)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	err = d.Set("source_interface", fmt.Sprintf("Loopback%v", id))
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(fmt.Sprintf("nve_%v", svc.Role))
	return diags
}

func resourceCiscoNativeNveRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Cisco NVE READ") // TODO
	var diags diag.Diagnostics
	return diags
}

func resourceCiscoNativeNveUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Cisco NVE UPDATE")
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
		Path:     "/data/Cisco-IOS-XE-native:native/interface/nve=1",
		Provider: c.Provider,
		Devices:  c.Devices.List(),
	}

	id := loopbackId(d.Get("source_interface").(string))
	d.Set("source_interface", id)

	roles := d.Get("roles").([]interface{})
	for _, role := range roles {
		svc.Role = role.(string)
		data := c.resourceCiscoNativeNveData(d)
		if data == nil {
			return diag.Errorf("No data in yang model")
		}

		if b, err := json.MarshalIndent(data, "", "\t"); err == nil {
			svc.Payload = string(b)
		}
		if svc.Provider.Get("debug").(bool) {
			debugJson(fmt.Sprintf("nve_%v", svc.Role), svc.Payload)
		}

		_, err = iosxe.MultiSession(svc)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	err = d.Set("source_interface", fmt.Sprintf("Loopback%v", id))
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(fmt.Sprintf("nve_%v", svc.Role))
	return diags
}

func resourceCiscoNativeNveDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Cisco NVE DELETE")
	var diags diag.Diagnostics
	var err error

	c, _ := meta.(*providerClient)
	svc := &service.Client{
		Method:   "DELETE",
		Path:     "/data/Cisco-IOS-XE-native:native/interface/nve=1",
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

func (svc *providerClient) resourceCiscoNativeNveData(d *schema.ResourceData) *nve.CiscoIOSXENativeNves {
	data := &nve.CiscoIOSXENativeNves{}
	nveData := &nve.CiscoIOSXENativeNve{}

	nveData.Name = 1
	nveData.Description = d.Get("description").(string)
	nveData.HostReachability.Protocol.Bgp = null()
	id, err := strconv.Atoi(d.Get("source_interface").(string))
	if err != nil {
		log.Panicln("[PANIC] Can't find Loopback ID", err)
	}
	nveData.SourceInterface.Loopback = id

	if v, ok := d.GetOk("vni"); ok {
		for vrf, vnis := range v.(map[string]interface{}) {
			vni := &nve.CiscoIOSXENativeNveVni{}

			vni.VniRange = svc.vniRanges(vnis)
			vni.Vrf = vrf

			nveData.MemberInOneLine.Member.Vni = append(nveData.MemberInOneLine.Member.Vni, *vni)
		}
	}
	if v, ok := d.GetOk("vni_ipv4_multicast_group"); ok {

		for mc, vnis := range v.(map[string]interface{}) {
			vni := &nve.CiscoIOSXENativeNveVni{}
			vniMulticast := &nve.CiscoIOSXENativeNveMcastGroup{}

			vni.VniRange = svc.vniRanges(vnis)
			vniMulticast.MulticastGroupMin = fmt.Sprintf("%v", mc)

			vni.McastGroup = vniMulticast
			nveData.Member.Vni = append(nveData.Member.Vni, *vni)
		}
	}
	if v, ok := d.GetOk("vni_ingress_replication"); ok {

		for _, id := range v.([]interface{}) {
			vni := &nve.CiscoIOSXENativeNveVni{}
			ingress := &nve.CiscoIOSXENativeNveIrCpConfig{}

			vni.VniRange = fmt.Sprintf("%v", id)
			ingress.IngressReplication = null()

			vni.IrCpConfig = ingress
			nveData.Member.Vni = append(nveData.Member.Vni, *vni)
		}
	}

	data.CiscoIOSXENativeNve = append(data.CiscoIOSXENativeNve, *nveData)
	return data
}

func (*providerClient) vniRanges(data interface{}) string {
	vni := fmt.Sprintf("%v", data)
	vniRange := strings.Split(vni, "-")
	if len(vniRange) == 1 {
		return fmt.Sprintf("%v", vniRange[0])
	}
	return fmt.Sprintf("%v-%v", vniRange[0], vniRange[1])

}
