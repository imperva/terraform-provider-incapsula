package incapsula

import (
	"fmt"
	"testing"
	"log"
	"strconv"
	"time"

	"math/rand"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const waitingRoomResourceName = "incapsula_waiting_room"
const waitingRoomConfigName = "testacc-terraform-waiting-room"
const waitingRoomResource = waitingRoomResourceName + "." + waitingRoomConfigName

func TestAccWaitingRoom_Basic(t *testing.T) {
	log.Printf("========================BEGIN TEST========================")
	log.Printf("[DEBUG]Running test resource_waiting_room_test.TestAccWaitingRoom_Basic")
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccWaitingRoomDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccWaitingRoomBasic(t),
				Check: resource.ComposeTestCheckFunc(
					//testCheckAccwaitingRoomAfterFullUpdate(),
					resource.TestCheckResourceAttr(waitingRoomResource, "enabled", "true"),
				),
			},
			{
				ResourceName:      waitingRoomResource,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccStateWaitingRoomID,
			},
		},
	})
}

func testAccWaitingRoomBasic(t *testing.T) string {
	return testAccCheckIncapsulaSiteConfigBasic(GenerateTestDomain(t)) + fmt.Sprintf(`
	resource "%s" "%s" {
		site_id = incapsula_site.testacc-terraform-site.id
		name = "testWaitingRoom%d"
		enabled = true
		entrance_rate_threshold = 100
	}`, waitingRoomResourceName, waitingRoomConfigName, rand.New(rand.NewSource(time.Now().UnixNano())).Intn(1000))
}

func testAccWaitingRoomDestroy(s *terraform.State) error {
	return nil
}

func testAccStateWaitingRoomID(s *terraform.State) (string, error) {
	fmt.Println(s)
	fmt.Println(s.RootModule().Resources)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "incapsula_waiting_room" {
			continue
		}

		waitingRoomID, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return "", fmt.Errorf("Error parsing ID %v to int", rs.Primary.ID)
		}
		siteID, err := strconv.Atoi(rs.Primary.Attributes["site_id"])
		if err != nil {
			return "", fmt.Errorf("Error parsing site_id %v to int", rs.Primary.Attributes["site_id"])
		}
		return fmt.Sprintf("%d/%d", siteID, waitingRoomID), nil
	}

	return "", fmt.Errorf("Error finding Waiting Room")
}