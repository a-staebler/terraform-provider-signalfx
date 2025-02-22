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

const newIntegrationAWSConfig = `
resource "signalfx_aws_integration" "aws_myteamXX" {
    name = "AWSFoo"
    enabled = false

		auth_method = "ExternalId"
		role_arn = "arn:aws:iam::XXX:role/SignalFx-Read-Role"
		regions = ["us-east-1"]
		poll_rate = 300
		import_cloud_watch = true
		enable_aws_usage = true

		custom_namespace_sync_rule {
			default_action = "Exclude"
			filter_action = "Include"
			filter_source = "filter('code', '200')"
			namespace = "fart"
		}

		namespace_sync_rule {
			default_action = "Exclude"
			filter_action = "Include"
			filter_source = "filter('code', '200')"
			namespace = "AWS/EC2"
		}
}
`

const updatedIntegrationAWSConfig = `
resource "signalfx_aws_integration" "aws_myteamXX" {
    name = "AWSFoo NEW"
    enabled = false

		auth_method = "ExternalId"
		role_arn = "arn:aws:iam::XXX:role/SignalFx-Read-Role"
		regions = ["us-east-1"]
		poll_rate = 300
		import_cloud_watch = true
		enable_aws_usage = true

		custom_namespace_sync_rule {
			default_action = "Exclude"
			filter_action = "Include"
			filter_source = "filter('code', '200')"
			namespace = "fart"
		}

		namespace_sync_rule {
			default_action = "Exclude"
			filter_action = "Include"
			filter_source = "filter('code', '200')"
			namespace = "AWS/EC2"
		}
}
`

func TestAccCreateUpdateIntegrationAWS(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccIntegrationAWSDestroy,
		Steps: []resource.TestStep{
			// Create It
			{
				Config: newIntegrationAWSConfig,
				Check:  testAccCheckIntegrationAWSResourceExists,
			},
			{
				ResourceName:      "signalfx_aws_integration.aws_myteamXX",
				ImportState:       true,
				ImportStateIdFunc: testAccStateIdFunc("signalfx_aws_integration.aws_myteamXX"),
				ImportStateVerify: true,
				// The API doesn't return this value, so blow it up
				ImportStateVerifyIgnore: []string{"role_arn"},
			},
			// Update It
			{
				Config: updatedIntegrationAWSConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIntegrationAWSResourceExists,
					resource.TestCheckResourceAttr("signalfx_aws_integration.aws_myteamXX", "name", "AWSFoo NEW"),
				),
			},
		},
	})
}

func testAccCheckIntegrationAWSResourceExists(s *terraform.State) error {
	client, _ := sfx.NewClient(os.Getenv("SFX_AUTH_TOKEN"))

	for _, rs := range s.RootModule().Resources {
		switch rs.Type {
		case "signalfx_aws_integration":
			integration, err := client.GetAWSCloudWatchIntegration(rs.Primary.ID)
			if integration == nil {
				return fmt.Errorf("Error finding integration %s: %s", rs.Primary.ID, err)
			}
		default:
			return fmt.Errorf("Unexpected resource of type: %s", rs.Type)
		}
	}

	return nil
}

func testAccIntegrationAWSDestroy(s *terraform.State) error {
	client, _ := sfx.NewClient(os.Getenv("SFX_AUTH_TOKEN"))
	for _, rs := range s.RootModule().Resources {
		switch rs.Type {
		case "signalfx_aws_integration":
			integration, _ := client.GetAWSCloudWatchIntegration(rs.Primary.ID)
			if integration != nil {
				return fmt.Errorf("Found deleted integration %s", rs.Primary.ID)
			}
		default:
			return fmt.Errorf("Unexpected resource of type: %s", rs.Type)
		}
	}

	return nil
}

func TestValidateAwsService(t *testing.T) {
	_, errors := validateAwsService("AWS/Logs", "")
	assert.Equal(t, 0, len(errors), "No errors for valid value")

	_, errors = validateAwsService("Fart", "")
	assert.Equal(t, 1, len(errors), "Errors for invalid value")
}

func TestValidateAuthMethod(t *testing.T) {
	_, errors := validateAuthMethod("ExternalId", "")
	assert.Equal(t, 0, len(errors), "No errors for valid value")

	_, errors = validateAuthMethod("SecurityToken", "")
	assert.Equal(t, 0, len(errors), "No errors for valid value")

	_, errors = validateAuthMethod("Fart", "")
	assert.Equal(t, 1, len(errors), "Errors for invalid value")
}

func TestValidateFilterAction(t *testing.T) {
	_, errors := validateFilterAction("Exclude", "")
	assert.Equal(t, 0, len(errors), "No errors for valid value")

	_, errors = validateFilterAction("Include", "")
	assert.Equal(t, 0, len(errors), "No errors for valid value")

	_, errors = validateFilterAction("Fart", "")
	assert.Equal(t, 1, len(errors), "Errors for invalid value")
}

func TestValidateAwsPollRate(t *testing.T) {
	_, errors := validateAwsPollRate(60, "")
	assert.Equal(t, 0, len(errors), "No errors for valid value")

	_, errors = validateAwsPollRate(300, "")
	assert.Equal(t, 0, len(errors), "No errors for valid value")

	_, errors = validateAwsPollRate(12, "")
	assert.Equal(t, 1, len(errors), "Errors for invalid value")
}
