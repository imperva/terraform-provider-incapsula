package incapsula

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
)

func updateHttp2Properties(client *Client, d *schema.ResourceData) error {
	enableHttp2ForNewSiteChanged := d.HasChange("enable_http2_for_new_sites") && d.Get("enable_http2_for_new_sites") != ""
	enableHttp2ToOriginForNewSitesChanged := d.HasChange("enable_http2_to_origin_for_new_sites") && d.Get("enable_http2_to_origin_for_new_sites") != ""

	log.Printf("[INFO] adi enable_http2_for_new_sites %v %v %v %v\n", d.HasChange("enable_http2_for_new_sites"), d.Get("enable_http2_for_new_sites"),
		d.HasChange("enable_http2_to_origin_for_new_sites"), d.Get("enable_http2_to_origin_for_new_sites"))

	if !enableHttp2ForNewSiteChanged && !enableHttp2ToOriginForNewSitesChanged {
		return nil
	}

	if d.Get("enable_http2_for_new_sites").(string) == "false" && d.Get("enable_http2_to_origin_for_new_sites").(string) == "true" {
		log.Printf("[ERROR] Could not update Incapsula account param enable_http2_for_new_sites with value (%s) and  enable_http2_to_origin_for_new_sites with value (%s) for account_id: %s",
			d.Get("enable_http2_for_new_sites"), d.Get("enable_http2_to_origin_for_new_sites"), d.Id())
		return fmt.Errorf("[ERROR] invalid values for enable_http2_for_new_sites and enable_http2_to_origin_for_new_sites")
	}

	updateParamsList := getParamsToUpdateInOrder(enableHttp2ForNewSiteChanged, enableHttp2ToOriginForNewSitesChanged, d)

	return updateParams(client, d, updateParamsList)
}

func getParamsToUpdateInOrder(enableHttp2ForNewSiteChanged bool, enableHttp2ToOriginForNewSitesChanged bool, d *schema.ResourceData) []string {

	updateParamsList := make([]string, 0)
	if enableHttp2ForNewSiteChanged && !enableHttp2ToOriginForNewSitesChanged {
		updateParamsList = append(updateParamsList, "enable_http2_for_new_sites")
	} else if !enableHttp2ForNewSiteChanged && enableHttp2ToOriginForNewSitesChanged {
		updateParamsList = append(updateParamsList, "enable_http2_to_origin_for_new_sites")
	} else if d.Get("enable_http2_to_origin_for_new_sites").(string) == "true" { // if the origin is true, then the client must be set first
		updateParamsList = append(updateParamsList, "enable_http2_for_new_sites", "enable_http2_to_origin_for_new_sites")
	} else {
		updateParamsList = append(updateParamsList, "enable_http2_to_origin_for_new_sites", "enable_http2_for_new_sites")
	}
	return updateParamsList
}

func updateParams(client *Client, d *schema.ResourceData, updateParams []string) error {
	for i := 0; i < len(updateParams); i++ {
		param := updateParams[i]
		if d.HasChange(param) && d.Get(param) != "" {
			log.Printf("[INFO] Updating Incapsula account param (%s) with value (%s) for account_id: %s\n", param, d.Get(param), d.Id())
			_, err := client.UpdateAccount(d.Id(), param, d.Get(param).(string))
			if err != nil {
				log.Printf("[ERROR] Could not update Incapsula account param (%s) with value (%t) for account_id: %s %s\n", param, d.Get(param).(bool), d.Id(), err)
				return err
			}
		}
	}
	return nil
}
