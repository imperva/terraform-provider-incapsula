package incapsula

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

const siemGenericHttpConnectionResourceType = "incapsula_siem_generic_http_connection"

func resourceSiemGenericHttpConnection() *schema.Resource {
	return &schema.Resource{
		Create: resourceSiemGenericHttpConnectionCreate,
		Read:   resourceSiemGenericHttpConnectionRead,
		Update: resourceSiemGenericHttpConnectionUpdate,
		Delete: resourceSiemGenericHttpConnectionDelete,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				idSlice := strings.Split(d.Id(), "/")
				if len(idSlice) != 2 || idSlice[0] == "" || idSlice[1] == "" {
					return nil, fmt.Errorf("unexpected format of ID (%q), expected account_id/connection_id", d.Id())
				}

				d.Set("account_id", idSlice[0])
				d.SetId(idSlice[1])

				return []*schema.ResourceData{d}, nil
			},
		},

		Schema: map[string]*schema.Schema{
			"account_id": {
				Description: "Client account id.",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			"connection_name": {
				Description: "Name of the connection.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"url": {
				Description: "HTTPS endpoint the logs are sent to.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"headers": {
				Description: "Additional non-secret HTTP headers sent with every request.",
				Type:        schema.TypeMap,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"disable_cert_verification": {
				Description: "Disable SSL certificate verification of the endpoint.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
			"auth_type": {
				Description:  "Authentication type: HEADERS or BASIC_AUTH.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{HttpAuthTypeHeaders, HttpAuthTypeBasicAuth}, false),
			},
			"auth_header": {
				Description: "Secret authentication headers. Used with auth_type HEADERS.",
				Type:        schema.TypeList,
				Optional:    true,
				MinItems:    1,
				MaxItems:    5,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"header_name": {
							Description: "Name of the header.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"secret_format": {
							Description: "Header value template containing the {secret} placeholder.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"secret": {
							Description:      "Secret substituted into secret_format.",
							Type:             schema.TypeString,
							Required:         true,
							Sensitive:        true,
							DiffSuppressFunc: suppressSensitiveDataPlaceholderDiff,
						},
					},
				},
			},
			"username": {
				Description: "Username. Used with auth_type BASIC_AUTH.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"password": {
				Description:      "Password. Used with auth_type BASIC_AUTH.",
				Type:             schema.TypeString,
				Optional:         true,
				Sensitive:        true,
				DiffSuppressFunc: suppressSensitiveDataPlaceholderDiff,
			},
		},
	}
}

func suppressSensitiveDataPlaceholderDiff(k, old, new string, d *schema.ResourceData) bool {
	return new == sensitiveDataPlaceholder
}

func siemGenericHttpConnectionResourceValidation(d *schema.ResourceData) error {
	authHeaderCount := len(d.Get("auth_header").([]interface{}))
	hasUsername := d.Get("username").(string) != ""
	hasPassword := d.Get("password").(string) != ""

	switch d.Get("auth_type").(string) {
	case HttpAuthTypeHeaders:
		if authHeaderCount == 0 {
			return fmt.Errorf("[ERROR] at least one auth_header must be provided when auth_type is %s for %s", HttpAuthTypeHeaders, siemGenericHttpConnectionResourceType)
		}
		if hasUsername || hasPassword {
			return fmt.Errorf("[ERROR] username and password are mutually exclusive with auth_type %s for %s", HttpAuthTypeHeaders, siemGenericHttpConnectionResourceType)
		}
	case HttpAuthTypeBasicAuth:
		if authHeaderCount > 0 {
			return fmt.Errorf("[ERROR] auth_header and auth_type %s are mutually exclusive for %s", HttpAuthTypeBasicAuth, siemGenericHttpConnectionResourceType)
		}
		if !hasUsername || !hasPassword {
			return fmt.Errorf("[ERROR] username and password must be provided when auth_type is %s for %s", HttpAuthTypeBasicAuth, siemGenericHttpConnectionResourceType)
		}
	}
	return nil
}

func buildSiemGenericHttpConnectionData(d *schema.ResourceData) SiemConnectionData {
	return SiemConnectionData{
		ID:             d.Id(),
		AssetID:        d.Get("account_id").(string),
		ConnectionName: d.Get("connection_name").(string),
		StorageType:    StorageTypeCustomerHttp,
		ConnectionInfo: HttpConnectionInfo{
			Url:                     d.Get("url").(string),
			Authentication:          buildHttpAuthenticationInfo(d),
			Headers:                 buildHttpHeaders(d),
			DisableCertVerification: d.Get("disable_cert_verification").(bool),
		},
	}
}

func buildHttpAuthenticationInfo(d *schema.ResourceData) HttpAuthenticationInfo {
	authType := d.Get("auth_type").(string)
	if authType == HttpAuthTypeBasicAuth {
		return HttpAuthenticationInfo{
			Type:     authType,
			Username: d.Get("username").(string),
			Password: d.Get("password").(string),
		}
	}

	authHeaders := []AuthHeaderInfo{}
	for _, rawAuthHeader := range d.Get("auth_header").([]interface{}) {
		authHeader := rawAuthHeader.(map[string]interface{})
		authHeaders = append(authHeaders, AuthHeaderInfo{
			HeaderName:   authHeader["header_name"].(string),
			SecretFormat: authHeader["secret_format"].(string),
			Secret:       authHeader["secret"].(string),
		})
	}
	return HttpAuthenticationInfo{Type: authType, AuthHeaders: authHeaders}
}

func buildHttpHeaders(d *schema.ResourceData) map[string]string {
	headers := map[string]string{}
	for name, value := range d.Get("headers").(map[string]interface{}) {
		headers[name] = value.(string)
	}
	return headers
}

func resourceSiemGenericHttpConnectionCreate(d *schema.ResourceData, m interface{}) error {
	if err := siemGenericHttpConnectionResourceValidation(d); err != nil {
		return err
	}

	client := m.(*Client)
	response, statusCode, err := client.CreateSiemConnection(&SiemConnection{Data: []SiemConnectionData{buildSiemGenericHttpConnectionData(d)}})
	if err != nil {
		return err
	}

	if (*statusCode == 201) && (response != nil) && (len(response.Data) == 1) {
		d.SetId(response.Data[0].ID)
		return resourceSiemGenericHttpConnectionRead(d, m)
	}
	return fmt.Errorf("[ERROR] Unsupported operation. Response status code: %d", *statusCode)
}

func resourceSiemGenericHttpConnectionRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	response, statusCode, err := client.ReadSiemConnection(d.Id(), d.Get("account_id").(string))
	// The client reports any non-200 as an error, so 404 is checked first.
	if statusCode != nil && *statusCode == 404 {
		d.SetId("")
		return nil
	}
	if err != nil {
		return err
	}
	if (*statusCode != 200) || (response == nil) || (len(response.Data) != 1) {
		return fmt.Errorf("[ERROR] Unsupported operation. Response status code: %d", *statusCode)
	}

	connection := response.Data[0]
	connectionInfo, ok := connection.ConnectionInfo.(HttpConnectionInfo)
	if !ok {
		return fmt.Errorf("[ERROR] connection %s is of storage type %s, expected %s", connection.ID, connection.StorageType, StorageTypeCustomerHttp)
	}

	d.Set("account_id", connection.AssetID)
	d.Set("connection_name", connection.ConnectionName)
	d.Set("url", connectionInfo.Url)
	d.Set("headers", connectionInfo.Headers)
	d.Set("disable_cert_verification", connectionInfo.DisableCertVerification)
	d.Set("auth_type", connectionInfo.Authentication.Type)
	// Secrets come back as ciphertext, so the practitioner's values in state are kept as-is.
	if connectionInfo.Authentication.Type == HttpAuthTypeBasicAuth {
		d.Set("username", connectionInfo.Authentication.Username)
		d.Set("auth_header", []interface{}{})
	} else {
		d.Set("auth_header", flattenAuthHeadersKeepingSecrets(connectionInfo.Authentication.AuthHeaders, d.Get("auth_header").([]interface{})))
	}
	return nil
}

func flattenAuthHeadersKeepingSecrets(authHeaders []AuthHeaderInfo, stateAuthHeaders []interface{}) []interface{} {
	flattened := make([]interface{}, 0, len(authHeaders))
	for index, authHeader := range authHeaders {
		secret := ""
		if index < len(stateAuthHeaders) {
			if stateAuthHeader, ok := stateAuthHeaders[index].(map[string]interface{}); ok {
				secret, _ = stateAuthHeader["secret"].(string)
			}
		}
		flattened = append(flattened, map[string]interface{}{
			"header_name":   authHeader.HeaderName,
			"secret_format": authHeader.SecretFormat,
			"secret":        secret,
		})
	}
	return flattened
}

func resourceSiemGenericHttpConnectionUpdate(d *schema.ResourceData, m interface{}) error {
	if err := siemGenericHttpConnectionResourceValidation(d); err != nil {
		return err
	}

	client := m.(*Client)
	_, _, err := client.UpdateSiemConnection(&SiemConnection{Data: []SiemConnectionData{buildSiemGenericHttpConnectionData(d)}})
	if err != nil {
		return err
	}
	return resourceSiemGenericHttpConnectionRead(d, m)
}

func resourceSiemGenericHttpConnectionDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	if _, err := client.DeleteSiemConnection(d.Id(), d.Get("account_id").(string)); err != nil {
		return err
	}
	d.SetId("")
	return nil
}
