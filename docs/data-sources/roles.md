# phase_roles (Data Source)

Lists all default and custom roles visible to the configured token.

## Example Usage

~~~hcl
data "phase_roles" "all" {}

output "service_role_ids" {
  value = [
    for role in data.phase_roles.all.roles :
    role.id if role.name == "Service"
  ]
}
~~~

## Attribute Reference

- roles — List of visible roles. Each item exports id, name, description, color, is_default, and created_at.

Default roles have is_default set to true. The list endpoint returns role summaries only; it does not include permission details.
