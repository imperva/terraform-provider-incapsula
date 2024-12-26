---
subcategory: "Application Performance and Delivery"
layout: "incapsula"
page_title: "incapsula_waiting_room"
description: |-
  Provides a Incapsula Waiting Room resource.
---

# incapsula_waiting_room

Provides a waiting room resource.
A waiting room controls the traffic to the website during peak periods when the origin server is unable to handle the load, and routes the website visitors to a virtual waiting room when their requests can't be handled immediately.
For full feature documentation, see [Set Up a Waiting Room](https://docs.imperva.com/bundle/cloud-application-security/page/waiting-room.htm).

**Note:** At least one of the threshold strategies (entrance_rate_threshold / concurrent_sessions_threshold) must be configured.

## Example Usage

```hcl
resource "incapsula_waiting_room" "example-waiting-room" {
    site_id = incapsula_site.example-site.id
    name = "Waiting room name"
    description = "Waiting room description"
    enabled = true
    html_template_base64 = filebase64("${"path/to/your/template.html"}")
    filter = <<EOF
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

* `account_id` - (Optional) The account to operate on. If not specified, operation will be performed on the account identified by the authentication parameters.

* `site_id` - (Required) Numeric identifier of the site to operate on.

* `name` - (Required) The waiting room name. Must be unique across all waiting rooms of the site.

* `description` - (Optional) The waiting room description.

* `enabled` - (Optional) Indicates if this waiting room is enabled or not. **Default:** true.

* `html_template_base64` - (Optional) The HTML template file path in Base64 format. A default template is used in case one isn't provided. The following placeholders can be used to insert dynamic information:
  * `$WAITING_ROOM_CONFIG$` - Calls a script that periodically updates the status of the user, and reloads the page when the user is allowed to enter the website from the waiting room. This parameter is mandatory and should not be modified or deleted.
  * `$WAITING_ROOM_LOADER$` - Used to validate the loading of the page. This parameter is mandatory and should not be modified or deleted.
  * `$WAITING_ROOM_WRAPPER$` - Used to validate the content of the template. This parameter is mandatory and should not be modified or deleted.
  * `$WAITING_ROOM_POSITION_IN_LINE$` - Used to display the user's position in the waiting room queue.
  * `$WAITING_ROOM_LAST_STATUS_UPDATE$` - Used to display the time of the last status update.
  * `$ESTIMATED_TIME_TO_WAIT$` - Estimated time to wait.

* `filter` - (Optional) The conditions that determine on which sessions this waiting room applies. For example, you can create a condition to apply the waiting room to a subset of your website, instead of to the entire website, such as: **URL contains "^/ShoppingCart"**. You can also use conditions to create waiting rooms for different visitor groups, such as visitors from different countries. For example, **CountryCode == GB**. See [Rule Filter Parameters](https://docs.imperva.com/bundle/cloud-application-security/page/rules/rule-parameters.htm) for more filtering options. **Default:** No filter. The room applies to the entire website and all users.

* `bots_action_in_queuing_mode` - (Optional) The waiting room bot handling action. Determines the waiting room behavior for legitimate bots trying to access your website during peak time. Applies only when the activation threshold has been passed and visitors are being sent to the queue. **Default:** `WAIT_IN_LINE`
Possible values:
  * `WAIT_IN_LINE` - Wait in line alongside regular users.
  * `BYPASS` - Bypass the queue.
  * `BLOCK` - Block this request.

* `entrance_rate_threshold` - (Optional) The waiting room is activated when new users per minute exceed the specified value. Minimum of 60 users per minute.

* `concurrent_sessions_threshold` - (Optional) The waiting room is activated when the number of active users reaches the specified value. Minimum value = 1.

* `inactivity_timeout` - (Optional, Mandatory if concurrent_sessions_threshold is used) Inactivity timeout, from 1 to 30 minutes. If waiting room conditions that limit the scope of the waiting room to a subset of the website have been defined, the user is considered active only when navigating the pages in scope of the conditions. A user who is inactive for a longer period of time is considered as having left the site. On returning to the site, the user needs to wait in line again if the waiting room is active. **Tip:** When enabling the concurrent_sessions_threshold, the inactivity timeout is very important. Once the site is at full capacity (the threshold has been passed), no new user can access the site until another user leaves and frees up space. To optimize the user experience, we recommend setting a balanced inactivity timeout value — long enough so that the user's session is still open if they return quickly, but not so long that it unnecessarily prevents access to other waiting visitors. The default timeout of 5 minutes is the recommended minimum value. **Default:** 5 minutes.

* `queue_inactivity_timeout` - (Optional) Queue inactivity timeout, from 1 to 10 minutes. A user in the waiting room who is inactive for a longer period of time is considered as having left the queue. On returning to the site, the user moves to the end of the queue and needs to wait in line again if the waiting room is active. **Default:** 1 minute.

## Attributes Reference

The following attributes are exported:

* `id` - Unique identifier in the API for the waiting room.
* `created_at` - (timestamp) When the waiting room was created.
* `modified_at` - (timestamp) When the waiting room was last modified.
* `last_modified_by` - (user email) Last user modifying the waiting room.
* `mode` - (QUEUING or NOT_QUEUING) Waiting room current mode.

## Import

Waiting rooms can be imported using the waiting room account_id, site_id and waiting_room_id separated by /, e.g.:

```
$ terraform import incapsula_waiting_room.example-waiting-room account_id/site_id/waiting_room_id
```