package incapsula

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"strconv"
)

const ignoreSensitiveVariableString = "Exported Certificate - data placeholder"

func validateInput(d *schema.ResourceData) (int, int, error) {
	siteIDStr := d.Get("site_id").(string)
	certificateIDStr := d.Get("certificate_id").(string)

	siteID, err := strconv.Atoi(siteIDStr)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to convert Site Id, actual value: %s, expected numeric id", siteIDStr)
	}

	certificateID, err := strconv.Atoi(certificateIDStr)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to convert certificate Id, actual value: %s, expected numeric id", certificateIDStr)
	}

	return siteID, certificateID, nil
}
