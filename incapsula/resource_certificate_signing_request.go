package incapsula

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceCertificateSigningRequest() *schema.Resource {
	return &schema.Resource{
		Create: resourceCertificateSigningRequestCreate,
		Read:   resourceCertificateSigningRequestRead,
		Update: resourceCertificateSigningRequestUpdate,
		Delete: resourceCertificateSigningRequestDelete,

		Schema: map[string]*schema.Schema{
			// Required Arguments
			"site_id": {
				Description: "Numeric identifier of the site to operate on.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			// Optional Arguments
			"domain": {
				Description: "common name. For example: example.com.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"email": {
				Description: "Email address. For example: joe@example.com.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"country": {
				Description: "The two-letter ISO code for the country where your organization is located.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"state": {
				Description: "The state/region where your organization is located. This should not be abbreviated.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"city": {
				Description: "The city where your organization is located.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"organization": {
				Description: "The legal name of your organization. This should not be abbreviated or include suffixes such as Inc., Corp., or LLC.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"organization_unit": {
				Description: "The division of your organization handling the certificate. For example, IT Department.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			// Computed Arguments
			"csr_content": {
				Description: "The certificate request data.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func resourceCertificateSigningRequestCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	certificateSigningRequestResponse, err := client.CreateCertificateSigningRequest(
		d.Get("site_id").(string),
		d.Get("domain").(string),
		d.Get("email").(string),
		d.Get("country").(string),
		d.Get("state").(string),
		d.Get("city").(string),
		d.Get("organization").(string),
		d.Get("organization_unit").(string),
	)

	if err != nil {
		return err
	}

	d.Set("csr_content", certificateSigningRequestResponse.CsrContent)

	d.SetId(d.Get("site_id").(string))

	return nil
}

func resourceCertificateSigningRequestRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceCertificateSigningRequestUpdate(d *schema.ResourceData, m interface{}) error {
	if d.HasChange("domain") ||
		d.HasChange("email") ||
		d.HasChange("country") ||
		d.HasChange("state") ||
		d.HasChange("city") ||
		d.HasChange("organization") ||
		d.HasChange("organization_unit") {
		return resourceCertificateSigningRequestCreate(d, m)
	}
	return nil
}

func resourceCertificateSigningRequestDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}
