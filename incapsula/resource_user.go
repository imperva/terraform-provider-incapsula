package incapsula

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
	"strconv"
	"strings"
	"time"
)

func resourceUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceUserCreate,
		Read:   resourceUserRead,
		Update: resourceUserUpdate,
		Delete: resourceUserDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			// Required Arguments
			"email": {
				Description: "Email address. For example: joe@example.com. example: userEmail@imperva.com",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"account_id": {
				Description: "Unique ID of the required account . example: 123456",
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
			},
			"first_name": {
				Description: "The first name of the user that was acted on. example: John",
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
			},
			"last_name": {
				Description: "The last name of the user that was acted on. example: Snow",
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
			},
			"role_ids": {
				Description: "List of role ids to add to the user.",
				Type:        schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Optional: true,
			},
			"role_names": {
				Description: "List of role names.",
				Type:        schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed: true,
			},
		},
	}
}

func resourceUserCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	email := d.Get("email").(string)
	accountId := d.Get("account_id").(int)

	log.Printf("[INFO] Creating Incapsula user for email: %s\n", email)

	UserAddResponse, err := client.AddUser(
		accountId,
		email,
		d.Get("role_ids").([]interface{}),
		d.Get("first_name").(string),
		d.Get("last_name").(string),
	)

	if err != nil {
		log.Printf("[ERROR] Could not create user for email: %s, %s\n", email, err)
		return err
	}

	// Set the User ID
	d.SetId(fmt.Sprintf("%s_%s", strconv.Itoa(accountId), email))
	log.Printf("[INFO] Created Incapsula user for email: %s userid: %d\n", email, UserAddResponse.UserID)

	// There may be a timing/race condition here
	// Set an arbitrary period to sleep
	time.Sleep(3 * time.Second)

	// Set the rest of the state from the resource read
	return resourceUserRead(d, m)
}

func resourceUserRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	userID := d.Id()
	stringSlice := strings.Split(userID, "_")
	accountID, _ := strconv.Atoi(stringSlice[0])
	email := stringSlice[1]
	log.Printf("[INFO] Reading Incapsula user : %s\n", userID)

	UserStatusResponse, err := client.UserStatus(accountID, email)

	if err != nil {
		log.Printf("[ERROR] Could not read Incapsula user: %s, %s\n", email, err)
		return err
	}

	log.Printf("[INFO]listRoles : %v\n", UserStatusResponse.Roles)

	listRolesids := make([]interface{}, len(UserStatusResponse.Roles))
	listRolesnames := make([]interface{}, len(UserStatusResponse.Roles))
	for i, v := range UserStatusResponse.Roles {
		log.Printf("[INFO]listRoles : %v\n", UserStatusResponse.Roles)
		listRolesids[i] = v.RoleID
		listRolesnames[i] = v.RoleName
	}

	d.Set("email", UserStatusResponse.Email)
	d.Set("account_id", UserStatusResponse.AccountID)
	d.Set("first_name", UserStatusResponse.FirstName)
	d.Set("last_name", UserStatusResponse.LastName)
	d.Set("role_ids", listRolesids)
	d.Set("role_names", listRolesnames)

	log.Printf("[INFO] Finished reading Incapsula user: %s\n", email)

	return nil
}

func resourceUserUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	email := d.Get("email").(string)
	accountId := d.Get("account_id").(int)

	log.Printf("[INFO] Creating Incapsula user for email: %s\n", email)

	UserUpdateResponse, err := client.UpdateUser(
		accountId,
		email,
		d.Get("role_ids").([]interface{}),
	)
	if err != nil {
		log.Printf("[ERROR] Could not update user for email: %s, %s\n", email, err)
		return err
	}

	log.Printf("[Info] New Roles for user %s : %+v\n", email, UserUpdateResponse.Roles)

	// There may be a timing/race condition here
	// Set an arbitrary period to sleep
	time.Sleep(3 * time.Second)

	// Set the rest of the state from the resource read
	return resourceUserRead(d, m)
}

func resourceUserDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	accountID := d.Get("account_id").(int)
	email := d.Get("email").(string)

	log.Printf("[INFO] Deleting Incapsula user: %s\n", email)

	err := client.DeleteUser(accountID, email)

	if err != nil {
		log.Printf("[ERROR] Could not delete Incapsula user: %s %s\n", email, err)
		return err
	}

	// Set the ID to empty
	// Implicitly clears the resource
	d.SetId("")

	log.Printf("[INFO] Deleted Incapsula user: %s\n", email)

	return nil
}
