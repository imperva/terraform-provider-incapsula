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

func resourceMtlsClientToImpervaCertificateSetings() *schema.Resource {
	return &schema.Resource{
		Create: resourceeMtlsClientToImpervaCertificateSetingsUpdate,
		Read:   resourceeMtlsClientToImpervaCertificateSetingsRead,
		Update: resourceeMtlsClientToImpervaCertificateSetingsUpdate,
		Delete: resourceeMtlsClientToImpervaCertificateSetingsDelete,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				d.Set("site_id", d.Id())
				log.Printf("[DEBUG] Importing Incapsula MTLS Client to Imperva Certificate Site Settings for Site ID %s", d.Id())
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
			"require_client_certificate": {
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
				Description: "Disables SSL session resumption for site. Needed when Incapsula Client Certificate is needed only for specific hosts/ports and site have clients that reuse TLS session across different hosts/ports. Default: false.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     isDisableSessionResumptionDefault,
			},
		},
	}
}

func resourceeMtlsClientToImpervaCertificateSetingsUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)

	siteIDStr := d.Get("site_id").(string)
	siteID, err := strconv.Atoi(siteIDStr)
	if err != nil {
		return fmt.Errorf("failed to convert Site Id for Incapsula MTLS Client to Imperva Certificate Site Settings resource, actual value: %s, expected numeric id", siteIDStr)
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
		Mandatory:                  d.Get("require_client_certificate").(bool),
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
	return resourceeMtlsClientToImpervaCertificateSetingsRead(d, m)
}

func resourceeMtlsClientToImpervaCertificateSetingsRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)

	siteIDStr := d.Get("site_id").(string)
	siteID, err := strconv.Atoi(siteIDStr)
	if err != nil {
		return fmt.Errorf("failed to convert Site Id for Incapsula MTLS Client to Imperva Certificate Site Settings resource, actual value: %s, expected numeric id", siteIDStr)
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

	if len(siteTlsSettings.Ports) == 0 {
		d.Set("ports", nil)
	} else {
		ports := &schema.Set{F: schema.HashInt}
		for i := range siteTlsSettings.Ports {
			ports.Add(siteTlsSettings.Ports[i])
		}
		d.Set("ports", ports)
	}

	if len(siteTlsSettings.Fingerprints) == 0 {
		d.Set("fingerprints", nil)
	} else {
		fingerprints := &schema.Set{F: schema.HashString}
		for i := range siteTlsSettings.Fingerprints {
			fingerprints.Add(siteTlsSettings.Fingerprints[i])
		}
		d.Set("fingerprints", fingerprints)
	}

	if len(siteTlsSettings.Hosts) == 0 {
		d.Set("hosts", nil)
	} else {
		hosts := &schema.Set{F: schema.HashString}
		for i := range siteTlsSettings.Hosts {
			hosts.Add(siteTlsSettings.Hosts[i])
		}
		d.Set("hosts", hosts)
	}

	if err != nil {
		return err
	}

	d.Set("require_client_certificate", siteTlsSettings.Mandatory)
	d.Set("is_ports_exception", siteTlsSettings.IsPortsException)
	d.Set("is_hosts_exception", siteTlsSettings.IsHostsException)
	d.Set("forward_to_origin", siteTlsSettings.ForwardToOrigin)
	d.Set("header_name", siteTlsSettings.HeaderName)
	d.Set("header_value", siteTlsSettings.HeaderValue)
	d.Set("is_disable_session_resumption", siteTlsSettings.IsDisableSessionResumption)

	d.SetId(siteIDStr)
	return nil
}

func resourceeMtlsClientToImpervaCertificateSetingsDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)

	siteIDStr := d.Get("site_id").(string)
	siteID, err := strconv.Atoi(siteIDStr)
	if err != nil {
		return fmt.Errorf("failed to convert Site Id for Incapsula MTLS Client to Imperva Certificate Site Settings resource, actual value: %s, expected numeric id", siteIDStr)
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
		return fmt.Errorf("Failed to destroy Incapsula MTLS Client to Imperva Certificate Site Settings resource for Site ID %s, error:\n%s", siteIDStr, err)
	}

	d.SetId("")
	return nil
}
