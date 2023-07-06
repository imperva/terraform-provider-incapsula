package incapsula

import (
	"testing"
)

const atoEndpointMitigationConfigurationResourceType = "incapsula_ato_endpoint_mitigation_configuration"
const atoEndpointMitigationConfigurationResourceName = "testacc-terraform-ato-endpoint-mitigation-configuration"
const atoSiteMitigationConfigurationResource = atoEndpointMitigationConfigurationResourceType + "." + atoEndpointMitigationConfigurationResourceName

/*
	The test endpoints do not exist. Currently, the onboarding process is done through the UI to ensure correctness

Upcoming: We will be adding the feature to copy endpoint configurations from other sites.
But since we want our tests to be site agnostic, we should not implement a specific site to copy the endpoint configuration from.
That would cause the entire tests to depend on a site which cannot be maintained by anyone who can contribute to this repository
When we enable the site onboarding feature to be done truly via APIs we can write the acceptance tests here.
Until then we have the client tests at resource_ato_endpoint_mitigation_configuration_test.go
*/
func TestAccIncapsulaATOEndpointMitigationConfiguration_basic(t *testing.T) {

}
