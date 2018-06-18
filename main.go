package main

import (
	"bytes"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/plugin"
	"github.com/hashicorp/terraform/terraform"
	"os/exec"
)

func resolveIPAddress() (ipAddress string, err error) {
	stdoutBytes, err := exec.Command("dig", "+short", "myip.opendns.com", "@resolver1.opendns.com").Output()
	if err != nil {
		return
	}
	ipAddress = string(bytes.TrimSpace(stdoutBytes))
	return
}

func resource() *schema.Resource {
	return &schema.Resource{
		Read: resourceRead,
		Schema: map[string]*schema.Schema{
			"wan_ip_address": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceRead(d *schema.ResourceData, meta interface{}) error {
	ipAddress, err := resolveIPAddress()
	if err != nil {
		return err
	}
	// you must set the Id; otherwise other resources can't access anything set in this method
	d.SetId("_network_info")
	d.Set("wan_ip_address", ipAddress)
	return nil
}

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		// DataSourcesMap is like ResourcesMap but points to Resources that only have Read and Schema implemented
		DataSourcesMap: map[string]*schema.Resource{
			"network_info": resource(),
		},
	}
}

func main() {
	opts := plugin.ServeOpts{ProviderFunc: Provider}
	plugin.Serve(&opts)
}
