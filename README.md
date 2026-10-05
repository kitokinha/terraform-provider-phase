# Terraform Provider for Phase

Manage Phase applications, secrets, and service accounts with Terraform.

## Quick start

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

resource "phase_application" "api" {
  name = "api"
}

resource "phase_service_account" "deploy" {
  name    = "deploy-bot"
  role_id = var.deployer_role_id
}

resource "phase_service_account_access" "deploy" {
  service_account_id = phase_service_account.deploy.id
  app_id             = phase_application.api.id
  environment_ids    = [var.production_environment_id]
}
```

Supply authentication through `PHASE_TOKEN` (recommended) or the provider's `phase_token` argument. Both personal access tokens and service tokens are supported.

See the full [provider documentation](docs/index.md) for configuration, all resources, data sources, imports, and service-account token handling.

## Development

```sh
go test -count=1 ./...
go build ./...
```

## License

Distributed under the [MIT License](LICENSE).
