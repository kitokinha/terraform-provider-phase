# phase_service_account (Resource)

Creates and manages a Phase service account. Manage its app/environment access with phase_service_account_access.

## Example Usage

```hcl
resource "phase_service_account" "deploy" {
  name       = "deploy-bot"
  role_id    = var.deployer_role_id
  token_name = "CI"
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
