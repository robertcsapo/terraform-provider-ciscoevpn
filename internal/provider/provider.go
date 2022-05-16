package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

type providerClient struct {
	Provider schema.ResourceData
	Devices  *schema.Set
}

func init() {
	// Set descriptions to support markdown syntax, this will be used in document generation
	// and the language server.
	schema.DescriptionKind = schema.StringMarkdown
}

func Provider() *schema.Provider {
	p := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"username": {
				Type:         schema.TypeString,
				Required:     true,
				Sensitive:    true,
				DefaultFunc:  schema.EnvDefaultFunc("EVPN_USERNAME", nil),
				ValidateFunc: validation.StringIsNotWhiteSpace,
				Description:  `The Username of the IOSXE switch. E.g.: "admin". This can also be set by environment variable "EVPN_USERNAME".`,
			},
			"password": {
				Type:         schema.TypeString,
				Required:     true,
				Sensitive:    true,
				DefaultFunc:  schema.EnvDefaultFunc("EVPN_PASSWORD", nil),
				ValidateFunc: validation.StringIsNotWhiteSpace,
				Description:  `The Password of the IOSXE switch. E.g.: "somePassword". This can also be set by environment variable "EVPN_PASSWORD".`,
			},
			"insecure": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Allow insecure TLS. Default: true, means the API call is insecure.",
			},
			"timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     30,
				Description: "Timeout for HTTP requests. Default value: 30.",
			},
			"ca_file": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("EVPN_CA_FILE", nil),
				Description: "The path to CA certificate file (PEM). In case, certificate is based on legacy CN instead of ASN, set env. variable `GODEBUG=x509ignoreCN=0`. This can also be set by environment variable `EVPN_CA_FILE`.",
			},
			"proxy_url": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("EVPN_PROXY_URL", nil),
				Description: "Proxy Server URL with port number. This can also be set by environment variable `EVPN_PROXY_URL`.",
			},
			"proxy_creds": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("EVPN_PROXY_CREDS", nil),
				Description: "Proxy credential in format `username:password`. This can also be set by environment variable `EVPN_PROXY_CREDS`.",
			},
			"debug": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Debug JSON Payloads in to debug folder",
			},
			"roles": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"spines": {
							Type:     schema.TypeSet,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"iosxe": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},
						"leafs": {
							Type:     schema.TypeSet,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"iosxe": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},
						"borders": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"iosxe": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},
					},
				},
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"ciscoevpn_loopback":                 resourceCiscoNativeLoopbackInterface(),
			"ciscoevpn_bgp_system":               resourceCiscoNativeBgpSystem(),
			"ciscoevpn_bgp_neighbor":             resourceCiscoNativeBgpNeighbor(),
			"ciscoevpn_bgp_neighbor_vrf_unicast": resourceCiscoNativeBgpNeighborVrfUnicast(),
			"ciscoevpn_bgp_vrf":                  resourceCiscoNativeBgpVrf(),
			"ciscoevpn_evpn":                     resourceCiscoNativeL2VpnEvpn(),
			"ciscoevpn_evpn_instance":            resourceCiscoNativeEvpnInstance(),
			"ciscoevpn_vrf":                      resourceCiscoNativeVrf(),
			"ciscoevpn_vlan":                     resourceCiscoNativeVlan(),
			"ciscoevpn_nve":                      resourceCiscoNativeNve(),
			"ciscoevpn_svi":                      resourceCiscoNativeSvi(),
			"ciscoevpn_subinterface":             resourceCiscoNativeSubInterface(),
			"ciscoevpn_dhcp":             resourceCiscoNativeDhcp(),
			"ciscoevpn_dhcp_helper":             resourceCiscoNativeDhcpHelper(),
		},
		DataSourcesMap: map[string]*schema.Resource{},
	}
	p.ConfigureContextFunc = providerConfigure(p)
	return p
}

func providerConfigure(p *schema.Provider) func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	return func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		var diags diag.Diagnostics
		if diags = providerValidateInput(d); diags.HasError() {
			return nil, diags
		}
		return &providerClient{
			Provider: *d,
			Devices:  d.Get("roles").(*schema.Set),
		}, diags
	}
}

func providerValidateInput(d *schema.ResourceData) diag.Diagnostics {
	var diags diag.Diagnostics
	if v, ok := d.GetOk("username"); !ok && v == "" {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Username is required",
			Detail:   "Username must be set for Cisco EVPN Provider",
		})
	}
	if v, ok := d.GetOk("password"); !ok && v == "" {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Password is required",
			Detail:   "Password must be set for Cisco EVPN Provider",
		})
	}
	if v, ok := d.GetOk("roles"); ok && v.(*schema.Set).Len() > 0 {
		for _, rolesMapRaw := range v.(*schema.Set).List() {
			rolesMap, _ := rolesMapRaw.(map[string]interface{})
			for k, roles := range rolesMap {
				if len(roles.(*schema.Set).List()) == 0 {
					// Borders is optional
					if k != "borders" {
						diags = append(diags, diag.Diagnostic{
							Severity: diag.Error,
							Summary:  fmt.Sprintf("%v roles is required", k),
							Detail:   fmt.Sprintf("Cisco devices (hosts) required to added as %v for Cisco EVPN Provider", k),
						})
					}
				}
			}
		}
	}
	return diags
}
