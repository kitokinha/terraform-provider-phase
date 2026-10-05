# phase_service_accounts (Data Source)

Lists all service accounts visible to the configured token. This is a summary listing and does not expose service-account tokens or app/environment grants.

## Example Usage

```hcl
data "phase_service_accounts" "all" {}

output "service_account_names" {
  value = [
    for account in data.phase_service_accounts.all.service_accounts : account.name
  ]
}
```

## Attribute Reference

- `service_accounts` — List of visible service accounts. Each item exports `id`, `name`, `role_id`, `role_name`, `created_at`, and `updated_at`.
