# Phase Terraform Provider

Manage Phase applications, secrets, and service accounts with Terraform.

## Contents

- [Getting started](#getting-started)
- [Provider configuration](#provider-configuration)
- [Resources](#resources)
- [Data sources](#data-sources)
- [Import](#import)

## Getting started

Configure the provider with a Phase personal access token or service token. Keep the token outside your Terraform files whenever possible.

```hcl
terraform {
  required_providers {
    phase = {
      source  = "kitokinha/phase"
      version = "1.0.0"
    }
  }
}

provider "phase" {}
```

```sh
export PHASE_TOKEN="pss_user:v1:..."
```

The following example creates an application, stores a secret, and creates a service account with access to one environment.

```hcl
resource "phase_application" "api" {
  name        = "api"
  description = "API configuration"
}

resource "phase_secret" "database_url" {
  app_id = phase_application.api.id
  env    = "production"
  key    = "DATABASE_URL"
  value  = var.database_url
  path   = "/database/"
  tags   = ["database"]
}

resource "phase_service_account" "deploy" {
  name       = "deploy-bot"
  role_id    = var.deployer_role_id
  token_name = "CI"

  access {
    app_id          = phase_application.api.id
    environment_ids = [var.production_environment_id]
  }
}

output "deploy_token" {
  value     = phase_service_account.deploy.initial_token
  sensitive = true
}
```

## Provider configuration

| Argument | Required | Description |
| --- | --- | --- |
| `phase_token` | No | Phase personal access token (`pss_user:...`) or service token (`pss_service:...`). It is sensitive. |
| `host` | No | Phase API host. Defaults to `https://api.phase.dev`. Custom hosts receive the `/service/public` path automatically. |
| `skip_tls_verification` | No | Disables TLS certificate validation for a custom host. Defaults to `false`; use only for trusted self-hosted installations. |

`phase_token` can be supplied directly or through the first available environment variable in this order: `PHASE_TOKEN`, `PHASE_SERVICE_TOKEN`, then `PHASE_PAT_TOKEN`.

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

## Resources

### `phase_application`

Manages a Phase application.

```hcl
resource "phase_application" "api" {
  name        = "api"
  description = "API configuration"
}
```

| Argument | Required | Description |
| --- | --- | --- |
| `name` | Yes | Application name. |
| `description` | No | Application description. |

Exports `id`, `created_at`, and `updated_at`.

### `phase_secret`

Manages one secret in an application environment. If a secret with the same application, environment, and key already exists, the provider updates it during creation rather than failing.

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

| Argument | Required | Description |
| --- | --- | --- |
| `app_id` | Yes, replacement | Application ID. Changing it creates a new resource. |
| `env` | Yes, replacement | Environment name, such as `development` or `production`. Changing it creates a new resource. |
| `key` | Yes | Secret key. |
| `value` | Yes | Secret value. Stored as sensitive state. |
| `path` | No | Secret path. Defaults to `/`. |
| `comment` | No | Secret description. |
| `tags` | No | Tags to assign to the secret. Tags must already exist in Phase. |
| `override` | No | Personal override block; requires a personal access token. |

`override` supports `value` (sensitive) and `is_active` (boolean). The override belongs to the authenticated user, not to the application.

Exports `id`, `version`, `created_at`, and `updated_at`.

### `phase_service_account`

Creates and manages a Phase service account. Its `access` blocks are authoritative: Terraform sends the complete desired app/environment access set to Phase. Do not manage access for the same service account from another configuration or resource.

```hcl
resource "phase_service_account" "deploy" {
  name       = "deploy-bot"
  role_id    = var.deployer_role_id
  token_name = "CI"

  access {
    app_id          = phase_application.api.id
    environment_ids = [var.development_environment_id, var.production_environment_id]
  }
}
```

| Argument | Required | Description |
| --- | --- | --- |
| `name` | Yes | Service account name. |
| `role_id` | Yes | ID of a non-global Phase role. |
| `token_name` | No, replacement | Name of the initial token. Phase creates this token only when the account is created. |
| `team_id` | No, replacement | Owning team ID. An account cannot move to a different team. |
| `access` | No | One or more app access blocks. Each block requires `app_id` and a non-empty `environment_ids` set. |

Exports `id`, `created_at`, `updated_at`, `initial_token`, `initial_bearer_token`, and `initial_token_id`.

> `initial_token` and `initial_bearer_token` are sensitive and returned by Phase only once, at creation. Save them in a secure destination. They cannot be recovered by refresh or import.

## Data sources

### `phase_service_accounts`

Lists all service accounts visible to the configured token. This is a summary listing; it does not expose tokens or environment grants.

```hcl
data "phase_service_accounts" "all" {}

output "service_account_names" {
  value = [
    for account in data.phase_service_accounts.all.service_accounts : account.name
  ]
}
```

Each item in `service_accounts` exports `id`, `name`, `role_id`, `role_name`, `created_at`, and `updated_at`.

### `phase_secrets`

Reads secrets from one application environment. The returned map is sensitive.

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

| Argument | Required | Description |
| --- | --- | --- |
| `app_id` | Yes | Application ID. |
| `env` | Yes | Environment name. |
| `path` | No | Path filter. |
| `key` | No | Restrict results to one secret key. |
| `tags` | No | Return secrets with any supplied tag. |

Exports `secrets`, a sensitive map from secret key to value. An active personal override for the authenticated user is returned in preference to the regular secret value.

## Import

Import an existing application by ID:

```sh
terraform import phase_application.api <application-id>
```

Import an existing secret with its application ID, environment, path, and key:

```sh
terraform import phase_secret.database_url "<app-id>:production:/database/:DATABASE_URL"
```

Import a service account by ID:

```sh
terraform import phase_service_account.deploy <service-account-id>
```

After import, add all configurable arguments to your `.tf` configuration and run `terraform plan`. Initial service-account tokens cannot be imported because Phase does not return them after creation.

## Development

```sh
go test -count=1 ./...
go build ./...
```
