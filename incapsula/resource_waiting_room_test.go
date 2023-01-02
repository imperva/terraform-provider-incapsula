package incapsula

import (
	"fmt"
	"log"
	"strconv"
	"testing"
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
		CheckDestroy: checkWaitingRoomDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccWaitingRoomBasic(t),
				Check: resource.ComposeTestCheckFunc(
					//testCheckAccwaitingRoomAfterFullUpdate(),
					resource.TestCheckResourceAttr(waitingRoomResource, "enabled", "true"),
					resource.TestCheckResourceAttr(waitingRoomResource, "entrance_rate_threshold", "100"),
					resource.TestCheckResourceAttr(waitingRoomResource, "concurrent_sessions_threshold", "150"),
				),
			},
			{
				ResourceName:      waitingRoomResource,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: getWaitingRoomImportString,
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
		concurrent_sessions_threshold = 150
	}`, waitingRoomResourceName, waitingRoomConfigName, rand.New(rand.NewSource(time.Now().UnixNano())).Intn(1000))
}

func checkWaitingRoomDestroy(state *terraform.State) error {
	client := testAccProvider.Meta().(*Client)

	for _, res := range state.RootModule().Resources {
		if res.Type != waitingRoomResourceName {
			continue
		}

		waitingRoomID := res.Primary.ID
		if waitingRoomID == "" {
			return fmt.Errorf("Incapsula Waiting Room ID does not exist")
		}

		siteID := res.Primary.Attributes["site_id"]
		if siteID == "" {
			return fmt.Errorf("Incapsula Waiting Room with id %s doesn't have site ID", waitingRoomID)
		}

		waitingRoomIdInt, err := strconv.ParseInt(waitingRoomID, 10, 64)
		if err != nil {
			return fmt.Errorf("Incapsula Waiting Room with id %s doesn't have numeric ID", waitingRoomID)
		}

		waitingRoomDTOResponse, _ := client.ReadWaitingRoom(siteID, waitingRoomIdInt)
		if waitingRoomDTOResponse == nil {
			return fmt.Errorf("Failed to check Waiting Room status (id=%s)", waitingRoomID)
		}
		if waitingRoomDTOResponse.Errors[0].Status != 404 {
			return fmt.Errorf("Incapsula Waiting Room with id %s still exists", waitingRoomID)
		}
	}

	return nil
}

func getWaitingRoomImportString(state *terraform.State) (string, error) {
	fmt.Println(state)
	fmt.Println(state.RootModule().Resources)
	for _, rs := range state.RootModule().Resources {
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
