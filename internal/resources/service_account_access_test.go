package resources

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestServiceAccountAccessSchema(t *testing.T) {
	resource := ServiceAccountAccess()
	for _, field := range []string{"service_account_id", "app_id", "environment_ids"} {
		if _, ok := resource.Schema[field]; !ok {
			t.Errorf("expected schema field %q", field)
		}
	}
	if !resource.Schema["service_account_id"].ForceNew || !resource.Schema["app_id"].ForceNew {
		t.Fatal("service account and app IDs must be ForceNew")
	}
}

func TestResourceServiceAccountAccessMutatesOneApp(t *testing.T) {
	requests := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		switch requests {
		case 1: // Create reads current grants.
			writeServiceAccountAccess(t, w, nil, "env-2")
		case 2: // Create writes app-1 and preserves app-2.
			assertAccessRequest(t, r, []string{"app-1", "app-2"}, "app-1", []string{"env-1"})
			_, _ = w.Write([]byte("{\"id\":\"sa-1\"}"))
		case 3: // Create reads the resulting grant.
			writeServiceAccountAccess(t, w, []string{"env-1"}, "env-2")
		case 4: // Update reads current grants.
			writeServiceAccountAccess(t, w, []string{"env-1"}, "env-2")
		case 5: // Update changes only app-1.
			assertAccessRequest(t, r, []string{"app-1", "app-2"}, "app-1", []string{"env-1", "env-3"})
			_, _ = w.Write([]byte("{\"id\":\"sa-1\"}"))
		case 6: // Update reads the resulting grant.
			writeServiceAccountAccess(t, w, []string{"env-1", "env-3"}, "env-2")
		case 7: // Delete reads current grants.
			writeServiceAccountAccess(t, w, []string{"env-1", "env-3"}, "env-2")
		case 8: // Delete removes app-1 and preserves app-2.
			assertAccessRequest(t, r, []string{"app-2"}, "", nil)
			_, _ = w.Write([]byte("{\"id\":\"sa-1\"}"))
		default:
			t.Fatalf("unexpected request %d: %s %s", requests, r.Method, r.URL.Path)
		}
	}))
	defer server.Close()

	resource := ServiceAccountAccess()
	data := schema.TestResourceDataRaw(t, resource.Schema, map[string]any{
		"service_account_id": "sa-1",
		"app_id":             "app-1",
		"environment_ids":    []any{"env-1"},
	})
	client := newApplicationTestClient(server)
	if diags := resourceServiceAccountAccessCreate(context.Background(), data, client); diags.HasError() {
		t.Fatalf("unexpected create diagnostics: %v", diags)
	}
	if err := data.Set("environment_ids", []any{"env-1", "env-3"}); err != nil {
		t.Fatal(err)
	}
	if diags := resourceServiceAccountAccessUpdate(context.Background(), data, client); diags.HasError() {
		t.Fatalf("unexpected update diagnostics: %v", diags)
	}
	if diags := resourceServiceAccountAccessDelete(context.Background(), data, client); diags.HasError() {
		t.Fatalf("unexpected delete diagnostics: %v", diags)
	}
	if data.Id() != "" {
		t.Fatalf("expected deleted state, got %q", data.Id())
	}
}

func writeServiceAccountAccess(t *testing.T, w http.ResponseWriter, firstEnvironments []string, secondEnvironment string) {
	t.Helper()
	apps := []map[string]any{
		{"id": "app-2", "environments": []map[string]string{{"id": secondEnvironment}}},
	}
	if len(firstEnvironments) > 0 {
		environments := make([]map[string]string, 0, len(firstEnvironments))
		for _, environment := range firstEnvironments {
			environments = append(environments, map[string]string{"id": environment})
		}
		apps = append([]map[string]any{{"id": "app-1", "environments": environments}}, apps...)
	}
	if err := json.NewEncoder(w).Encode(map[string]any{"id": "sa-1", "apps": apps}); err != nil {
		t.Fatal(err)
	}
}

func assertAccessRequest(t *testing.T, r *http.Request, expectedApps []string, changedApp string, expectedEnvironments []string) {
	t.Helper()
	if r.Method != http.MethodPut || r.URL.Path != "/v1/service-accounts/sa-1/access/" {
		t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
	}
	var body struct {
		Apps []struct {
			ID           string   `json:"id"`
			Environments []string `json:"environments"`
		} `json:"apps"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		t.Fatal(err)
	}
	if len(body.Apps) != len(expectedApps) {
		t.Fatalf("expected apps %v, got %#v", expectedApps, body.Apps)
	}
	found := map[string][]string{}
	for _, app := range body.Apps {
		found[app.ID] = app.Environments
	}
	for _, appID := range expectedApps {
		if _, ok := found[appID]; !ok {
			t.Fatalf("missing app %q in %#v", appID, found)
		}
	}
	if changedApp != "" && !sameStrings(found[changedApp], expectedEnvironments) {
		t.Fatalf("expected environments %v for %s, got %v", expectedEnvironments, changedApp, found[changedApp])
	}
}

func sameStrings(actual, expected []string) bool {
	if len(actual) != len(expected) {
		return false
	}
	found := make(map[string]bool, len(actual))
	for _, value := range actual {
		found[value] = true
	}
	for _, value := range expected {
		if !found[value] {
			return false
		}
	}
	return true
}
