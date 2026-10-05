package datasource

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/kitokinha/terraform-provider/internal/client"
	"github.com/kitokinha/terraform-provider/internal/serviceaccounts"
)

// ServiceAccounts returns all service accounts the configured token can see.
func ServiceAccounts() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceServiceAccountsRead,
		Schema: map[string]*schema.Schema{
			"service_accounts": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{Schema: map[string]*schema.Schema{
					"id":         {Type: schema.TypeString, Computed: true},
					"name":       {Type: schema.TypeString, Computed: true},
					"role_id":    {Type: schema.TypeString, Computed: true},
					"role_name":  {Type: schema.TypeString, Computed: true},
					"created_at": {Type: schema.TypeString, Computed: true},
					"updated_at": {Type: schema.TypeString, Computed: true},
				}},
			},
		},
	}
}

func dataSourceServiceAccountsRead(_ context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	accounts, err := serviceaccounts.List(meta.(*client.PhaseClient))
	if err != nil {
		return diag.FromErr(err)
	}

	values := make([]any, 0, len(accounts))
	for _, account := range accounts {
		values = append(values, map[string]any{
			"id": account.ID, "name": account.Name,
			"role_id": account.Role.ID, "role_name": account.Role.Name,
			"created_at": account.CreatedAt, "updated_at": account.UpdatedAt,
		})
	}
	if err := d.Set("service_accounts", values); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(fmt.Sprintf("service-accounts-%d", len(accounts)))
	return nil
}
