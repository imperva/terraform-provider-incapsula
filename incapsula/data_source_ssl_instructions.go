package incapsula

import (
	"bytes"
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"golang.org/x/exp/slices"
	"strconv"
	"time"
)

func dataSourceSSLInstructions() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSSLInstructionsRead,
		Description: "Provides data about SSL instructions",

		Schema: map[string]*schema.Schema{
			"site_id": {
				Description: "Numeric identifier of the site to operate on.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"domain_ids": {
				Description: "domain ids.",
				Type:        schema.TypeSet,
				Required:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"managed_certificate_settings_id": {
				Description: "Numeric identifier of the site certificate id.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"validation_timeout": {
				Description: "Length of time in seconds to wait for SSL validation.",
				Type:        schema.TypeInt,
				Required:    false,
				Default:     200,
			},
			// Computed Attributes
			"instructions": {
				Description: "A set of SSL instructions.",
				Computed:    true,
				Type:        schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"domain_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"san_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:        schema.TypeString,
							Description: "Numeric identifier of the certificate on Imperva service",
							Computed:    true,
						},
						"value": {
							Type:        schema.TypeString,
							Description: "Status of the certificate",
							Computed:    true,
						},
					}},
			},
		},
	}
}

func dataSourceSSLInstructionsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)
	siteIdStr := d.Get("site_id").(string)
	siteId, _ := strconv.Atoi(siteIdStr)
	timeoutStr := d.Get("validation_timeout").(string)
	timeout, _ := strconv.Atoi(timeoutStr)
	waitForInstructions(client, siteId, timeout)
	sSLInstructionsResponse, err := client.GetSiteSSLInstructions(siteId)
	if err != nil {
		return diag.Errorf("Error request site SSL instructions: %v", err)
	}
	if sSLInstructionsResponse.Data != nil && len(sSLInstructionsResponse.Data) > 0 {
		populateInstructionsFromDTO(d, sSLInstructionsResponse.Data)
	}
	d.SetId(siteIdStr)
	return nil
}

func getSslValidationStatus(client *Client, siteId int, dtoData []SiteDomainDetails) bool {
	siteCertificateV3Response, _ := client.GetSiteCertificateRequestStatus(siteId, nil)
	b := siteCertificateV3Response != nil && siteCertificateV3Response.Data != nil && len(siteCertificateV3Response.Data) > 0
	b = b && siteCertificateV3Response.Data[0].CertificatesDetails != nil && validateSans(siteCertificateV3Response.Data[0], dtoData)
	return b
}

func waitForInstructions(client *Client, siteId int, timeout int) diag.Diagnostics {
	siteDomainDetailsDto, err := client.GetWebsiteDomains(strconv.Itoa(siteId))
	if err != nil {
		return diag.Errorf("Error getting site %d domains: %v", siteId, err)
	}
	if siteDomainDetailsDto.Data != nil && len(siteDomainDetailsDto.Data) > 0 {
		// Base case
		if getSslValidationStatus(client, siteId, siteDomainDetailsDto.Data) {
			return nil
		}
		// Try with delay until timeout
		limit := time.Duration(timeout) * time.Second
		start := time.Now()
		current := time.Now()
		for current.Sub(start) <= limit {
			time.Sleep(10 * time.Second)
			current = time.Now()

			if getSslValidationStatus(client, siteId, siteDomainDetailsDto.Data) {
				return nil
			}
		}
		panic("managed certificate was not created as expected")
	}
	return nil
}

func validateSans(siteCertificateDTO SiteCertificateDTO, siteDomainDetails []SiteDomainDetails) bool {
	if len(siteDomainDetails) == 0 { //todo handle multiple certificates
		return len(siteCertificateDTO.CertificatesDetails) == 0 || validateEmptyCertificates(siteCertificateDTO.CertificatesDetails)
	} else {
		return len(siteCertificateDTO.CertificatesDetails) > 0 && validateDomainsSans(siteDomainDetails, siteCertificateDTO.CertificatesDetails)
	}
}

func validateEmptyCertificates(certificatesDetails []CertificateDTO) bool {
	for _, t := range certificatesDetails {
		if t.Sans != nil && len(t.Sans) > 0 && !allSansAreDeleted(t.Sans) {
			return false
		}
	}
	return true
}

func allSansAreDeleted(sans []SanDTO) bool {
	deleteStatuses := []string{"REMOVED", "DELETED_PENDING_PUBLICATION", "DELETED_LOCALLY", "CANCELED_LOCALLY", "PENDING_CANCELATION", "PENDING_DELETION"}
	for _, t := range sans {
		if !slices.Contains(deleteStatuses, t.Status) {
			return false
		}
	}
	return true
}

func validateDomainsSans(siteDomainDetails []SiteDomainDetails, certificateDTO []CertificateDTO) bool {
	return validateSanForEachDomain(siteDomainDetails, certificateDTO) && validateRedundantSansAreDeleted(siteDomainDetails, certificateDTO)
}

func validateSanForEachDomain(siteDomainDetails []SiteDomainDetails, certificateDTO []CertificateDTO) bool {
	for _, t := range siteDomainDetails {
		if !validateSanForDomain(t, certificateDTO) {
			return false
		}
	}
	return true
}

func validateSanForDomain(domain SiteDomainDetails, certificates []CertificateDTO) bool {
	for _, t := range certificates {
		if t.Sans != nil && len(t.Sans) > 0 {
			for _, s := range t.Sans {
				if slices.Contains(s.DomainIds, domain.Id) && s.SanValue == domain.Domain {
					return true
				}
			}
		}
	}
	return false
}

func validateRedundantSansAreDeleted(domains []SiteDomainDetails, certificates []CertificateDTO) bool {
	res := true
	for _, t := range certificates {
		if t.Sans != nil && len(t.Sans) > 0 {
			for _, s := range t.Sans {
				res = res && validateExistDomainOrDeletedSan(s, domains)
			}
		}
	}
	return res
}

func validateExistDomainOrDeletedSan(s SanDTO, domains []SiteDomainDetails) bool {
	deleteStatuses := []string{"REMOVED", "DELETED_PENDING_PUBLICATION", "DELETED_LOCALLY", "CANCELED_LOCALLY", "PENDING_CANCELATION", "PENDING_DELETION"}
	for _, t := range domains {
		if t.Domain == s.SanValue && slices.Contains(s.DomainIds, t.Id) {
			return true
		}
	}
	return !slices.Contains(deleteStatuses, s.Status)
}

func populateInstructionsFromDTO(d *schema.ResourceData, instructionsDTO []InstructionsDTO) {
	instructions := &schema.Set{F: getHashSanId}
	for _, v := range instructionsDTO {
		ins := map[string]interface{}{}
		ins["name"] = v.Domain
		ins["type"] = v.RecordType
		ins["value"] = v.VerificationCode
		ins["domain_id"] = v.RelatedSansDetails[0].DomainIds[0]
		ins["san_id"] = v.RelatedSansDetails[0].SanId
		instructions.Add(ins)
	}
	d.Set("instructions", instructions)
}

func getHashSanId(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})

	if v, ok := m["san_id"]; ok {
		buf.WriteString(fmt.Sprintf("%s-", strconv.Itoa(v.(int))))
	}
	return PositiveHash(buf.String())
}
