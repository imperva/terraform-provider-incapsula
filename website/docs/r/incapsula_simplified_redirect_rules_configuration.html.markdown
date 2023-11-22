# incapsula_simplified_redirect_rules_configuration
Provides the delivery simplified redirect rules configuration for a specific site. The order of rules execution (a.k.a. priority) is the same as the order they are defined in the resource configuration.
Currently there delivery SIMPLIFIED_REDIRECT rule types:

* **SIMPLIFIED_REDIRECT** - Redirect requests with 30X response. (this category doesn't support condition triggers, and needs to be enbled at the account level before being used)

## Example Usage

## `SIMPLIFIED_REDIRECT` RULES

```hcl
resource "incapsula_delivery_rules_configuration" "simplified-redirect-rules" {
  site_id = incapsula_site.example-site.id
  rule {
    from = "/1"
    to = "$scheme://www.example.com/$city"
    response_code = "302"
    rule_name = "New delivery simplified redirect rule",
    enabled = "true"
  }

  rule {
    ...
  }
}
```

### Argument Reference
* `site_id` - (Required) Numeric identifier of the site to operate on.
* `rule_name` - (Required) Rule name.
* `response_code` - (Required) Redirect status code. Valid values are `302`, `301`, `303`, `307`, `308`.
* `from` - (Required) URL pattern to rewrite. **Note**: this field must be unique among other rules of the same category.
* `to` - (Required) URL pattern to change to.
* `enabled` - (Optional) Boolean that enables the rule. Default value is true.

## Import

Delivery rules configuration can be imported using the site_id and category separated by /, e.g.:

```
$ terraform import delivery_rules_configuration.demo site_id/category
```