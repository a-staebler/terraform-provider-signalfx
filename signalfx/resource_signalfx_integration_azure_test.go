package signalfx

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"

	"github.com/stretchr/testify/assert"

	sfx "github.com/signalfx/signalfx-go"
)

const newIntegrationAzureConfig = `
resource "signalfx_azure_integration" "azure_myteamXX" {
    name = "AzureFoo"
    enabled = false

    environment = "azure"

		poll_rate = 300

    secret_key = "XXX"

    app_id = "YYY"

    tenant_id = "ZZZ"

    services = [ "microsoft.sql/servers/elasticpools" ]

    subscriptions = [ "microsoft.sql/servers/elasticpools" ]
}
`

const updatedIntegrationAzureConfig = `
resource "signalfx_azure_integration" "azure_myteamXX" {
    name = "AzureFoo NEW"
    enabled = false

    environment = "azure"

		poll_rate = 300

    secret_key = "XXX"

    app_id = "YYY"

    tenant_id = "ZZZ"

    services = [ "microsoft.sql/servers/elasticpools" ]

    subscriptions = [ "microsoft.sql/servers/elasticpools" ]
}
`

func TestAccCreateUpdateIntegrationAzure(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccIntegrationAzureDestroy,
		Steps: []resource.TestStep{
			// Create It
			{
				Config: newIntegrationAzureConfig,
				Check:  testAccCheckIntegrationAzureResourceExists,
			},
			{
				ResourceName:      "signalfx_azure_integration.azure_myteamXX",
				ImportState:       true,
				ImportStateIdFunc: testAccStateIdFunc("signalfx_azure_integration.azure_myteamXX"),
				ImportStateVerify: true,
				// The API doesn't return this value, so blow it up
				ImportStateVerifyIgnore: []string{"app_id", "secret_key"},
			},
			// Update It
			{
				Config: updatedIntegrationAzureConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIntegrationAzureResourceExists,
					resource.TestCheckResourceAttr("signalfx_azure_integration.azure_myteamXX", "name", "AzureFoo NEW"),
				),
			},
		},
	})
}

func testAccCheckIntegrationAzureResourceExists(s *terraform.State) error {
	client, _ := sfx.NewClient(os.Getenv("SFX_AUTH_TOKEN"))

	for _, rs := range s.RootModule().Resources {
		switch rs.Type {
		case "signalfx_azure_integration":
			integration, err := client.GetAzureIntegration(rs.Primary.ID)
			if integration == nil {
				return fmt.Errorf("Error finding integration %s: %s", rs.Primary.ID, err)
			}
		default:
			return fmt.Errorf("Unexpected resource of type: %s", rs.Type)
		}
	}

	return nil
}

func testAccIntegrationAzureDestroy(s *terraform.State) error {
	client, _ := sfx.NewClient(os.Getenv("SFX_AUTH_TOKEN"))
	for _, rs := range s.RootModule().Resources {
		switch rs.Type {
		case "signalfx_azure_integration":
			integration, _ := client.GetAzureIntegration(rs.Primary.ID)
			if integration != nil {
				return fmt.Errorf("Found deleted integration %s", rs.Primary.ID)
			}
		default:
			return fmt.Errorf("Unexpected resource of type: %s", rs.Type)
		}
	}

	return nil
}

func TestValidateAzureService(t *testing.T) {
	_, errors := validateAzureService("microsoft.batch/batchaccounts", "")
	assert.Equal(t, 0, len(errors), "No errors for valid value")

	_, errors = validateAzureService("Fart", "")
	assert.Equal(t, 1, len(errors), "Errors for invalid value")
}

func TestValidateAzureEnvironment(t *testing.T) {
	_, errors := validateAzureEnvironment("azure", "")
	assert.Equal(t, 0, len(errors), "No errors for valid value")

	_, errors = validateAzureEnvironment("azure_us_government", "")
	assert.Equal(t, 0, len(errors), "No errors for valid value")

	_, errors = validateAzureEnvironment("Fart", "")
	assert.Equal(t, 1, len(errors), "Errors for invalid value")
}

func TestValidateAzurePollRate(t *testing.T) {
	_, errors := validateAzurePollRate(60, "")
	assert.Equal(t, 0, len(errors), "No errors for valid value")

	_, errors = validateAzurePollRate(300, "")
	assert.Equal(t, 0, len(errors), "No errors for valid value")

	_, errors = validateAzurePollRate(12, "")
	assert.Equal(t, 1, len(errors), "Errors for invalid value")
}
