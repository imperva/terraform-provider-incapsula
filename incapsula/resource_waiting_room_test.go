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
	var waitingRoomDTOResponse WaitingRoomDTO

	log.Printf("========================BEGIN TEST========================")
	log.Printf("[DEBUG]Running test resource_waiting_room_test.TestAccWaitingRoom_Basic")
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckWaitingRoomDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccWaitingRoomBasic(t, 100, 200),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWaitingRoomExists(&waitingRoomDTOResponse),
					testAccCheckWaitingRoomThresholds(&waitingRoomDTOResponse, 100, 200),
					resource.TestCheckResourceAttr(waitingRoomResource, "enabled", "true"),
					resource.TestCheckResourceAttr(waitingRoomResource, "entrance_rate_threshold", "100"),
					resource.TestCheckResourceAttr(waitingRoomResource, "concurrent_sessions_threshold", "200"),
				),
			},
			{
				Config: testAccWaitingRoomEntranceRateOnly(t, 70),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWaitingRoomExists(&waitingRoomDTOResponse),
					testAccCheckWaitingRoomThresholds(&waitingRoomDTOResponse, 70, 0),
					resource.TestCheckResourceAttr(waitingRoomResource, "enabled", "true"),
					resource.TestCheckResourceAttr(waitingRoomResource, "entrance_rate_threshold", "70"),
					resource.TestCheckNoResourceAttr(waitingRoomResource, "concurrent_sessions_threshold"),
				),
			},
			{
				Config: testAccWaitingRoomConcurrentSessionsOnly(t, 50),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWaitingRoomExists(&waitingRoomDTOResponse),
					testAccCheckWaitingRoomThresholds(&waitingRoomDTOResponse, 0, 50),
					resource.TestCheckResourceAttr(waitingRoomResource, "enabled", "true"),
					resource.TestCheckNoResourceAttr(waitingRoomResource, "entrance_rate_threshold"),
					resource.TestCheckResourceAttr(waitingRoomResource, "concurrent_sessions_threshold", "50"),
				),
			},
			{
				ResourceName:      waitingRoomResource,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccGetWaitingRoomImportString,
			},
		},
	})
}

func testAccCheckWaitingRoomExists(waitingRoomDTOresponse *WaitingRoomDTO) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[waitingRoomResource]
		if !ok {
			return fmt.Errorf("Not found: %s", waitingRoomResource)
		}

		waitingRoomID, err := strconv.ParseInt(rs.Primary.ID, 10, 64)
		if err != nil {
			return fmt.Errorf("Error parsing ID %s to int", rs.Primary.ID)
		}

		siteID := rs.Primary.Attributes["site_id"]
		if siteID == "" {
			return fmt.Errorf("Incapsula Waiting Room with id %d doesn't have site ID", waitingRoomID)
		}

		client := testAccProvider.Meta().(*Client)

		response, _ := client.ReadWaitingRoom(siteID, waitingRoomID)
		if response == nil {
			return fmt.Errorf("Failed to retrieve Waiting Room (id=%d)", waitingRoomID)
		}

		if response.Errors == nil && response.Data != nil {
			*waitingRoomDTOresponse = response.Data[0]
			return nil
		}
		return fmt.Errorf("resource %s was not updated correctly after full update", waitingRoomResource)
	}
}

func testAccWaitingRoomBasic(t *testing.T, entranceRate int, concurrentSessions int) string {
	return testAccCheckIncapsulaSiteConfigBasic(GenerateTestDomain(t)) + fmt.Sprintf(`
	resource "%s" "%s" {
		site_id = incapsula_site.testacc-terraform-site.id
		name = "testWaitingRoom%d"
		enabled = true
		entrance_rate_threshold = %d
		concurrent_sessions_threshold = %d
	}`, waitingRoomResourceName, waitingRoomConfigName, rand.New(rand.NewSource(time.Now().UnixNano())).Intn(1000), entranceRate, concurrentSessions)
}

func testAccWaitingRoomEntranceRateOnly(t *testing.T, entranceRate int) string {
	return testAccCheckIncapsulaSiteConfigBasic(GenerateTestDomain(t)) + fmt.Sprintf(`
	resource "%s" "%s" {
		site_id = incapsula_site.testacc-terraform-site.id
		name = "testWaitingRoom%d"
		enabled = true
		entrance_rate_threshold = %d
	}`, waitingRoomResourceName, waitingRoomConfigName, rand.New(rand.NewSource(time.Now().UnixNano())).Intn(1000), entranceRate)
}

func testAccWaitingRoomConcurrentSessionsOnly(t *testing.T, concurrentSessions int) string {
	return testAccCheckIncapsulaSiteConfigBasic(GenerateTestDomain(t)) + fmt.Sprintf(`
	resource "%s" "%s" {
		site_id = incapsula_site.testacc-terraform-site.id
		name = "testWaitingRoom%d"
		enabled = true
		concurrent_sessions_threshold = %d
	}`, waitingRoomResourceName, waitingRoomConfigName, rand.New(rand.NewSource(time.Now().UnixNano())).Intn(1000), concurrentSessions)
}

func testAccCheckWaitingRoomThresholds(waitingRoomDTOresponse *WaitingRoomDTO, entranceRate int, concurrentSessions int) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		if entranceRate == 0 {
			if waitingRoomDTOresponse.EntranceRateEnabled != false || waitingRoomDTOresponse.EntranceRateThreshold != 0 {
				return fmt.Errorf("Expected disabled entrance threshold, found %d:", waitingRoomDTOresponse.EntranceRateThreshold)
			}
		} else if waitingRoomDTOresponse.EntranceRateEnabled == false || waitingRoomDTOresponse.EntranceRateThreshold == 0 {
			return fmt.Errorf("Entrance rate disabled (should hav been %d)", entranceRate)
		}

		if concurrentSessions == 0 {
			if waitingRoomDTOresponse.ConcurrentSessionsEnabled != false || waitingRoomDTOresponse.ConcurrentSessionsThreshold != 0 {
				return fmt.Errorf("Expected disabled concurrent sessions, found %d:", waitingRoomDTOresponse.ConcurrentSessionsThreshold)
			}
		} else if waitingRoomDTOresponse.ConcurrentSessionsEnabled == false || waitingRoomDTOresponse.ConcurrentSessionsThreshold == 0 {
			return fmt.Errorf("Concurrent sessions disabled (should hav been %d)", concurrentSessions)
		}

		return nil
	}
}

func testAccCheckWaitingRoomDestroy(state *terraform.State) error {
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

func testAccGetWaitingRoomImportString(state *terraform.State) (string, error) {
	fmt.Println(state)
	fmt.Println(state.RootModule().Resources)
	for _, rs := range state.RootModule().Resources {
		if rs.Type != "incapsula_waiting_room" {
			continue
		}

		waitingRoomID, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return "", fmt.Errorf("Error parsing ID %s to int", rs.Primary.ID)
		}
		siteID, err := strconv.Atoi(rs.Primary.Attributes["site_id"])
		if err != nil {
			return "", fmt.Errorf("Error parsing site_id %s to int", rs.Primary.Attributes["site_id"])
		}
		return fmt.Sprintf("%d/%d", siteID, waitingRoomID), nil
	}

	return "", fmt.Errorf("Error finding Waiting Room")
}
