package datasource

import (
	"context"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestDataSourceServiceAccountsRead(t *testing.T) {
	client := newTestSecretsDataSourceClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/v1/service-accounts/" {
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		_, _ = w.Write([]byte(`{"data":[{"id":"sa-1","name":"deploy-bot","role":{"id":"role-1","name":"Service"},"createdAt":"2026-01-01T00:00:00Z","updatedAt":"2026-01-02T00:00:00Z"}]}`))
	}))
	d := schema.TestResourceDataRaw(t, ServiceAccounts().Schema, map[string]any{})
	diags := dataSourceServiceAccountsRead(context.Background(), d, client)
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}
	accounts := d.Get("service_accounts").([]any)
	if len(accounts) != 1 {
		t.Fatalf("expected one service account, got %#v", accounts)
	}
	account := accounts[0].(map[string]any)
	if account["id"] != "sa-1" || account["role_id"] != "role-1" {
		t.Fatalf("unexpected account: %#v", account)
	}
}
