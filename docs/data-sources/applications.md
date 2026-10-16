# phase_applications (Data Source)

Lists every Phase app visible to the configured token, including both server-side encryption (SSE) and end-to-end encryption (E2EE) apps.

## Example Usage

~~~hcl
data "phase_applications" "all" {}

output "sse_application_ids" {
  value = [
    for application in data.phase_applications.all.applications :
    application.id if application.sse_enabled
  ]
}
~~~

## Attribute Reference

- applications — List of visible applications. Each item exports id, name, description, sse_enabled, created_at, and updated_at.

Only applications with sse_enabled set to true support secret and environment write operations through the Phase REST API.
