---
layout: "signalfx"
page_title: "SignalFx: signalfx_resource"
sidebar_current: "docs-signalfx-resource-dashboard-group"
description: |-
  Allows Terraform to create and manage SignalFx Dashboard Groups
---

# Resource: signalfx_dashboard_group

In the SignalFx web UI, a [dashboard group](https://developers.signalfx.com/v2/docs/dashboard-group-model) is a collection of dashboards.

**NOTE:** Dashboard groups cannot be accessed directly, but just via a dashboard contained in them. This is the reason why make show won't show any of yours dashboard groups.

## Example Usage

```terraform
resource "signalfx_dashboard_group" "mydashboardgroup0" {
    name = "My team dashboard group"
    description = "Cool dashboard group"
}
```

## Example Usage With Mirrored Dashboards

```terraform
resource "signalfx_dashboard_group" "mydashboardgroup_withmirrors" {
    name = "My team dashboard group"
    description = "Cool dashboard group"

    // You can add as many of these as you like. Make sure your account
    // supports this feature!
    dashboard {
      dashboard_id = "${signalfx_dashboard.gc_dashboard.id}"
      name_override = "GC For My Service"
      description_override = "Garbage Collection dashboard maintained by JVM team"

      filter_override {
        property = "service"
        values = [ "myservice" ]
        negated = false
      }

      variable_override {
        property = "region"
        values = ["us-west1"]
        values_suggested = ["us-west-1", "us-east-1"]
      }
    }
}
```

## Argument Reference

The following arguments are supported in the resource block:

* `name` - (Required) Name of the dashboard group.
* `description` - (Required) Description of the dashboard group.
* `teams` - (Optional) Team IDs to associate the dashboard group to.
* `dashboard` - (Optional) [Mirrored dashboards](https://docs.signalfx.com/en/latest/dashboards/dashboard-mirrors.html) in this dashboard group. **Note:** This feature is not present in all accounts. Please contact support if you are unsure.
  * `dashboard_id` - (Required) The dashboard id to mirror
  * `name_override` - (Optional) The name that will override the original dashboards's name.
  * `description_override` - (Optional) The description that will override the original dashboards's description.
  * `filter_override` - (Optional) The description that will override the original dashboards's description.
    * `property` - (Required) The name of a dimension to filter against.
    * `values` - (Required) A list of values to be used with the `property`, they will be combined via `OR`.
    * `negated` - (Optional) If true,  only data that does not match the specified value of the specified property appear in the event overlay. Defaults to `false`.
  * `filter_override` - (Optional) The description that will override the original dashboards's description.
    * `property` - (Required) A metric time series dimension or property name.
    * `values` - (Optional) (Optional) List of of strings (which will be treated as an OR filter on the property).
    * `values_suggested` - (Optional) A list of strings of suggested values for this variable; these suggestions will receive priority when values are autosuggested for this variable.
