---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "flagsmith_feature Resource - terraform-provider-flagsmith"
subcategory: ""
description: |-
  Flagsmith Feature/ Remote config
---

# flagsmith_feature (Resource)

Flagsmith Feature/ Remote config

## Example Usage

```terraform
resource "flagsmith_feature" "new_mv_feature" {
  feature_name = "new_mv_feature"
  project_uuid = "10421b1f-5f29-4da9-abe2-30f88c07c9e8"
  description  = "This is a new multivariate feature"
  type         = "MULTIVARIATE"
}

resource "flagsmith_feature" "new_standard_feature" {
  feature_name = "new_standard_feature"
  project_uuid = "10421b1f-5f29-4da9-abe2-30f88c07c9e8"
  description  = "This is a new standard feature"
  type         = "STANDARD"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `feature_name` (String) Name of the feature
- `project_uuid` (String) UUID of project the feature belongs to
- `type` (String) Type of the feature, can be STANDARD, or MULTIVARIATE

### Optional

- `default_enabled` (Boolean) Determines if the feature is enabled by default. If unspecified, it will default to false
- `description` (String) Description of the feature
- `initial_value` (String) Determines the initial value of the feature.
- `is_archived` (Boolean) Can be used to archive/unarchive a feature. If unspecified, it will default to false
- `owners` (Set of Number) List of user IDs who are owners of the feature

### Read-Only

- `id` (Number) ID of the feature
- `project_id` (Number) ID of the project
- `uuid` (String) UUID of the feature

## Import

Import is supported using the following syntax:

```shell
terraform import flagsmith_feature.some_feature <feature_uuid>
```