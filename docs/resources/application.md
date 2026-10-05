# phase_application (Resource)

Manages a Phase application.

## Example Usage

```hcl
resource "phase_application" "api" {
  name        = "api"
  description = "API configuration"
}
```

## Argument Reference

| Argument | Required | Description |
| --- | --- | --- |
| `name` | Yes | Application name. |
| `description` | No | Application description. |

## Attribute Reference

In addition to the arguments above, this resource exports:

- `id` — Phase application ID.
- `created_at` — Creation timestamp.
- `updated_at` — Last-update timestamp.

## Import

```sh
terraform import phase_application.api <application-id>
```
