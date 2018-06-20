package main

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"io/ioutil"
	"net/http"
	"os/exec"
	"regexp"
)

// One of dozens of URLs that return the request's IP address as the only content
// Others that seems likely to have decent longevity:
// - whatismyip.akamai.com (but SSL cert is wrong),
// - checkip.dyndns.org (contains other content; SSL hangs)
const defaultHttpUrl = "https://checkip.amazonaws.com"

// Match an IPv4 address in "quad-dotted" format (four decimal octets separated by periods)
func extractIPAddress(b []byte) ([]byte, error) {
	re, err := regexp.Compile("\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}")
	if err != nil {
		return nil, err
	}
	return re.Find(b), nil
}

func dnsResolver() ([]byte, error) {
	return exec.Command("dig", "+short", "myip.opendns.com", "@resolver1.opendns.com").Output()
}

func httpResolver(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

func resourceRead(d *schema.ResourceData, meta interface{}) error {
	// GetOkExists is required since the 'zero value' for bool is false
	dns, dnsOk := d.GetOkExists("dns")
	http, httpOk := d.GetOk("http")
	// Call out to one of dnsResolver or httpResolver based on the given configuration
	var resolution []byte
	var err error
	var resolver string

	if httpOk && http.(string) != "" {
		// resolve via HTTP if http = "<something>"
		// if http is explicitly set to something, we ignore dns;
		// it wouldn't make sense to set dns = true, and it might be worth a warning,
		// but we'll just silently go with http in that case
		resolution, err = httpResolver(http.(string))
		resolver = "http"
	} else if dnsOk && dns.(bool) == false {
		// also resolve via HTTP if dns = false (in which case use defaultHttpUrl)
		resolution, err = httpResolver(defaultHttpUrl)
		resolver = "http"
	} else {
		// resolve via DNS if <nothing is configured>, or if dns = true
		resolution, err = dnsResolver()
		resolver = "dns"
	}
	// catch any of those calls returning an error
	if err != nil {
		return err
	}

	ipAddress, err := extractIPAddress(resolution)
	if err != nil {
		return err
	}
	if ipAddress == nil {
		return fmt.Errorf("Failed to match IPv4 address in '%s' resolver's response: %q", resolver, resolution)
	}
	// We must set the Id; otherwise other resources can't access anything set in this method
	d.SetId("_network_info-" + resolver)
	d.Set("wan_ip_address", string(ipAddress))
	return nil
}

func resource() *schema.Resource {
	return &schema.Resource{
		Read: resourceRead,
		Schema: map[string]*schema.Schema{
			"dns": {
				Type:        schema.TypeBool,
				Description: "To resolve your IP Address with (Open)DNS, omit or set this to true",
				Optional:    true,
			},
			"http": {
				Type:        schema.TypeString,
				Description: "To resolve your IP Address with HTTP, set this to a URL that will return your IP(v4) Address",
				Optional:    true,
			},
			"wan_ip_address": {
				Type:        schema.TypeString,
				Description: "IP Address of calling machine as seen on the current wide area network (WAN)",
				Computed:    true,
			},
		},
	}
}

func Provider() *schema.Provider {
	return &schema.Provider{
		// DataSourcesMap is like ResourcesMap but points to Resources that only have Read and Schema implemented
		DataSourcesMap: map[string]*schema.Resource{
			"network_info": resource(),
		},
	}
}
