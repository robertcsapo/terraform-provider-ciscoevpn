package iosxe

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/CiscoDevNet/iosxe-go-client/client"
	"github.com/CiscoDevNet/iosxe-go-client/container"
	"github.com/CiscoDevNet/iosxe-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/robertcsapo/terraform-provider-ciscoevpn/internal/provider/service"
)

type sessionClient struct {
	Client  *client.V2
	Host    string
	Service *service.Client
}

func SingleSession(svc *service.Client) (string, error) {
	var err error
	var body string
	s := &sessionClient{
		Host:    fmt.Sprintf("%v", svc.Device),
		Service: svc,
	}
	body, err = s.methods(svc.Method)
	if err != nil {
		log.Println("[DEBUG] ERROR SingleSession: ", err)
		return body, err
	}
	return body, nil
}

func MultiSession(svc *service.Client) (map[string]string, error) {
	var err error
	var body string
	var data = make(map[string]string)
	hosts := HostRoles(svc.Devices, svc.Role)
	if len(hosts) == 0 {
		log.Panicln("[PANIC] No hosts found, using role: ", svc.Role)
	}
	for _, host := range hosts {
		s := &sessionClient{
			Host:    fmt.Sprintf("%v", host),
			Service: svc,
		}
		body, err = s.methods(svc.Method)
		if err != nil {
			log.Println("[DEBUG] ERROR MultiSession: ", err)
			return data, err
		}
		data[s.Host] = body

	}
	return data, nil
}

func HostRoles(devices []interface{}, role string) []interface{} {
	var hosts []interface{}
	for _, rolesMapRaw := range devices {
		rolesMap, _ := rolesMapRaw.(map[string]interface{})
		for k, roles := range rolesMap {
			if k == role {
				for _, v := range roles.(*schema.Set).List() {
					for _, host := range v.(map[string]interface{}) {
						hosts = append(hosts, host.([]interface{})...)
					}
				}
			}
		}
	}
	return hosts
}

func NewClient(host string, d schema.ResourceData) (*client.V2, diag.Diagnostics) {
	var diags diag.Diagnostics
	iosxeV2Client, err := client.NewV2(
		host,
		d.Get("username").(string),
		d.Get("password").(string),
		d.Get("timeout").(int),
		d.Get("insecure").(bool),
		d.Get("proxy_url").(string),
		d.Get("proxy_creds").(string),
		d.Get("ca_file").(string),
	)

	if err != nil {
		log.Printf("[DEBUG] ERROR: %v\n", err)
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create Cisco IOS-XE client",
			Detail:   fmt.Sprintf("Unable to create Cisco IOSXE client. Error - %v", err),
		})
		return nil, diags
	}
	return iosxeV2Client, diags

}

func (s *sessionClient) methods(method string) (string, error) {
	var err error
	var body string
	var resp *http.Response
	var container *container.Container
	var httpRetry int
	var httpSleep time.Duration

	httpRetry = 20
	httpSleep = 10 * time.Second

	iosxeGM := &models.GenericModel{
		JSONPayload: s.Service.Payload,
	}

	c, _ := NewClient(
		fmt.Sprintf("https://%v", s.Host),
		s.Service.Provider,
	)

	switch method {
	case "GET":
		log.Println("[DEBUG] IOS-XE GET on: ", s.Host)
		for i := 0; ; i++ {
			resp, container, err = c.Get(s.Service.Path, nil)
			if resp != nil {
				if resp.StatusCode == 409 {
					log.Println(s.httpErrorMsg(resp.StatusCode))
				} else {
					break
				}
			}
			if i >= (httpRetry - 1) {
				break
			}
			log.Printf("[DEBUG] IOS-XE Retry: (%v/%v) waiting: %v\n", i, httpRetry, httpSleep)
			time.Sleep(httpSleep)
		}
		if err != nil {
			log.Println("[DEBUG] ERROR GET: ", err, resp.Status)
			return body, err
		}
		body = container.String()
	case "PATCH":
		log.Println("[DEBUG] IOS-XE PATCH on: ", s.Host)
		for i := 0; ; i++ {
			resp, _, err = c.PatchRaw(s.Service.Path, s.Service.Payload)
			if resp != nil {
				if resp.StatusCode == 409 {
					log.Println(s.httpErrorMsg(resp.StatusCode))
				} else {
					break
				}
			}
			if i >= (httpRetry - 1) {
				break
			}
			log.Printf("[DEBUG] IOS-XE Retry: (%v/%v) waiting: %v\n", i, httpRetry, httpSleep)
			time.Sleep(httpSleep)
		}
		if err != nil {
			log.Println("[DEBUG] ERROR PATCH: ", err, resp.Status)
			return body, err
		}
	case "UPDATE":
		log.Println("[DEBUG] IOS-XE UPDATE on: ", s.Host)
		for i := 0; ; i++ {
			resp, err = c.Update(s.Service.Path, iosxeGM)
			if resp != nil {
				if resp.StatusCode == 409 {
					log.Println(s.httpErrorMsg(resp.StatusCode))
				} else {
					break
				}
			}
			if i >= (httpRetry - 1) {
				break
			}
			log.Printf("IOS-XE Retry: (%v/%v) waiting: %v\n", i, httpRetry, httpSleep)
			time.Sleep(httpSleep)
		}
		if err != nil {
			log.Println("[DEBUG] ERROR UPDATE: ", err, resp.Status)
			return body, err
		}
	case "DELETE":
		log.Printf("[DEBUG] IOS-XE DELETE on: %v %v\n", s.Host, s.Service.Path)
		for i := 0; ; i++ {
			resp, err = c.Delete(s.Service.Path)
			// TODO
			// catch this?
			// panic: runtime error: invalid memory address or nil pointer dereference
			if resp != nil {
				if resp.StatusCode == 409 {
					log.Println(s.httpErrorMsg(resp.StatusCode))
				} else {
					break
				}
			}
			if i >= (httpRetry - 1) {
				break
			}
			log.Printf("[DEBUG] IOS-XE Retry: (%v/%v) waiting: %v\n", i, httpRetry, httpSleep)
			time.Sleep(httpSleep)
		}
		if err != nil {
			log.Println("[DEBUG] ERROR DELETE: ", err, resp.Status)
			return body, err
		}
	}
	return body, nil
}

func (*sessionClient) httpErrorMsg(status int) string {
	var msg string
	if status == 409 {
		msg = "[DEBUG] IOS-XE configuration database is unavailable"
	}
	return msg
}
