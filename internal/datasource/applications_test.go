package datasource

import (
	"context"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestDataSourceApplicationsRead(t *testing.T) {
	phaseClient := newTestSecretsDataSourceClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/v1/apps/" {
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		_, _ = w.Write([]byte("{\"data\":[{\"id\":\"app-1\",\"name\":\"Backend\",\"description\":\"Backend API\",\"sseEnabled\":true,\"createdAt\":\"2026-01-01T00:00:00Z\",\"updatedAt\":\"2026-01-02T00:00:00Z\"}]}"))
	}))

	d := schema.TestResourceDataRaw(t, Applications().Schema, map[string]any{})
	diags := dataSourceApplicationsRead(context.Background(), d, phaseClient)
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}
	items := d.Get("applications").([]any)
	if len(items) != 1 {
		t.Fatalf("expected one application, got %#v", items)
	}
	application := items[0].(map[string]any)
	if application["id"] != "app-1" || application["sse_enabled"] != true {
		t.Fatalf("unexpected application: %#v", application)
	}
}
