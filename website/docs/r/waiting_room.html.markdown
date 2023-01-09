---
layout: "incapsula"
page_title: "Incapsula: waiting-room"
sidebar_current: "docs-incapsula-waiting-room"
description: |-
  Provides a Incapsula Waiting Room resource.
---

# incapsula_waiting_room

Provides a waiting room resource.
A waiting room controls the traffic to the website during peak periods when the origin server is unable to handle the load, and route the website visitors to a virtual waiting room when their requests can't be handled immediately.

**Note:** at least one of the threshold strategies (entrance_rate_threshold / concurrent_sessions_threshold) must be configured.

## Example Usage

```hcl
resource "incapsula_waiting_room" "example-waiting-room" {
    site_id = incapsula_site.example-site.id
    name = "Waiting room name"
    description = "Waiting room description"
    enabled = true
    html_template_base64 = filebase64("${"path/to/your/template.html"}")
    filter = >>>EOF
        URL == "/example"
    EOF
    bots_action_in_queuing_mode = "WAIT_IN_LINE"
    entrance_rate_threshold = 600
    concurrent_sessions_threshold = 600
    inactivity_timeout = 30
    queue_inactivity_timeout = 1
}
```

## Argument Reference

The following arguments are supported:

* `site_id` - (Required) Numeric identifier of the site to operate on.
* `name` - (Required) The waiting room name. Must be unique across all waiting room of the site.
* `description` - (Optional) The waiting room description.
* `enabled` - (Optional) whether this waiting room is enabled or not. **default:** true.
* `html_template_base64` - (Optional) The HTML template file path. A default Incapsula template is used in case one isn't provided. The following placeholders can be used to insert dynamic information:
* `filter` - (Optional) The condition that determines on which sessions this waiting room applies. **default:** no filter (i.e. the room applies to the whole website and all users)
* `bots_action_in_queuing_mode` - (Optional) The waiting room bot handling action. Determines the waiting room behavior for legitimate bots trying to access your website during peak time. Applies only when the activation threshold has been passed and visitors are being sent to the queue. Possible values:
- `WAIT_IN_LINE` - Wait in line alongside regular users.
- `BYPASS` - Bypass the queue.
- `BLOCK` - Block this request.
**default:** `WAIT_IN_LINE`
* `entrance_rate_threshold` - (Optional) The entrance rate activation threshold of the waiting room. The waiting room is activated when sessions per minute exceed the specified value. Minimum of 60 users per minute.
* `concurrent_sessions_threshold` - (Optional) The active users activation threshold of the waiting room. The waiting room is activated when number of active users reached specified value. Must be a positive number.
* `inactivity_timeout` - (Optional, Mandatory if concurrentSessionsThreshold is used) Inactivity timeout, from 1 to 30 minutes. If waiting room conditions that limit the scope of the waiting room to a subset of the website have been defined, the user is considered active only when navigating the pages in scope of the conditions. A user who is inactive for a longer period of time is considered as having left the site. On returning to the site, the user needs to wait in line again if the waiting room is active. **Default:** 5 minutes.
* `queue_inactivity_timeout` - (Optional) Queue inactivity timeout, from 1 to 10 minutes. A user in the waiting room who is inactive for a longer period of time is considered as having left the queue. On returning to the site, the user moves to the end of the queue and needs to wait in line again if the waiting room is active. **default:** 1 minute.

## Attributes Reference

The following attributes are exported:

* `id` - Unique identifier in the API for the Waiting Room.
* `account_id` - The ID of the account that created the waiting room.
* `created_at` - (timestamp) When the waiting room was created.
* `modified_at` - (timestamp) When the waiting room was last modified.
* `last_modified_by` - (mail) Last user modifying the waiting room.
* `mode` - (QUEUING or NOT_QUEUING) Waiting room current mode.

## Import

Waiting rooms can be imported using the waiting room site_id and waiting_room_id separated by /, e.g.:

```
$ terraform import incapsula_waiting_room.demo site_id/waiting_room_id
```