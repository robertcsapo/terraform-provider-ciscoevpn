package service

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

type Client struct {
	Method   string
	Path     string
	Payload  string
	Provider schema.ResourceData
	Device   string
	Devices  []interface{}
	Role     string
}
