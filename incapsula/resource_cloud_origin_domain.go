package incapsula

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
	"net"
	"strconv"
	"strings"
)

var supportedCloudOriginRegions = []string{
	"us-east-1", "us-east-2", "us-west-1", "us-west-2",
	"eu-west-1", "eu-west-2", "eu-west-3", "eu-central-1", "eu-north-1",
	"ap-northeast-1", "ap-northeast-2", "ap-southeast-1", "ap-southeast-2",
	"ap-south-1", "sa-east-1",
}

func resourceCloudOriginDomain() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCloudOriginDomainCreate,
		ReadContext:   resourceCloudOriginDomainRead,
		UpdateContext: resourceCloudOriginDomainUpdate,
		DeleteContext: resourceCloudOriginDomainDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceCloudOriginDomainImport,
		},
		Schema: map[string]*schema.Schema{
			"site_id": {
				Description: "Numeric identifier of the site.",
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
			},
			"domain": {
				Description: "The origin domain (FQDN). Must be unique per site. Maximum 253 characters.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					domain := val.(string)

					// Check if domain is an IP address (not allowed)
					if net.ParseIP(domain) != nil {
						errs = append(errs, fmt.Errorf("%q must be a fully qualified domain name (FQDN), not an IP address, got: %s", key, domain))
						return
					}

					// Check domain length
					if len(domain) > 253 {
						errs = append(errs, fmt.Errorf("%q must be maximum 253 characters, got: %d characters", key, len(domain)))
						return
					}

					// Basic FQDN validation: must contain at least one dot
					if !strings.Contains(domain, ".") {
						errs = append(errs, fmt.Errorf("%q must be a fully qualified domain name (FQDN) with at least one dot, got: %s", key, domain))
						return
					}

					return
				},
			},
			"region": {
				Description: "AWS region where the cloud origin is located. Supported regions: us-east-1, us-east-2, us-west-1, us-west-2, eu-west-1, eu-west-2, eu-west-3, eu-central-1, eu-north-1, ap-northeast-1, ap-northeast-2, ap-southeast-1, ap-southeast-2, ap-south-1, sa-east-1",
				Type:        schema.TypeString,
				Required:    true,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					region := val.(string)
					found := false
					for _, supported := range supportedCloudOriginRegions {
						if region == supported {
							found = true
							break
						}
					}
					if !found {
						errs = append(errs, fmt.Errorf("%q must be one of %v, got: %s", key, supportedCloudOriginRegions, region))
					}
					return
				},
			},
			"port": {
				Description: "Port number for the origin. Valid range: 1-65535. Default: 443",
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     443,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					port := val.(int)
					if port < 1 || port > 65535 {
						errs = append(errs, fmt.Errorf("%q must be between 1 and 65535, got: %d", key, port))
					}
					return
				},
			},
			"imperva_origin_domain": {
				Description: "The Imperva-managed origin domain that is used to route traffic to the cloud origin.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"status": {
				Description: "The status of the cloud origin domain (e.g., PENDING, ACTIVE).",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"created_at": {
				Description: "Timestamp when the cloud origin domain was created.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"updated_at": {
				Description: "Timestamp when the cloud origin domain was last updated.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func resourceCloudOriginDomainCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)

	siteID := d.Get("site_id").(int)
	domain := d.Get("domain").(string)
	region := d.Get("region").(string)
	port := d.Get("port").(int)

	log.Printf("[INFO] Creating Incapsula cloud origin domain: %s for site: %d\n", domain, siteID)

	response, err := client.CreateCloudOriginDomain(siteID, 0, domain, region, port)
	if err != nil {
		return diag.Errorf("[ERROR] Could not create Incapsula cloud origin domain: %s for site: %d: %s\n", domain, siteID, err)
	}

	// Set the resource ID as site_id/origin_id
	originID := response.Value.OriginID
	syntheticID := fmt.Sprintf("%d/%d", siteID, originID)
	d.SetId(syntheticID)

	log.Printf("[INFO] Created Incapsula cloud origin domain: %s with ID: %d for site: %d\n", domain, originID, siteID)

	// Refresh state with full data from server
	return resourceCloudOriginDomainRead(ctx, d, m)
}

func resourceCloudOriginDomainRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)

	// Parse the composite ID
	if !strings.Contains(d.Id(), "/") {
		return diag.Errorf("[ERROR] Invalid cloud origin domain resource ID format: %s", d.Id())
	}

	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return diag.Errorf("[ERROR] Invalid cloud origin domain resource ID format: %s", d.Id())
	}

	siteID, err := strconv.Atoi(parts[0])
	if err != nil {
		return diag.Errorf("[ERROR] Invalid site ID in resource ID: %s", parts[0])
	}

	originID, err := strconv.Atoi(parts[1])
	if err != nil {
		return diag.Errorf("[ERROR] Invalid origin ID in resource ID: %s", parts[1])
	}

	log.Printf("[INFO] Reading Incapsula cloud origin domain: %d for site: %d\n", originID, siteID)

	response, err := client.GetCloudOriginDomain(siteID, 0, originID)
	if err != nil {
		log.Printf("[ERROR] Could not read Incapsula cloud origin domain: %d for site: %d: %s\n", originID, siteID, err)
		d.SetId("")
		return nil
	}

	// Set the attributes
	d.Set("site_id", siteID)
	d.Set("domain", response.Value.Domain)
	d.Set("region", response.Value.Region)
	d.Set("port", response.Value.Port)
	d.Set("imperva_origin_domain", response.Value.ImpervaOriginDomain)
	d.Set("status", response.Value.Status)
	d.Set("created_at", response.Value.CreatedAt)
	d.Set("updated_at", response.Value.UpdatedAt)

	return nil
}

func resourceCloudOriginDomainUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)

	// Parse the composite ID
	if !strings.Contains(d.Id(), "/") {
		return diag.Errorf("[ERROR] Invalid cloud origin domain resource ID format: %s", d.Id())
	}

	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return diag.Errorf("[ERROR] Invalid cloud origin domain resource ID format: %s", d.Id())
	}

	siteID, err := strconv.Atoi(parts[0])
	if err != nil {
		return diag.Errorf("[ERROR] Invalid site ID in resource ID: %s", parts[0])
	}

	originID, err := strconv.Atoi(parts[1])
	if err != nil {
		return diag.Errorf("[ERROR] Invalid origin ID in resource ID: %s", parts[1])
	}

	// Only region and port can be updated
	if d.HasChange("region") || d.HasChange("port") {
		region := d.Get("region").(string)
		port := d.Get("port").(int)

		log.Printf("[INFO] Updating Incapsula cloud origin domain: %d for site: %d\n", originID, siteID)

		_, err := client.UpdateCloudOriginDomain(siteID, 0, originID, region, port)
		if err != nil {
			return diag.Errorf("[ERROR] Could not update Incapsula cloud origin domain: %d for site: %d: %s\n", originID, siteID, err)
		}

		log.Printf("[INFO] Updated Incapsula cloud origin domain: %d for site: %d\n", originID, siteID)
	}

	// Refresh state with full data from server
	return resourceCloudOriginDomainRead(ctx, d, m)
}

func resourceCloudOriginDomainDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)

	// Parse the composite ID
	if !strings.Contains(d.Id(), "/") {
		return diag.Errorf("[ERROR] Invalid cloud origin domain resource ID format: %s", d.Id())
	}

	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return diag.Errorf("[ERROR] Invalid cloud origin domain resource ID format: %s", d.Id())
	}

	siteID, err := strconv.Atoi(parts[0])
	if err != nil {
		return diag.Errorf("[ERROR] Invalid site ID in resource ID: %s", parts[0])
	}

	originID, err := strconv.Atoi(parts[1])
	if err != nil {
		return diag.Errorf("[ERROR] Invalid origin ID in resource ID: %s", parts[1])
	}

	log.Printf("[INFO] Deleting Incapsula cloud origin domain: %d for site: %d\n", originID, siteID)

	err = client.DeleteCloudOriginDomain(siteID, 0, originID)
	if err != nil {
		return diag.Errorf("[ERROR] Could not delete Incapsula cloud origin domain: %d for site: %d: %s\n", originID, siteID, err)
	}

	log.Printf("[INFO] Deleted Incapsula cloud origin domain: %d for site: %d\n", originID, siteID)

	d.SetId("")
	return nil
}

func resourceCloudOriginDomainImport(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	// Expect import ID format: site_id/origin_id
	importID := d.Id()

	if importID == "" {
		return nil, fmt.Errorf("expected import ID in format 'site_id/origin_id'")
	}

	if !strings.Contains(importID, "/") {
		return nil, fmt.Errorf("invalid import ID format %q: expected 'site_id/origin_id'", importID)
	}

	parts := strings.SplitN(importID, "/", 2)
	if len(parts) != 2 || strings.TrimSpace(parts[0]) == "" || strings.TrimSpace(parts[1]) == "" {
		return nil, fmt.Errorf("invalid import ID format %q: expected 'site_id/origin_id'", importID)
	}

	siteID := strings.TrimSpace(parts[0])
	originID := strings.TrimSpace(parts[1])

	// Set the ID
	d.SetId(fmt.Sprintf("%s/%s", siteID, originID))

	// Trigger a read to populate all attributes
	return []*schema.ResourceData{d}, nil
}
