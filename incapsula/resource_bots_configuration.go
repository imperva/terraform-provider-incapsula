package incapsula

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
	"sort"
)

var canceledGoodBots = "canceled_good_bots"
var badBots = "bad_bots"
var botId = "id"

func resourceBotsConfiguration() *schema.Resource {
	return &schema.Resource{
		Create: resourceBotsConfigurationCreate,
		Read:   resourceBotsConfigurationRead,
		Update: resourceBotsConfigurationCreate,
		Delete: resourceBotsConfigurationDelete,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				d.Set("site_id", d.Id())
				return []*schema.ResourceData{d}, nil
			},
		},

		Schema: map[string]*schema.Schema{
			// Required Arguments
			"site_id": {
				Description: "Numeric identifier of the site to operate on.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			// Optional Arguments
			canceledGoodBots: {
				Description: "List of bots for Canceled Good Bots configuration",
				Type:        schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Optional: true,
			},
			badBots: {
				Description: "List of bots for Bad Bots configuration",
				Type:        schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Optional: true,
			},
		},
	}
}

func populateFromConfBots(d *schema.ResourceData, resourceKey string) []BotStruct {
	botsConf := d.Get(resourceKey).(*schema.Set)

	var botsStructs = make([]BotStruct, len(botsConf.List()))

	var dcInd = 0
	for _, botConf := range botsConf.List() {
		botId := botConf.(int)

		botsStructs[dcInd] = BotStruct{}
		if botId != 0 {
			botsStructs[dcInd].ID = &botId
		}
		dcInd++
	}

	log.Printf("[DEBUG] populateFromConfBots - botsStructs: %+v\n", botsStructs)
	return botsStructs
}

func populateFromConfBotsConfigurationDTO(d *schema.ResourceData) BotsConfigurationDTO {
	requestDTO := BotsConfigurationDTO{}
	requestDTO.Data = make([]BotsStruct, 1)
	requestDTO.Data[0].CanceledGoodBots = populateFromConfBots(d, canceledGoodBots)
	requestDTO.Data[0].BadBots = populateFromConfBots(d, badBots)

	log.Printf("[DEBUG] populateFromConfBotsConfigurationDTO - requestDTO: %+v\n", requestDTO)
	return requestDTO
}

func resourceBotsConfigurationCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)

	requestDTO := populateFromConfBotsConfigurationDTO(d)

	responseDTO, err := client.UpdateBotAccessControlConfiguration(d.Get("site_id").(string), requestDTO)
	if err != nil {
		return fmt.Errorf("Error updating Bots configuration for site (%s): %s",
			d.Get("site_id"), err)
	}

	if responseDTO.Errors != nil && len(responseDTO.Errors) > 0 {
		return fmt.Errorf("Error updating Bots configuration for site (%s): %s",
			d.Get("site_id"), responseDTO.Errors)
	}

	// Set the dc ID
	d.SetId(d.Get("site_id").(string))

	return resourceBotsConfigurationRead(d, m)
}

func resourceBotsConfigurationRead(d *schema.ResourceData, m interface{}) error {
	// Implement by reading the ListBotsResponse for the bots
	client := m.(*Client)

	responseDTO, err := client.GetBotAccessControlConfiguration(d.Get("site_id").(string))
	if err != nil {
		return fmt.Errorf("Error getting Bots configuration for site (%s): %s", d.Get("site_id"), err)
	}

	if responseDTO.Errors != nil && len(responseDTO.Errors) > 0 {
		if responseDTO.Errors[0].Status == "404" {
			log.Printf("[INFO] Incapsula Site with ID %s has already been deleted: %s\n", d.Get("site_id"), responseDTO.Errors)
			d.SetId("")
			return nil
		}

		out, err := json.Marshal(responseDTO.Errors)
		if err != nil {
			panic(err)
		}
		return fmt.Errorf("Error getting Bots configuration for site (%s): %s", d.Get("site_id"), string(out))
	}

	canceledGoodBotsList := make([]int, 0)
	for _, bot := range responseDTO.Data[0].CanceledGoodBots {
		canceledGoodBotsList = append(canceledGoodBotsList, *bot.ID)
	}
	sort.Ints(canceledGoodBotsList)
	d.Set(canceledGoodBots, canceledGoodBotsList)

	badBotsList := make([]int, 0)
	for _, bot := range responseDTO.Data[0].BadBots {
		badBotsList = append(badBotsList, *bot.ID)
	}
	sort.Ints(badBotsList)
	d.Set(badBots, badBotsList)

	log.Printf("[DEBUG] resourceBotsConfigurationRead - canceledGoodBots: %+v\n", canceledGoodBotsList)
	log.Printf("[DEBUG] resourceBotsConfigurationRead - badBotsList: %+v\n", badBotsList)

	return nil
}

func resourceBotsConfigurationDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)

	responseDTO, err := client.GetBotAccessControlConfiguration(d.Get("site_id").(string))
	if err != nil {
		return fmt.Errorf("Error deleting Bots configuration for site (%s): %s", d.Get("site_id"), err)
	}

	if responseDTO.Errors != nil && len(responseDTO.Errors) > 0 && responseDTO.Errors[0].Status != "404" {
		out, err := json.Marshal(responseDTO.Errors)
		if err != nil {
			panic(err)
		}
		return fmt.Errorf("Error deleting Bots configuration for site (%s): %s", d.Get("site_id"), string(out))
	}

	d.SetId("")
	return nil
}
