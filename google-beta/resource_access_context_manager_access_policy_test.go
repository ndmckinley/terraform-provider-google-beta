package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

// Access Context Manager tests need to run serially
func TestAccAccessContextManagerAccessPolicy_basic(t *testing.T) {
	org := getTestOrgFromEnv(t)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAccessContextManagerAccessPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAccessContextManagerAccessPolicy_basic(org, "my policy"),
			},
			{
				ResourceName:      "google_access_context_manager_access_policy.test-access",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAccessContextManagerAccessPolicy_basic(org, "my new policy"),
			},
			{
				ResourceName:      "google_access_context_manager_access_policy.test-access",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckAccessContextManagerAccessPolicyDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "google_access_context_manager_access_policy" {
			continue
		}

		config := testAccProvider.Meta().(*Config)

		url, err := replaceVarsForTest(rs, "https://accesscontextmanager.googleapis.com/v1beta/accessPolicies/{{name}}")
		if err != nil {
			return err
		}

		_, err = sendRequest(config, "GET", url, nil)
		if err == nil {
			return fmt.Errorf("AccessPolicy still exists at %s", url)
		}
	}

	return nil
}

func testAccAccessContextManagerAccessPolicy_basic(org, title string) string {
	return fmt.Sprintf(`
resource "google_access_context_manager_access_policy" "test-access" {
  parent = "organizations/%s"
  title  = "%s"
}
`, org, title)
}