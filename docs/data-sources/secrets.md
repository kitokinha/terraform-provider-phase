# phase_secrets (Data Source)

Reads secrets from one Phase application environment. The returned secret map is sensitive.

## Example Usage

```hcl
data "phase_secrets" "database" {
  app_id = phase_application.api.id
  env    = "production"
  path   = "/database/"
  tags   = ["database"]
}

output "database_url" {
  value     = data.phase_secrets.database.secrets["DATABASE_URL"]
  sensitive = true
}
```

## Argument Reference

| Argument | Required | Description |
| --- | --- | --- |
| `app_id` | Yes | Application ID. |
| `env` | Yes | Environment name. |
| `path` | No | Path filter. |
| `key` | No | Restricts results to one secret key. |
| `tags` | No | Returns secrets matching any supplied tag. |

## Attribute Reference

- `secrets` — Sensitive map from secret key to value.

If the authenticated user has an active personal override, the data source returns that value instead of the regular secret value.
