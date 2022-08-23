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
				Description: "The Site ID of the the site the API security is configured on.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"mandatory": {
				Description: "The certificate file in base64 format.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     mandatoryDefault,
			},
			"ports": {
				Description: "The certificate file in base64 format.", //todo KATRIN change
				Type:        schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Optional: true,
			},
			"is_ports_exception": {
				Description: "The certificate file in base64 format.", //todo KATRIN change
				Type:        schema.TypeBool,
				Default:     isPortsExceptionDefault,
				Optional:    true,
			},
			"hosts": {
				Description: "The certificate file in base64 format.", //todo KATRIN change
				Optional:    true,
				Type:        schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"is_hosts_exception": {
				Description:  "The certificate file in base64 format.", //todo KATRIN change
				Type:         schema.TypeBool,
				Default:      isHostsExceptionDefault,
				RequiredWith: []string{"hosts"},
				Optional:     true,
			},
			"fingerprints": {
				Description: "The certificate file in base64 format.", //todo KATRIN change
				Type:        schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
			},
			"forward_to_origin": {
				Description: "The certificate file in base64 format.", //todo KATRIN change
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     forwardToOriginDefault,
			},
			"header_name": {
				Description: "The certificate file in base64 format.", //todo KATRIN change
				Type:        schema.TypeString,
				Default:     headerNameDefault,
				Optional:    true,
			},
			"header_value": {
				Description:  "The certificate file in base64 format.", //todo KATRIN change
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"FULL_CERT", "COMMON_NAME", "FINGERPRINT", "SERIAL_NUMBER"}, false),
				Optional:     true,
				Default:      headerValueDefault,
			},
			"is_disable_session_resumption": {
				Description: "The certificate file in base64 format.", //todo KATRIN change
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
		return fmt.Errorf("failed to convert Site Id for Incapsula Site to Imperva to Origin mutual TLS Certificate Association resource, actual value: %s, expected numeric id", siteIDStr)
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

	siteTlsSettings, err := client.UpdateSiteTlsSetings(siteID, payload)

	if err != nil {
		return err
	}

	portsRes := &schema.Set{F: schema.HashInt}
	for i := range siteTlsSettings.Ports {
		portsRes.Add(siteTlsSettings.Ports[i])
	}

	hostsRes := &schema.Set{F: schema.HashString}
	for i := range siteTlsSettings.Hosts {
		hostsRes.Add(siteTlsSettings.Hosts[i])
	}

	fingerprintsRes := &schema.Set{F: schema.HashString}
	for i := range siteTlsSettings.Fingerprints {
		fingerprintsRes.Add(siteTlsSettings.Fingerprints[i])
	}

	// TODO: Setting this to arbitrary value as there is only one cert for each site.
	d.SetId(siteIDStr)
	d.Set("mandatory", siteTlsSettings.Mandatory)
	d.Set("ports", portsRes)
	d.Set("is_ports_exception", siteTlsSettings.IsPortsException)
	d.Set("hosts", hostsRes)
	d.Set("is_hosts_exception", siteTlsSettings.IsHostsException)
	d.Set("fingerprints", fingerprintsRes)
	d.Set("forward_to_origin", siteTlsSettings.ForwardToOrigin)
	d.Set("header_name", siteTlsSettings.HeaderName)
	d.Set("header_value", siteTlsSettings.HeaderValue)
	d.Set("is_disable_session_resumption", siteTlsSettings.IsDisableSessionResumption)
	return nil
}

func resourceSiteTlsSetingsRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)

	siteIDStr := d.Get("site_id").(string)
	siteID, err := strconv.Atoi(siteIDStr)
	if err != nil {
		return fmt.Errorf("failed to convert Site Id for Incapsula Site to Imperva to Origin mutual TLS Certificate Association resource, actual value: %s, expected numeric id", siteIDStr)
	}

	siteTlsSettings, err := client.GetSiteTlsSettings(
		siteID,
	)

	//katrin todo change error
	if err != nil {
		log.Printf("[ERROR] Could not update Incapsula API-security Site Configuration on site id: %d - %s\n", d.Get("site_id"), err)
		return err
	}

	log.Printf("%v", siteTlsSettings)

	ports := &schema.Set{F: schema.HashInt}
	for i := range siteTlsSettings.Ports {
		ports.Add(siteTlsSettings.Ports[i])
	}

	hosts := &schema.Set{F: schema.HashString}
	for i := range siteTlsSettings.Hosts {
		hosts.Add(siteTlsSettings.Hosts[i])
	}

	fingerprints := &schema.Set{F: schema.HashString}
	for i := range siteTlsSettings.Fingerprints {
		fingerprints.Add(siteTlsSettings.Fingerprints[i])
	}

	//mTLSCertificateData, err := client.GetClientCaCertificate(d.Id())
	if err != nil {
		return err
	}

	d.Set("mandatory", siteTlsSettings.Mandatory)
	d.Set("ports", ports)
	d.Set("is_ports_exception", siteTlsSettings.IsPortsException)
	d.Set("hosts", hosts)
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
		return fmt.Errorf("failed to convert Site Id for Incapsula Site to Imperva to Origin mutual TLS Certificate Association resource, actual value: %s, expected numeric id", siteIDStr)
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

	_, err = client.UpdateSiteTlsSetings(siteID, payload)
	if err != nil {
		return fmt.Errorf("Failed to destroy Incapsula Site TLS Settings resource for Site ID %s, error:\n%s", siteIDStr, err)
	}

	d.SetId("")
	return nil
}
