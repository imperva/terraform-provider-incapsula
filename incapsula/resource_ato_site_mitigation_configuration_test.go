package incapsula

import (
	"testing"
)

const atoSiteMitigationConfigurationResourceType = "incapsula_ato_site_mitigation_configuration"
const atoSiteMitigationConfigurationResourceName = "testacc-terraform-ato-site-mitigation-configuration"
const atoSiteMitigationConfigurationResource = atoSiteMitigationConfigurationResourceType + "." + atoSiteMitigationConfigurationResourceName

/*
	The test endpoints do not exist. Currently, the onboarding process is done through the UI to ensure correctness

Upcoming: We will be adding the feature to copy endpoint configurations from other sites.
But since we want our tests to be site agnostic, we should not implement a specific site to copy the endpoint configuration from.
That would cause the entire tests to depend on a site which cannot be maintained by anyone who can contribute to this repository
When we enable the site onboarding feature to be done truly via APIs we can write the acceptance tests here.
Until then we have the client tests at resource_ato_site_mitigation_configuration_test.go
*/
func TestAccIncapsulaATOSiteMitigationConfiguration_basic(t *testing.T) {

}
