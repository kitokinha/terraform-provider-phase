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
	for _, field := range []string{"name", "role_id", "token_name", "team_id", "access", "initial_token", "initial_bearer_token", "initial_token_id"} {
		if _, ok := resource.Schema[field]; !ok {
			t.Errorf("expected schema field %q", field)
		}
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
			_, _ = w.Write([]byte(`{"id":"sa-1","name":"deploy-bot","role":{"id":"role-1"},"initialToken":{"id":"token-1","token":"secret-token","bearerToken":"ServiceAccount secret-token"}}`))
		case 2:
			if r.Method != http.MethodPut || r.URL.Path != "/v1/service-accounts/sa-1/access/" {
				t.Fatalf("unexpected access request: %s %s", r.Method, r.URL.Path)
			}
			var body map[string]any
			_ = json.NewDecoder(r.Body).Decode(&body)
			if apps := body["apps"].([]any); len(apps) != 1 {
				t.Fatalf("expected one app grant, got %#v", body)
			}
			_, _ = w.Write([]byte(`{"id":"sa-1"}`))
		case 3:
			if r.Method != http.MethodGet || r.URL.Path != "/v1/service-accounts/sa-1/" {
				t.Fatalf("unexpected read request: %s %s", r.Method, r.URL.Path)
			}
			_, _ = w.Write([]byte(`{"id":"sa-1","name":"deploy-bot","role":{"id":"role-1","name":"Service"},"createdAt":"2026-01-01T00:00:00Z","updatedAt":"2026-01-01T00:00:00Z","apps":[{"id":"app-1","environments":[{"id":"dev"},{"id":"prod"}]}]}`))
		default:
			t.Fatalf("unexpected request %d", requests)
		}
	}))
	defer server.Close()
	resource := ServiceAccount()
	data := schema.TestResourceDataRaw(t, resource.Schema, map[string]any{"name": "deploy-bot", "role_id": "role-1", "token_name": "CI", "access": []any{map[string]any{"app_id": "app-1", "environment_ids": []any{"dev", "prod"}}}})
	diags := resourceServiceAccountCreate(context.Background(), data, newApplicationTestClient(server))
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}
	if data.Id() != "sa-1" {
		t.Fatalf("expected ID sa-1, got %q", data.Id())
	}
	if data.Get("initial_token").(string) != "secret-token" {
		t.Fatal("expected initial token in state")
	}
	if data.Get("access").(*schema.Set).Len() != 1 {
		t.Fatal("expected access in state")
	}
}

func TestResourceServiceAccountUpdateClearsAccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPut && r.URL.Path == "/v1/service-accounts/sa-1/" {
			_, _ = w.Write([]byte(`{"id":"sa-1"}`))
			return
		}
		if r.Method == http.MethodPut && r.URL.Path == "/v1/service-accounts/sa-1/access/" {
			var body map[string]any
			_ = json.NewDecoder(r.Body).Decode(&body)
			if apps, ok := body["apps"].([]any); !ok || len(apps) != 0 {
				t.Fatalf("expected empty access list, got %#v", body)
			}
			_, _ = w.Write([]byte(`{"id":"sa-1"}`))
			return
		}
		if r.Method == http.MethodGet && r.URL.Path == "/v1/service-accounts/sa-1/" {
			_, _ = w.Write([]byte(`{"id":"sa-1","name":"deploy-bot","role":{"id":"role-1"},"apps":[]}`))
			return
		}
		t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
	}))
	defer server.Close()
	resource := ServiceAccount()
	data := schema.TestResourceDataRaw(t, resource.Schema, map[string]any{"name": "deploy-bot", "role_id": "role-1", "access": []any{map[string]any{"app_id": "app-1", "environment_ids": []any{"dev"}}}})
	data.SetId("sa-1")
	if err := data.Set("access", []any{}); err != nil {
		t.Fatal(err)
	}
	diags := resourceServiceAccountUpdate(context.Background(), data, newApplicationTestClient(server))
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags)
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
