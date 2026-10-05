package datasource

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/kitokinha/terraform-provider/internal/client"
	"github.com/kitokinha/terraform-provider/internal/environments"
)

// Environments returns all environments visible to the configured token for an app.
func Environments() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceEnvironmentsRead,
		Schema: map[string]*schema.Schema{
			"app_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the Phase application.",
			},
			"environments": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{Schema: map[string]*schema.Schema{
					"id":         {Type: schema.TypeString, Computed: true},
					"name":       {Type: schema.TypeString, Computed: true},
					"env_type":   {Type: schema.TypeString, Computed: true},
					"index":      {Type: schema.TypeInt, Computed: true},
					"created_at": {Type: schema.TypeString, Computed: true},
					"updated_at": {Type: schema.TypeString, Computed: true},
				}},
			},
		},
	}
}

func dataSourceEnvironmentsRead(_ context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	appID := d.Get("app_id").(string)
	items, err := environments.List(meta.(*client.PhaseClient), appID)
	if err != nil {
		return diag.FromErr(err)
	}

	values := make([]any, 0, len(items))
	for _, environment := range items {
		values = append(values, map[string]any{
			"id": environment.ID, "name": environment.Name,
			"env_type": environment.EnvType, "index": environment.Index,
			"created_at": environment.CreatedAt, "updated_at": environment.UpdatedAt,
		})
	}
	if err := d.Set("environments", values); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(appID)
	return nil
}
