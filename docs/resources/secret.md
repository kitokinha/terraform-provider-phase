# phase_secret (Resource)

Manages one secret in a Phase application environment. If a secret with the same application, environment, and key already exists, the provider updates it during creation instead of failing.

## Example Usage

```hcl
resource "phase_secret" "database_url" {
  app_id  = phase_application.api.id
  env     = "production"
  key     = "DATABASE_URL"
  value   = var.database_url
  path    = "/database/"
  comment = "Managed by Terraform"
  tags    = ["database", "credentials"]
}
```

## Argument Reference

| Argument | Required | Description |
| --- | --- | --- |
| `app_id` | Yes | Application ID. Changing it creates a new resource. |
| `env` | Yes | Environment name, such as `development` or `production`. Changing it creates a new resource. |
| `key` | Yes | Secret key. |
| `value` | Yes | Secret value. Stored as sensitive state. |
| `path` | No | Secret path. Defaults to `/`. |
| `comment` | No | Secret description. |
| `tags` | No | Tags to assign to the secret. Tags must already exist in Phase. |
| `override` | No | Personal override block; requires a personal access token. |

The optional `override` block supports:

- `value` — Override value, stored as sensitive state.
- `is_active` — Whether the override is active.

An override belongs to the authenticated user, not to the application.

## Attribute Reference

In addition to the arguments above, this resource exports `id`, `version`, `created_at`, and `updated_at`.

## Import

```sh
terraform import phase_secret.database_url "<app-id>:production:/database/:DATABASE_URL"
```
