package resources

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestServiceAccountSchema(t *testing.T) {
	resource := ServiceAccount()
	for _, field := range []string{"name", "role_id", "token_name", "team_id", "initial_token", "initial_bearer_token", "initial_token_id"} {
		if _, ok := resource.Schema[field]; !ok {
			t.Errorf("expected schema field %q", field)
		}
	}
	if _, exists := resource.Schema["access"]; exists {
		t.Fatal("access must be managed by ServiceAccountAccess")
	}
	if !resource.Schema["initial_token"].Sensitive || !resource.Schema["initial_bearer_token"].Sensitive {
		t.Fatal("initial token values must be sensitive")
	}
}

func TestResourceServiceAccountCreateAndRead(t *testing.T) {
	requests := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		switch requests {
		case 1:
			if r.Method != http.MethodPost || r.URL.Path != "/v1/service-accounts/" {
				t.Fatalf("unexpected create request: %s %s", r.Method, r.URL.Path)
			}
			var body map[string]any
			_ = json.NewDecoder(r.Body).Decode(&body)
			if body["name"] != "deploy-bot" || body["role_id"] != "role-1" || body["token_name"] != "CI" {
				t.Fatalf("unexpected create body: %#v", body)
			}
			_, _ = w.Write([]byte("{\"id\":\"sa-1\",\"initialToken\":{\"id\":\"token-1\",\"token\":\"secret-token\",\"bearerToken\":\"ServiceAccount secret-token\"}}"))
		case 2:
			if r.Method != http.MethodGet || r.URL.Path != "/v1/service-accounts/sa-1/" {
				t.Fatalf("unexpected read request: %s %s", r.Method, r.URL.Path)
			}
			_, _ = w.Write([]byte("{\"id\":\"sa-1\",\"name\":\"deploy-bot\",\"role\":{\"id\":\"role-1\"}}"))
		default:
			t.Fatalf("unexpected request %d", requests)
		}
	}))
	defer server.Close()

	data := schema.TestResourceDataRaw(t, ServiceAccount().Schema, map[string]any{"name": "deploy-bot", "role_id": "role-1", "token_name": "CI"})
	diags := resourceServiceAccountCreate(context.Background(), data, newApplicationTestClient(server))
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}
	if data.Id() != "sa-1" || data.Get("initial_token").(string) != "secret-token" {
		t.Fatalf("unexpected service account state: ID=%q token=%q", data.Id(), data.Get("initial_token"))
	}
}

func TestResourceServiceAccountReadNotFoundRemovesState(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusNotFound) }))
	defer server.Close()
	data := schema.TestResourceDataRaw(t, ServiceAccount().Schema, map[string]any{"name": "deploy-bot", "role_id": "role-1"})
	data.SetId("sa-1")
	diags := resourceServiceAccountRead(context.Background(), data, newApplicationTestClient(server))
	if diags.HasError() || data.Id() != "" {
		t.Fatalf("expected state removal, diagnostics=%v ID=%q", diags, data.Id())
	}
}
