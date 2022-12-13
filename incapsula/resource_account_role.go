package incapsula

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
	"strconv"
)

func resourceAccountRole() *schema.Resource {
	return &schema.Resource{
		Create: resourceAccountRoleCreate,
		Read:   resourceAccountRoleRead,
		Update: resourceAccountRoleUpdate,
		Delete: resourceAccountRoleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			// Required Arguments
			"account_id": {
				Description: "Numeric identifier of the account to operate on.",
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
			},
			"name": {
				Description: "The role name.",
				Type:        schema.TypeString,
				Required:    true,
			},

			// Optional Arguments
			"permissions": {
				Description: "List of account permission keys that the role contains.",
				Type:        schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
			},
			"description": {
				Description: "The role description.",
				Type:        schema.TypeString,
				Optional:    true,
			},
		},
	}
}

func populateRolePermissions(d *schema.ResourceData) []string {
	permissions := d.Get("permissions").(*schema.Set)

	var permissionsSlice = make([]string, len(permissions.List()))

	var dcInd = 0
	for _, permission := range permissions.List() {
		permissionKey := permission.(string)

		if permissionKey != "" {
			permissionsSlice[dcInd] = permissionKey
		}
		dcInd++
	}

	log.Printf("[DEBUG] populateRolePermissions - RolePermissions: %+v\n", permissionsSlice)
	return permissionsSlice
}

func populateRoleDetailsDTO(d *schema.ResourceData) RoleDetailsBasicDTO {
	requestDTO := RoleDetailsBasicDTO{}
	requestDTO.RoleName = d.Get("name").(string)
	requestDTO.RoleAbilities = populateRolePermissions(d)
	requestDTO.RoleDescription = d.Get("description").(string)

	log.Printf("[DEBUG] populateRoleDetailsDTO - RoleDetailsDTO: %+v\n", requestDTO)
	return requestDTO
}

func resourceAccountRoleCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	accountId := d.Get("account_id").(int)

	log.Printf("[INFO] Creating Incapsula account role for account: %d\n", accountId)

	roleDetailsBasicDTO := populateRoleDetailsDTO(d)

	requestDTO := RoleDetailsCreateDTO{}
	requestDTO.RoleDetailsBasicDTO = roleDetailsBasicDTO
	requestDTO.AccountId = d.Get("account_id").(int)

	responseDTO, err := client.AddAccountRole(requestDTO)

	if err != nil {
		log.Printf("[ERROR] Could not create Incapsula account role: %s\n", err)
		return err
	}

	// Set the Account Role ID
	d.SetId(strconv.Itoa(responseDTO.RoleId))
	log.Printf("[INFO] Created Incapsula account role: %s, Id: %d\n", responseDTO.RoleName, responseDTO.RoleId)

	// Set the rest of the state from the resource read
	return resourceAccountRoleRead(d, m)
}

func resourceAccountRoleRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)

	roleID, _ := strconv.Atoi(d.Id())

	log.Printf("[INFO] Reading Incapsula account role ID: %d\n", roleID)

	accountRoleResponse, err := client.GetAccountRole(roleID)

	// Account object may have been deleted
	if accountRoleResponse != nil && accountRoleResponse.ErrorCode == 1047 {
		log.Printf("[INFO] Incapsula Account Role with ID %d does not exist: %s\n", roleID, err)
		d.SetId("")
		return nil
	}

	if err != nil {
		log.Printf("[ERROR] Could not read Incapsula account role ID: %d, %s\n", roleID, err)
		return err
	}

	d.Set("account_id", accountRoleResponse.AccountId)
	d.Set("name", accountRoleResponse.RoleName)
	d.Set("description", accountRoleResponse.RoleDescription)
	abilitiesList := make([]string, 0)
	for _, roleAbility := range accountRoleResponse.RoleAbilities {
		abilitiesList = append(abilitiesList, roleAbility.AbilityKey)
	}
	d.Set("permissions", abilitiesList)

	log.Printf("[INFO] Finished reading Incapsula account role id: %d\n", roleID)

	return nil
}

func resourceAccountRoleUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	roleID, _ := strconv.Atoi(d.Id())

	log.Printf("[INFO] Updating Incapsula account role with id: %d\n", roleID)

	requestDTO := populateRoleDetailsDTO(d)
	responseDTO, err := client.UpdateAccountRole(roleID, requestDTO)

	if err != nil {
		log.Printf("[ERROR] Could not create Incapsula account role: %s\n", err)
		return err
	}

	log.Printf("[INFO] Updated Incapsula account role: %s, Id: %d\n", responseDTO.RoleName, responseDTO.RoleId)
	// Set the rest of the state from the resource read
	return resourceAccountRoleRead(d, m)
}

func resourceAccountRoleDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	roleID, _ := strconv.Atoi(d.Id())

	log.Printf("[INFO] Deleting Incapsula account role with id: %d\n", roleID)

	err := client.DeleteAccountRole(roleID)

	if err != nil {
		log.Printf("[ERROR] Could not delete Incapsula account role with id: %d, %s\n", roleID, err)
		return err
	}

	// Set the ID to empty
	// Implicitly clears the resource
	d.SetId("")

	log.Printf("[INFO] Deleted Incapsula account role with id: %d\n", roleID)

	return nil
}
