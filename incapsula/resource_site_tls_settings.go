package incapsula

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"log"
	"strconv"
)

const mandatoryDefault = false
const isPortsExceptionDefault = false
const isHostsExceptionDefault = false
const forwardToOriginDefault = false
const headerNameDefault = "clientCertificateInfo"
const headerValueDefault = "FULL_CERT"
const isDisableSessionResumptionDefault = false

func resourceSiteTlsSetings() *schema.Resource {
	return &schema.Resource{
		Create: resourceSiteTlsSetingsUpdate,
		Read:   resourceSiteTlsSetingsRead,
		Update: resourceSiteTlsSetingsUpdate,
		Delete: resourceSiteTlsSetingsDelete,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				d.Set("site_id", d.Id())
				log.Printf("[DEBUG] Importing Site TLS Settings for Site ID %s", d.Id())
				return []*schema.ResourceData{d}, nil
			},
		},
		Schema: map[string]*schema.Schema{
			// Required Arguments
			"site_id": {
				Description: "Numeric identifier of the site to operate on.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"mandatory": {
				Description: "When set to true, the end user is required to present the client certificate in order to access the site. Default - false.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     mandatoryDefault,
			},
			"ports": {
				Description: "The ports on which client certificate authentication is supported. If left empty, client certificates are supported on all ports. Default: empty list",
				Type:        schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Optional: true,
			},
			"is_ports_exception": {
				Description: "When set to true, client certificates are not supported on the ports listed in the Ports field ('blacklisted'). Default - false.",
				Type:        schema.TypeBool,
				Default:     isPortsExceptionDefault,
				Optional:    true,
			},
			"hosts": {
				Description: "The hosts on which client certificate authentication is supported. If left empty, client certificates are supported on all hosts. Default: empty list.",
				Optional:    true,
				Type:        schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"is_hosts_exception": {
				Description: "When set to true, client certificates are not supported on the hosts listed in the Hosts field ('blacklisted'). Default - false.",
				Type:        schema.TypeBool,
				Default:     isHostsExceptionDefault,
				Optional:    true,
			},
			"fingerprints": {
				Description: "Permitted client certificate fingerprints. If left empty, all fingerprints are permitted. Default - empty list.",
				Type:        schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
			},
			"forward_to_origin": {
				Description: "When set to true, the contents specified in headerValue are sent to the origin server in the header specified by headerName. Default - false. If parameter is set to true, specify of `header_name`, `header_value` are required.", //todo KATRIN change
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     forwardToOriginDefault,
			},
			"header_name": {
				Description: "The name of the header to send header content in. By default, the header name is 'clientCertificateInfo'. Specifying this parameter is relevant only if `forward_to_origin` is set to true.", //todo KATRIN change
				Type:        schema.TypeString,
				Default:     headerNameDefault,
				Optional:    true,
			},
			"header_value": {
				Description:  "The content to send in the header specified by headerName. One of the following: FULL_CERT (for full certificate in Base64) COMMON_NAME (for certificate's common name (CN)) FINGERPRINT (for the certificate fingerprints in SHA1) SERIAL_NUMBER (for the certificate's serial number). Specifying this parameter is relevant only if `forward_to_origin` is set to true.", //todo KATRIN change
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"FULL_CERT", "COMMON_NAME", "FINGERPRINT", "SERIAL_NUMBER"}, false),
				Optional:     true,
				Default:      headerValueDefault,
			},
			"is_disable_session_resumption": {
				Description: "", //todo KATRIN change
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     isDisableSessionResumptionDefault,
			},
		},
	}
}

func resourceSiteTlsSetingsUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)

	siteIDStr := d.Get("site_id").(string)
	siteID, err := strconv.Atoi(siteIDStr)
	//todo katrin  edit error
	if err != nil {
		return fmt.Errorf("failed to convert Site Id for Incapsula Site TLS Settings resource, actual value: %s, expected numeric id", siteIDStr)
	}

	ports := []int{}
	for _, entry := range d.Get("ports").(*schema.Set).List() {
		ports = append(ports, entry.(int))
	}

	hosts := []string{}
	for _, entry := range d.Get("hosts").(*schema.Set).List() {
		hosts = append(hosts, entry.(string))
	}

	fingerprints := []string{}
	for _, entry := range d.Get("fingerprints").(*schema.Set).List() {
		fingerprints = append(fingerprints, entry.(string))
	}

	payload := SiteTlsSettings{
		Mandatory:                  d.Get("mandatory").(bool),
		Ports:                      ports,
		IsPortsException:           d.Get("is_ports_exception").(bool),
		Hosts:                      hosts,
		IsHostsException:           d.Get("is_hosts_exception").(bool),
		Fingerprints:               fingerprints,
		ForwardToOrigin:            d.Get("forward_to_origin").(bool),
		IsDisableSessionResumption: d.Get("is_disable_session_resumption").(bool),
	}

	if d.Get("forward_to_origin").(bool) == true {
		payload.HeaderName = d.Get("header_name").(string)
		payload.HeaderValue = d.Get("header_value").(string)
	}

	err = client.UpdateSiteTlsSetings(siteID, payload)

	if err != nil {
		return err
	}

	d.SetId(siteIDStr)
	return resourceSiteTlsSetingsRead(d, m)
}

func resourceSiteTlsSetingsRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)

	siteIDStr := d.Get("site_id").(string)
	siteID, err := strconv.Atoi(siteIDStr)
	if err != nil {
		return fmt.Errorf("failed to convert Site Id for Incapsula Site TLS Settings resource, actual value: %s, expected numeric id", siteIDStr)
	}

	siteTlsSettings, objectExists, err := client.GetSiteTlsSettings(
		siteID,
	)
	if !objectExists && !d.IsNewResource() {
		d.SetId("")
		return nil
	}

	if err != nil {
		return err
	}

	ports := &schema.Set{F: schema.HashInt}
	for i := range siteTlsSettings.Ports {
		ports.Add(siteTlsSettings.Ports[i])
	}

	fingerprints := &schema.Set{F: schema.HashString}
	for i := range siteTlsSettings.Fingerprints {
		fingerprints.Add(siteTlsSettings.Fingerprints[i])
	}

	if len(siteTlsSettings.Hosts) == 0 {
		log.Print("setting hosts to nil value")
		d.Set("hosts", nil)
	} else {
		hosts := &schema.Set{F: schema.HashString}
		for i := range siteTlsSettings.Hosts {
			hosts.Add(siteTlsSettings.Hosts[i])
		}
		log.Printf("hostsRes: %v", hosts)
		d.Set("hosts", hosts)

	}

	if err != nil {
		return err
	}

	d.Set("mandatory", siteTlsSettings.Mandatory)
	d.Set("ports", ports)
	d.Set("is_ports_exception", siteTlsSettings.IsPortsException)
	d.Set("is_hosts_exception", siteTlsSettings.IsHostsException)
	d.Set("fingerprints", fingerprints)
	d.Set("forward_to_origin", siteTlsSettings.ForwardToOrigin)
	d.Set("header_name", siteTlsSettings.HeaderName)
	d.Set("header_value", siteTlsSettings.HeaderValue)
	d.Set("is_disable_session_resumption", siteTlsSettings.IsDisableSessionResumption)

	d.SetId(siteIDStr)
	return nil
}

func resourceSiteTlsSetingsDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)

	siteIDStr := d.Get("site_id").(string)
	siteID, err := strconv.Atoi(siteIDStr)
	//todo katrin  edit error
	if err != nil {
		return fmt.Errorf("failed to convert Site Id for Incapsula Site TLS Settings resource, actual value: %s, expected numeric id", siteIDStr)
	}

	payload := SiteTlsSettings{
		Mandatory:                  mandatoryDefault,
		Ports:                      []int{},
		IsPortsException:           isPortsExceptionDefault,
		Hosts:                      []string{},
		IsHostsException:           isHostsExceptionDefault,
		Fingerprints:               []string{},
		ForwardToOrigin:            forwardToOriginDefault,
		IsDisableSessionResumption: isDisableSessionResumptionDefault,
		HeaderName:                 headerNameDefault,
		HeaderValue:                headerValueDefault,
	}

	err = client.UpdateSiteTlsSetings(siteID, payload)
	if err != nil {
		return fmt.Errorf("Failed to destroy Incapsula Site TLS Settings resource for Site ID %s, error:\n%s", siteIDStr, err)
	}

	d.SetId("")
	return nil
}
