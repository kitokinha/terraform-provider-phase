package datasource

import (
	"context"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestDataSourceRolesRead(t *testing.T) {
	phaseClient := newTestSecretsDataSourceClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/v1/roles/" {
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		_, _ = w.Write([]byte("{\"data\":[{\"id\":\"role-1\",\"name\":\"Service\",\"description\":null,\"color\":\"\",\"isDefault\":true,\"createdAt\":\"2026-01-01T00:00:00Z\"}]}"))
	}))

	d := schema.TestResourceDataRaw(t, Roles().Schema, map[string]any{})
	diags := dataSourceRolesRead(context.Background(), d, phaseClient)
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}
	items := d.Get("roles").([]any)
	if len(items) != 1 {
		t.Fatalf("expected one role, got %#v", items)
	}
	role := items[0].(map[string]any)
	if role["id"] != "role-1" || role["is_default"] != true || role["description"] != "" {
		t.Fatalf("unexpected role: %#v", role)
	}
}
