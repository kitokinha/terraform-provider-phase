# phase_service_account_access (Resource)

Manages the environments granted to one application for one Phase service account.

## Example Usage

Use one resource for each application. These resources can share the same service account.

~~~hcl
resource "phase_service_account_access" "backend" {
  service_account_id = phase_service_account.deploy.id
  app_id             = phase_application.backend.id
  environment_ids    = [var.development_environment_id]
}

resource "phase_service_account_access" "frontend" {
  service_account_id = phase_service_account.deploy.id
  app_id             = phase_application.frontend.id
  environment_ids = [
    var.development_environment_id,
    var.production_environment_id,
  ]
}
~~~

## Lifecycle behavior

Phase accepts the complete access list in one request, rather than a single-app mutation. To make one resource behave as a single application grant, the provider does the following while holding a lock for that service account:

1. Reads the current access list.
2. Adds, replaces, or removes only the configured application.
3. Sends the resulting complete list back to Phase.

This means that deleting the frontend resource above removes only frontend access; it preserves backend access.

> Do not declare two phase_service_account_access resources with the same service_account_id and app_id. They would represent the same remote grant.

## Argument Reference

| Argument | Required | Description |
| --- | --- | --- |
| service_account_id | Yes | ID of the service account. Changing it creates a new resource. |
| app_id | Yes | Application ID. Changing it creates a new resource. |
| environment_ids | Yes | Non-empty set of environment IDs this service account can access in the application. |

## Import

~~~sh
terraform import phase_service_account_access.backend <service-account-id>:<app-id>
~~~

## Migration from the earlier aggregate resource

The resource now represents one application grant, rather than the complete access list. Replace each access block with its own phase_service_account_access resource.
