package incapsula

import (
	"context"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceRoleAbilities() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRoleAbilitiesRead,

		Schema: map[string]*schema.Schema{
			// Computed Attributes
			"can_add_site": {
				Description: "Add sites",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"can_edit_site": {
				Description: "Modify site settings",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"can_edit_account": {
				Description: "Edit account settings",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"can_add_user": {
				Description: "Manage users",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"can_manage_api_key": {
				Description: "Manage API keys",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"can_manage_account_sub_accounts": {
				Description: "Manage account sub-accounts",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"can_edit_domain": {
				Description: "Modify DNS zone settings",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"can_add_domain": {
				Description: "Add DNS zones",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"can_view_infra_protect_setting": {
				Description: "View Infra Protect settings",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"can_run_connectivity_reports": {
				Description: "Allow user to run connectivity reports",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"can_purge_cache": {
				Description: "Purge cache",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"can_edit_single_ip": {
				Description: "Edit single IP",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"can_edit_roles": {
				Description: "Manage users roles",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"can_view_audit_trail": {
				Description: "View audit trail",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"can_view_client_certificates": {
				Description: "View client CA certificates",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"can_view_policy": {
				Description: "View policy",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"can_assign_client_certificates": {
				Description: "Manage client CA certificates for site",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"can_delete_policy_exception": {
				Description: "Delete exception from policy",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"can_delete_policy": {
				Description: "Delete policy",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"can_add_policy": {
				Description: "Add/Duplicate policy",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"can_edit_client_certificates": {
				Description: "Manage client CA certificates for account",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"can_edit_policy": {
				Description: "Edit policy",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"can_edit_policy_exception": {
				Description: "Edit exception in policy",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"can_apply_policy_to_assets": {
				Description: "Apply policy to assets",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"can_add_policy_exception": {
				Description: "Add exception to Policy",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func dataSourceRoleAbilitiesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Generate ID
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	// Populate values (doing this to improve TF usability)
	d.Set("can_add_site", "canAddSite")
	d.Set("can_edit_site", "canEditSite")
	d.Set("can_edit_account", "canEditAccount")
	d.Set("can_add_user", "canAddUser")
	d.Set("can_manage_api_key", "canManageApiKey")
	d.Set("can_manage_account_sub_accounts", "canManageAccountSubAccounts")
	d.Set("can_edit_domain", "canEditDomain")
	d.Set("can_add_domain", "canAddDomain")
	d.Set("can_view_infra_protect_setting", "canViewInfraProtectSetting")
	d.Set("can_run_connectivity_reports", "canRunConnectivityReports")
	d.Set("can_purge_cache", "canPurgeCache")
	d.Set("can_edit_single_ip", "canEditcanEditSingleIpSite")
	d.Set("can_edit_roles", "canEditRoles")
	d.Set("can_view_audit_trail", "canViewAuditTrail")
	d.Set("can_view_client_certificates", "canViewClientCertificates")
	d.Set("can_view_policy", "canViewPolicy")
	d.Set("can_assign_client_certificates", "canAssignClientCertificates")
	d.Set("can_delete_policy_exception", "canDeletePolicyException")
	d.Set("can_delete_policy", "canDeletePolicy")
	d.Set("can_add_policy", "canAddPolicy")
	d.Set("can_edit_client_certificates", "canEditClientCertificates")
	d.Set("can_edit_policy", "canEditPolicy")
	d.Set("can_edit_policy_exception", "canEditPolicyException")
	d.Set("can_apply_policy_to_assets", "canApplyPolicyToAssets")
	d.Set("can_add_policy_exception", "canAddPolicyException")

	return nil
}
