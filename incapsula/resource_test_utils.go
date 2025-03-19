package incapsula

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// extracts the primary attributes of a resource from the state.
func extractAttributes(s *terraform.State) map[string]map[string]string {
	attributes := make(map[string]map[string]string)
	for name, rs := range s.RootModule().Resources {
		attributes[name] = rs.Primary.Attributes
	}
	return attributes
}

func printState() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		attrs := extractAttributes(s)
		for resourceName, values := range attrs {
			fmt.Printf("[INFO] Resource: %s\nAttributes: %+v\n", resourceName, values)
		}
		return nil
	}
}
