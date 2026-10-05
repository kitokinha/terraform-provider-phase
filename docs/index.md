# Phase Provider

Manage Phase applications, secrets, and service accounts with Terraform.

## Example Usage

```hcl
terraform {
  required_providers {
    phase = {
      source  = "kitokinha/phase"
      version = "1.4.0"
    }
  }
}

provider "phase" {}
```

Set a Phase personal access token or service token outside your Terraform configuration:

```sh
export PHASE_TOKEN="pss_user:v1:..."
```

## Argument Reference

| Argument | Optional | Description |
| --- | --- | --- |
| `phase_token` | Yes | Phase personal access token (`pss_user:...`) or service token (`pss_service:...`). This value is sensitive. |
| `host` | Yes | Phase API host. Defaults to `https://api.phase.dev`. The provider adds `/service/public` to custom hosts. |
| `skip_tls_verification` | Yes | Disables TLS certificate validation for a custom host. Defaults to `false`. |

The provider checks these environment variables for `phase_token`, in order: `PHASE_TOKEN`, `PHASE_SERVICE_TOKEN`, and `PHASE_PAT_TOKEN`.

```hcl
provider "phase" {
  phase_token = var.phase_token
}
```

For a self-hosted installation:

```hcl
provider "phase" {
  host = "https://phase.example.com"
}
```

See the resource and data-source pages for configuration examples and reference details.

## Resources

- [`phase_application`](resources/application.md)
- [`phase_secret`](resources/secret.md)
- [`phase_service_account`](resources/service_account.md)
- [phase_service_account_access](resources/service_account_access.md)

## Data Sources

- [phase_environments](data-sources/environments.md)
- [phase_roles](data-sources/roles.md)
- [`phase_secrets`](data-sources/secrets.md)
- [`phase_service_accounts`](data-sources/service_accounts.md)
