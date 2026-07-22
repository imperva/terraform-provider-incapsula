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

func resourceCloudOriginDomain() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCloudOriginDomainCreate,
		ReadContext:   resourceCloudOriginDomainRead,
		DeleteContext: resourceCloudOriginDomainDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceCloudOriginDomainImport,
		},
		CustomizeDiff: func(ctx context.Context, d *schema.ResourceDiff, m interface{}) error {
			if d.Id() == "" {
				return nil
			}
			for _, f := range []string{"region", "port", "origin_tls_policy"} {
				if d.HasChange(f) {
					return fmt.Errorf(
						"%q cannot be modified after creation. To change it, run "+
							"`terraform destroy` against this resource and re-create it "+
							"(traffic routed through the imperva_origin_domain will be interrupted).",
						f)
				}
			}
			return nil
		},
		Schema: map[string]*schema.Schema{
			"account_id": {
				Description: "Numeric identifier of the account to operate on.",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
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

					if net.ParseIP(domain) != nil {
						errs = append(errs, fmt.Errorf("%q must be a fully qualified domain name (FQDN), not an IP address, got: %s", key, domain))
						return
					}

					if len(domain) > 253 {
						errs = append(errs, fmt.Errorf("%q must be maximum 253 characters, got: %d characters", key, len(domain)))
						return
					}

					if !strings.Contains(domain, ".") {
						errs = append(errs, fmt.Errorf("%q must be a fully qualified domain name (FQDN) with at least one dot, got: %s", key, domain))
						return
					}

					return
				},
			},
			"region": {
				Description: "The cloud region where the origin is located (e.g., us-east-1 for AWS). Immutable after creation.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"port": {
				Description: "Port number for the origin. Must be 443 or in the range 1024-65535. Default: 443. Immutable after creation.",
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     443,
				ForceNew:    true,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					port := val.(int)
					if port != 443 && (port < 1024 || port > 65535) {
						errs = append(errs, fmt.Errorf("%q must be 443 or between 1024 and 65535, got: %d", key, port))
					}
					return
				},
			},
			"origin_tls_policy": {
				Description: "Minimum TLS version for the origin connection. Supported values: SSLv3, TLS_1_0, TLS_1_1, TLS_1_2. Immutable after creation.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"imperva_origin_domain": {
				Description: "The Imperva-managed origin domain that is used to route traffic to the cloud origin.",
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

func parseCloudOriginID(id string) (accountID string, siteID int, originID int, err error) {
	parts := strings.Split(id, "/")
	if len(parts) != 3 {
		err = fmt.Errorf("invalid cloud origin domain resource ID format %q: expected 'account_id/site_id/origin_id'", id)
		return
	}

	accountID = parts[0]

	siteID, err = strconv.Atoi(parts[1])
	if err != nil {
		err = fmt.Errorf("invalid site ID in resource ID: %s", parts[1])
		return
	}

	originID, err = strconv.Atoi(parts[2])
	if err != nil {
		err = fmt.Errorf("invalid origin ID in resource ID: %s", parts[2])
		return
	}

	return
}

func resourceCloudOriginDomainCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)

	siteID := d.Get("site_id").(int)
	accountID, _ := d.Get("account_id").(string)
	domain := d.Get("domain").(string)
	region := d.Get("region").(string)
	port := d.Get("port").(int)
	tlsPolicy := d.Get("origin_tls_policy").(string)

	log.Printf("[INFO] Creating Incapsula cloud origin domain: %s for site: %d\n", domain, siteID)

	response, err := client.CreateCloudOriginDomain(siteID, accountID, domain, region, port, tlsPolicy)
	if err != nil {
		return diag.Errorf("[ERROR] Could not create Incapsula cloud origin domain: %s for site: %d: %s\n", domain, siteID, err)
	}

	if len(response.Data) == 0 {
		return diag.Errorf("[ERROR] Empty response when creating cloud origin domain: %s for site: %d", domain, siteID)
	}

	originID := response.Data[0].ID
	d.SetId(fmt.Sprintf("%s/%d/%d", accountID, siteID, originID))

	log.Printf("[INFO] Created Incapsula cloud origin domain: %s with ID: %d for site: %d\n", domain, originID, siteID)

	return resourceCloudOriginDomainRead(ctx, d, m)
}

func resourceCloudOriginDomainRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)

	accountID, siteID, originID, err := parseCloudOriginID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Reading Incapsula cloud origin domain: %d for site: %d\n", originID, siteID)

	response, err := client.GetCloudOriginDomain(siteID, originID, accountID)
	if err != nil {
		log.Printf("[ERROR] Could not read Incapsula cloud origin domain: %d for site: %d: %s\n", originID, siteID, err)
		return diag.Errorf("[ERROR] Could not read Incapsula cloud origin domain: %d for site: %d: %s", originID, siteID, err)
	}

	if len(response.Data) == 0 {
		return diag.Errorf("[ERROR] Empty response when reading cloud origin domain: %d for site: %d", originID, siteID)
	}

	origin := response.Data[0]
	d.Set("account_id", accountID)
	d.Set("site_id", siteID)
	d.Set("domain", origin.OriginDomain)
	d.Set("region", origin.Region)
	d.Set("port", origin.OriginConfig.Port)
	d.Set("origin_tls_policy", origin.OriginConfig.OriginTlsPolicy)
	d.Set("imperva_origin_domain", origin.ImpervaOriginDomain)
	d.Set("created_at", origin.CreatedAt)
	d.Set("updated_at", origin.UpdatedAt)

	return nil
}

func resourceCloudOriginDomainDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)

	accountID, siteID, originID, err := parseCloudOriginID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Deleting Incapsula cloud origin domain: %d for site: %d\n", originID, siteID)

	err = client.DeleteCloudOriginDomain(siteID, originID, accountID)
	if err != nil {
		return diag.Errorf("[ERROR] Could not delete Incapsula cloud origin domain: %d for site: %d: %s\n", originID, siteID, err)
	}

	log.Printf("[INFO] Deleted Incapsula cloud origin domain: %d for site: %d\n", originID, siteID)

	d.SetId("")
	return nil
}

func resourceCloudOriginDomainImport(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	importID := d.Id()
	if importID == "" {
		return nil, fmt.Errorf("expected import ID in format 'account_id/site_id/origin_id'")
	}

	parts := strings.Split(importID, "/")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid import ID format %q: expected 'account_id/site_id/origin_id'", importID)
	}

	accountID := strings.TrimSpace(parts[0])
	siteID := strings.TrimSpace(parts[1])
	originID := strings.TrimSpace(parts[2])

	if accountID == "" || siteID == "" || originID == "" {
		return nil, fmt.Errorf("invalid import ID format %q: expected 'account_id/site_id/origin_id'", importID)
	}

	if _, err := strconv.Atoi(siteID); err != nil {
		return nil, fmt.Errorf("invalid site_id %q: must be numeric", siteID)
	}
	if _, err := strconv.Atoi(originID); err != nil {
		return nil, fmt.Errorf("invalid origin_id %q: must be numeric", originID)
	}

	d.SetId(fmt.Sprintf("%s/%s/%s", accountID, siteID, originID))
	d.Set("account_id", accountID)

	return []*schema.ResourceData{d}, nil
}
