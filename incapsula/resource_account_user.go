package incapsula

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
	"net/mail"
	"strconv"
	"strings"
	"time"
)

const sleepTimeSeconds = 2

func resourceAccountUser() *schema.Resource {
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
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					email := val.(string)
					if _, err := mail.ParseAddress(email); err != nil {
						errs = append(errs, fmt.Errorf("%q is invalid, got: %s", key, email))
					}
					return
				},
			},
			"account_id": {
				Description: "Unique ID of the required account . example: 123456",
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
			},

			// Optional Arguments
			"first_name": {
				Description: "The first name of the user that was acted on. example: John",
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(string)
					if len(v) < 2 {
						errs = append(errs, fmt.Errorf("%q should have at least 2 characters, got: %s", key, v))
					}
					return
				},
			},
			"last_name": {
				Description: "The last name of the user that was acted on. example: Snow",
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(string)
					if len(v) < 2 {
						errs = append(errs, fmt.Errorf("%q should have at least 2 characters, got: %s", key, v))
					}
					return
				},
			},
			"role_ids": {
				Description: "List of role ids to add to the user.",
				Type:        schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Optional: true,
			},

			// Computed Arguments
			"role_names": {
				Description: "List of role names.",
				Type:        schema.TypeSet,
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

	roleIds := d.Get("role_ids").(*schema.Set)
	UserAddResponse, err := client.AddAccountUser(
		accountId,
		email,
		d.Get("first_name").(string),
		d.Get("last_name").(string),
		roleIds.List(),
	)

	if err != nil {
		log.Printf("[ERROR] Could not create user for email: %s, %s\n", email, err)
		return err
	}

	// Set the User ID
	d.SetId(fmt.Sprintf("%s/%s", strconv.Itoa(accountId), email))
	log.Printf("[INFO] Created Incapsula user for email: %s userid: %s\n", email, UserAddResponse.Data.UserID)

	// There may be a timing/race condition here
	// Set an arbitrary period to sleep
	log.Printf("[DEBUG] Avoid timing/race condition, sleeping %d seconds\n", sleepTimeSeconds)
	time.Sleep(sleepTimeSeconds * time.Second)

	// Set the rest of the state from the resource read
	return resourceUserRead(d, m)
}

func resourceUserRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	userID := d.Id()
	stringSlice := strings.Split(userID, "/")
	if len(stringSlice) != 2 {
		return fmt.Errorf("Error parsing ID, actual value: %s, expected numeric id and string seperated by '/'\n", stringSlice)
	}
	accountID, _ := strconv.Atoi(stringSlice[0])
	email := stringSlice[1]
	log.Printf("[INFO] Reading Incapsula user : %s\n", userID)

	userStatusResponse, err := client.GetAccountUser(accountID, email)

	if err != nil {
		log.Printf("[ERROR] Could not read Incapsula user: %s, %s\n", email, err)
		return err
	}

	log.Printf("[INFO]listRoles : %v\n", userStatusResponse.Data.Roles)

	listRolesIds := make([]int, len(userStatusResponse.Data.Roles))
	listRolesNames := make([]string, len(userStatusResponse.Data.Roles))
	for i, v := range userStatusResponse.Data.Roles {
		listRolesIds[i] = v.RoleID
		listRolesNames[i] = v.RoleName
	}

	d.Set("email", userStatusResponse.Data.Email)
	d.Set("account_id", userStatusResponse.Data.AccountID)
	d.Set("first_name", userStatusResponse.Data.FirstName)
	d.Set("last_name", userStatusResponse.Data.LastName)
	d.Set("role_ids", listRolesIds)
	d.Set("role_names", listRolesNames)

	log.Printf("[INFO] Finished reading Incapsula user: %s\n", email)

	return nil
}

func resourceUserUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	email := d.Get("email").(string)
	accountId := d.Get("account_id").(int)

	log.Printf("[INFO] Creating Incapsula user for email: %s\n", email)

	roleIds := d.Get("role_ids").(*schema.Set)
	userUpdateResponse, err := client.UpdateAccountUser(
		accountId,
		email,
		roleIds.List(),
	)
	if err != nil {
		log.Printf("[ERROR] Could not update user for email: %s, %s\n", email, err)
		return err
	}

	log.Printf("[Info] New Roles for user %s : %+v\n", email, userUpdateResponse.Data[0].Roles)

	// There may be a timing/race condition here
	// Set an arbitrary period to sleep
	log.Printf("[DEBUG] Avoid timing/race condition, sleeping %d seconds\n", sleepTimeSeconds)
	time.Sleep(sleepTimeSeconds * time.Second)

	// Set the rest of the state from the resource read
	return resourceUserRead(d, m)
}

func resourceUserDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	accountID := d.Get("account_id").(int)
	email := d.Get("email").(string)

	log.Printf("[INFO] Deleting Incapsula user: %s\n", email)

	err := client.DeleteAccountUser(accountID, email)

	if err != nil {
		log.Printf("[ERROR] Could not delete Incapsula user: %s %s\n", email, err)
		return err
	}

	// Set the ID to empty
	// Implicitly clears the resource
	d.SetId("")

	// There may be a timing/race condition here
	// Set an arbitrary period to sleep
	log.Printf("[DEBUG] Avoid timing/race condition, sleeping %d seconds\n", sleepTimeSeconds)
	time.Sleep(sleepTimeSeconds * time.Second)

	log.Printf("[INFO] Deleted Incapsula user: %s\n", email)

	return nil
}
