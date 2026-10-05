# phase_service_account (Resource)

Creates and manages a Phase service account and its app/environment access.

## Example Usage

```hcl
resource "phase_service_account" "deploy" {
  name       = "deploy-bot"
  role_id    = var.deployer_role_id
  token_name = "CI"

  access {
    app_id          = phase_application.api.id
    environment_ids = [
      var.development_environment_id,
      var.production_environment_id,
    ]
  }
}

output "deploy_token" {
  value     = phase_service_account.deploy.initial_token
  sensitive = true
}
```

## Argument Reference

| Argument | Required | Description |
| --- | --- | --- |
| `name` | Yes | Service account name. |
| `role_id` | Yes | ID of a non-global Phase role. |
| `token_name` | No | Name of the initial token. Changing it creates a new account. |
| `team_id` | No | Owning team ID. Changing it creates a new account. |
| `access` | No | App/environment access block. |

Each `access` block supports:

- `app_id` — Application ID.
- `environment_ids` — Non-empty set of environment IDs to grant access to.

> Access blocks are authoritative. Terraform sends the complete desired access set to Phase, so do not manage access for the same service account from another configuration or resource.

## Attribute Reference

In addition to the arguments above, this resource exports:

- `id` — Service account ID.
- `created_at` and `updated_at` — Timestamps returned by Phase.
- `initial_token` — Initial service token (sensitive).
- `initial_bearer_token` — Initial bearer token (sensitive).
- `initial_token_id` — Initial token ID.

> Phase returns the initial token values only at creation. Save them in a secure destination: they cannot be recovered by refresh or import.

## Import

```sh
terraform import phase_service_account.deploy <service-account-id>
```

Imported service accounts do not contain the initial token values.
