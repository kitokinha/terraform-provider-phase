package datasource

import (
	"context"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestDataSourceEnvironmentsRead(t *testing.T) {
	phaseClient := newTestSecretsDataSourceClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/v1/environments/" {
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		if appID := r.URL.Query().Get("app_id"); appID != "app-1" {
			t.Fatalf("expected app_id app-1, got %q", appID)
		}
		_, _ = w.Write([]byte("{\"data\":[{\"id\":\"env-1\",\"name\":\"Development\",\"envType\":\"dev\",\"index\":0,\"createdAt\":\"2026-01-01T00:00:00Z\",\"updatedAt\":\"2026-01-02T00:00:00Z\"}]}"))
	}))

	d := schema.TestResourceDataRaw(t, Environments().Schema, map[string]any{"app_id": "app-1"})
	diags := dataSourceEnvironmentsRead(context.Background(), d, phaseClient)
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}
	items := d.Get("environments").([]any)
	if len(items) != 1 {
		t.Fatalf("expected one environment, got %#v", items)
	}
	environment := items[0].(map[string]any)
	if environment["id"] != "env-1" || environment["env_type"] != "dev" {
		t.Fatalf("unexpected environment: %#v", environment)
	}
}
