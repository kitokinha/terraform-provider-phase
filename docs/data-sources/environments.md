# phase_environments (Data Source)

Lists every environment the configured token can access in one Phase application.

## Example Usage

~~~hcl
data "phase_environments" "api" {
  app_id = phase_application.api.id
}

output "production_environment_ids" {
  value = [
    for environment in data.phase_environments.api.environments :
    environment.id if environment.env_type == "prod"
  ]
}
~~~

## Argument Reference

| Argument | Required | Description |
| --- | --- | --- |
| app_id | Yes | Phase application ID. |

## Attribute Reference

- environments — List of environments visible to the configured token. Each item exports id, name, env_type, index, created_at, and updated_at.

The possible env_type values are dev, staging, prod, and custom.
