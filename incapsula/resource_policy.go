package incapsula

import (
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourcePolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourcePolicyCreate,
		Read:   resourcePolicyRead,
		Update: resourcePolicyUpdate,
		Delete: resourcePolicyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			// Required Arguments
			"name": {
				Description: "The policy name.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"enabled": {
				Description: "Enables the policy.",
				Type:        schema.TypeBool,
				Required:    true,
			},
			"policy_type": {
				Description: "The policy type. Possible values: ACL, WHITELIST.",
				Type:        schema.TypeString,
				Required:    true,
			},

			// Optional Arguments
			"account_id": {
				Description: "The Account ID of the policy.",
				Type:        schema.TypeInt,
				Computed:    true,
				Optional:    true,
			},
			"description": {
				Description: "The policy description.",
				Type:        schema.TypeString,
				Optional:    true,
			},
		},
	}
}

func resourcePolicyCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)

	policyLite := PolicyLite{
		Name:           d.Get("name").(string),
		Enabled:        d.Get("enabled").(bool),
		PolicyType:     d.Get("policy_type").(string),
		AccountID:      d.Get("account_id").(int),
		Description:    d.Get("description").(string),
		PolicySettings: make([]int, 0),
	}

	policyAddResponse, err := client.AddPolicy(&policyLite)

	if err != nil {
		log.Printf("[ERROR] Could not create Incapsula policy: %s - %s\n", policyLite.Name, err)
		return err
	}

	// Set the policyID
	policyID := strconv.Itoa(policyAddResponse.Value.ID)
	d.SetId(policyID)
	log.Printf("[INFO] Created Incapsula policy with ID: %s\n", policyID)

	return resourceDataCenterRead(d, m)
}

func resourcePolicyRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)

	policyID := d.Get("id").(string)
	policyGetResponse, err := client.GetPolicy(policyID)

	if err != nil {
		log.Printf("[ERROR] Could not get Incapsula policy: %s - %s\n", policyID, err)
		return err
	}

	// Set computed values
	d.Set("name", policyGetResponse.Value.Name)
	d.Set("enabled", policyGetResponse.Value.Enabled)
	d.Set("policy_type", policyGetResponse.Value.PolicyType)
	d.Set("account_id", policyGetResponse.Value.AccountID)
	d.Set("description", policyGetResponse.Value.Description)

	return nil
}

func resourcePolicyUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)

	id, err := strconv.Atoi(d.Get("id").(string))
	if err != nil {
		return err
	}

	policyLite := PolicyLite{
		ID:             id,
		Name:           d.Get("name").(string),
		Enabled:        d.Get("enabled").(bool),
		PolicyType:     d.Get("policy_type").(string),
		AccountID:      d.Get("account_id").(int),
		Description:    d.Get("description").(string),
		PolicySettings: make([]int, 0),
	}

	_, err = client.UpdatePolicy(&policyLite)

	if err != nil {
		return err
	}

	return nil
}

func resourcePolicyDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)

	err := client.DeletePolicy(d.Id())

	if err != nil {
		return err
	}

	// Set the ID to empty
	// Implicitly clears the resource
	d.SetId("")

	return nil
}
