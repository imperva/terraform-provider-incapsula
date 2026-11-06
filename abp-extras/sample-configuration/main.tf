terraform {
  required_providers {
    incapsula = {
      source = "imperva/incapsula"
    }
  }
}

provider "incapsula" {}

// Look up account ID corresponding to the provided API key.
// `id` and `account_id` are the same value, so either can be referenced.
data "incapsula_abp_account" "current" {}

# All configuration put in a separate module so it can be referenced as
# a dependency for the publishing part. When any resource in the module
# changes its state publishing will be triggered
module "abp" {
  source     = "./abp"
  account_id = data.incapsula_abp_account.current.id
}

# `depends_on` defines what triggers publishing: a change to anything listed
# here creates a new preflight, which is then published. It can be a module, as
# here, or a plain list of resources.
#
# A resource left out of this list does not trigger publishing on its own, but
# publishing activates the whole account, so it is included in the next publish
# anything else triggers.
data "incapsula_abp_pending_changes" "current" {
  depends_on = [module.abp]
}

# Create a preflight (a snapshot of the configuration)
#
# `pending_hash` is optional, but it should always come from the data source
# above: it both detects the changes and orders the preflight after them.
#
# A snapshot that can no longer be published — because it was discarded, or
# because the account was published from elsewhere — is detected and replaced
# automatically on the next apply. A snapshot that is only outdated is kept, and
# the reason is reported in `warnings`.
resource "incapsula_abp_preflight" "current" {
  account_id   = data.incapsula_abp_account.current.id
  pending_hash = data.incapsula_abp_pending_changes.current.hash
}

# Publish the specific preflight (make the corresponding configuration active).
#
# Publishing covers the whole account and is asynchronous: the apply returns
# once it has been accepted, not once it has finished.
#
# `terraform destroy` deletes the objects but publishes nothing, so publish the
# account again afterwards to activate the teardown.
resource "incapsula_abp_publish" "publish" {
  preflight_id = incapsula_abp_preflight.current.id
}

output "encrypted_secret" {
  value = module.abp.my_credential_encrypted_secret
}
