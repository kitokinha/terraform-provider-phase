package resources

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/kitokinha/terraform-provider/internal/client"
	"github.com/kitokinha/terraform-provider/internal/serviceaccounts"
)

var serviceAccountAccessLocks sync.Map

// ServiceAccountAccess manages the environments granted to one app for one service account.
func ServiceAccountAccess() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceServiceAccountAccessCreate,
		ReadContext:   resourceServiceAccountAccessRead,
		UpdateContext: resourceServiceAccountAccessUpdate,
		DeleteContext: resourceServiceAccountAccessDelete,
		Importer:      &schema.ResourceImporter{StateContext: serviceAccountAccessImportState},
		Schema: map[string]*schema.Schema{
			"service_account_id": {Type: schema.TypeString, Required: true, ForceNew: true},
			"app_id":             {Type: schema.TypeString, Required: true, ForceNew: true},
			"environment_ids": {
				Type: schema.TypeSet, Required: true, MinItems: 1,
				Elem: &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceServiceAccountAccessCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	accountID := d.Get("service_account_id").(string)
	if err := mutateServiceAccountAccess(meta.(*client.PhaseClient), accountID, d.Get("app_id").(string), environmentIDs(d), false); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(serviceAccountAccessID(accountID, d.Get("app_id").(string)))
	return resourceServiceAccountAccessRead(ctx, d, meta)
}

func resourceServiceAccountAccessRead(_ context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	accountID := d.Get("service_account_id").(string)
	appID := d.Get("app_id").(string)
	account, err := serviceaccounts.Read(meta.(*client.PhaseClient), accountID)
	if err != nil {
		if client.IsNotFound(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}
	for _, app := range account.Apps {
		if app.ID != appID {
			continue
		}
		if err := d.Set("environment_ids", environmentIDsFromApp(app)); err != nil {
			return diag.FromErr(err)
		}
		return nil
	}
	d.SetId("")
	return nil
}

func resourceServiceAccountAccessUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	accountID := d.Get("service_account_id").(string)
	if err := mutateServiceAccountAccess(meta.(*client.PhaseClient), accountID, d.Get("app_id").(string), environmentIDs(d), false); err != nil {
		return diag.FromErr(err)
	}
	return resourceServiceAccountAccessRead(ctx, d, meta)
}

func resourceServiceAccountAccessDelete(_ context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	accountID := d.Get("service_account_id").(string)
	err := mutateServiceAccountAccess(meta.(*client.PhaseClient), accountID, d.Get("app_id").(string), nil, true)
	if err != nil && !client.IsNotFound(err) {
		return diag.FromErr(err)
	}
	d.SetId("")
	return nil
}

// mutateServiceAccountAccess reads the complete grant list, changes only appID,
// then replaces the complete list as required by the Phase API.
func mutateServiceAccountAccess(c *client.PhaseClient, accountID, appID string, desiredEnvironments []string, remove bool) error {
	return withServiceAccountAccessLock(accountID, func() error {
		account, err := serviceaccounts.Read(c, accountID)
		if err != nil {
			return err
		}

		apps := make([]serviceaccounts.AccessApp, 0, len(account.Apps)+1)
		found := false
		for _, app := range account.Apps {
			if app.ID != appID {
				apps = append(apps, serviceaccounts.AccessApp{ID: app.ID, Environments: environmentIDsFromApp(app)})
				continue
			}
			found = true
			if !remove {
				apps = append(apps, serviceaccounts.AccessApp{ID: appID, Environments: desiredEnvironments})
			}
		}
		if !remove && !found {
			apps = append(apps, serviceaccounts.AccessApp{ID: appID, Environments: desiredEnvironments})
		}
		_, err = serviceaccounts.SetAccess(c, accountID, serviceaccounts.AccessRequest{Apps: apps})
		return err
	})
}

func environmentIDs(d *schema.ResourceData) []string {
	values := d.Get("environment_ids").(*schema.Set).List()
	ids := make([]string, 0, len(values))
	for _, value := range values {
		ids = append(ids, value.(string))
	}
	return ids
}

func environmentIDsFromApp(app serviceaccounts.App) []string {
	ids := make([]string, 0, len(app.Environments))
	for _, environment := range app.Environments {
		ids = append(ids, environment.ID)
	}
	return ids
}

func serviceAccountAccessID(accountID, appID string) string {
	return accountID + ":" + appID
}

func withServiceAccountAccessLock(accountID string, fn func() error) error {
	value, _ := serviceAccountAccessLocks.LoadOrStore(accountID, &sync.Mutex{})
	lock := value.(*sync.Mutex)
	lock.Lock()
	defer lock.Unlock()
	return fn()
}

func serviceAccountAccessImportState(_ context.Context, d *schema.ResourceData, _ any) ([]*schema.ResourceData, error) {
	accountID, appID, found := strings.Cut(d.Id(), ":")
	if !found || accountID == "" || appID == "" {
		return nil, fmt.Errorf("import ID must use the format <service-account-id>:<app-id>")
	}
	if err := d.Set("service_account_id", accountID); err != nil {
		return nil, err
	}
	if err := d.Set("app_id", appID); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}
