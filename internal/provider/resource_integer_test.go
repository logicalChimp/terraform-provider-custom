package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccResourceIntegerBasic(t *testing.T) {
	t.Parallel()
	resource.UnitTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testRandomIntegerBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccResourceIntegerBasic("sequential_integer.integer_1"),
				),
			},
			{
				ResourceName:      "sequential_integer.integer_1",
				ImportState:       true,
				ImportStateId:     "3,1,3",
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccResourceIntegerUpdate(t *testing.T) {
	t.Parallel()
	resource.UnitTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testRandomIntegerBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccResourceIntegerBasic("sequential_integer.integer_1"),
				),
			},
			{
				Config: testRandomIntegerBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccResourceIntegerUpdate("sequential_integer.integer_1"),
				),
			},
		},
	})
}

func TestAccResourceIntegerBig(t *testing.T) {
	t.Parallel()
	resource.UnitTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testRandomIntegerBig,
			},
			{
				ResourceName:      "sequential_integer.integer_1",
				ImportState:       true,
				ImportStateId:     "7227701560655103598,7227701560655103597,7227701560655103598",
				ImportStateVerify: true,
			},
		},
	})
}

func testAccResourceIntegerBasic(id string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[id]
		if !ok {
			return fmt.Errorf("Not found: %s", id)
		}
		result := rs.Primary.Attributes["result"]

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		if result == "" {
			return fmt.Errorf("Result not found")
		}

		if result != "1" {
			return fmt.Errorf("Invalid result %s. Did not initialise to 'min' value", result)
		}
		return nil
	}
}

func testAccResourceIntegerUpdate(id string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, _ := s.RootModule().Resources[id]
		// if !ok {
		// 	return fmt.Errorf("Not found: %s", id)
		// }
		result := rs.Primary.Attributes["result"]

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		if result == "" {
			return fmt.Errorf("Result not found")
		}

		if result != "2" {
			return fmt.Errorf("Invalid result %s. Did not update sequentially", result)
		}
		return nil
	}
}

const (
	testRandomIntegerBasic = `
resource "sequential_integer" "integer_1" {
   min  = 1
   max  = 3
}
`

	testRandomIntegerBig = `
resource "sequential_integer" "integer_1" {
   max  = 7227701560655103598
   min  = 7227701560655103597
}`
)
