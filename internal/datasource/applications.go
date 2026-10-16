package datasource

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/kitokinha/terraform-provider/internal/applications"
	"github.com/kitokinha/terraform-provider/internal/client"
)

// Applications returns all apps visible to the configured token.
func Applications() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceApplicationsRead,
		Schema: map[string]*schema.Schema{
			"applications": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{Schema: map[string]*schema.Schema{
					"id":          {Type: schema.TypeString, Computed: true},
					"name":        {Type: schema.TypeString, Computed: true},
					"description": {Type: schema.TypeString, Computed: true},
					"sse_enabled": {Type: schema.TypeBool, Computed: true},
					"created_at":  {Type: schema.TypeString, Computed: true},
					"updated_at":  {Type: schema.TypeString, Computed: true},
				}},
			},
		},
	}
}

func dataSourceApplicationsRead(_ context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	items, err := applications.ListApplications(meta.(*client.PhaseClient))
	if err != nil {
		return diag.FromErr(err)
	}

	values := make([]any, 0, len(items))
	for _, application := range items {
		values = append(values, map[string]any{
			"id": application.ID, "name": application.Name,
			"description": application.Description, "sse_enabled": application.SseEnabled,
			"created_at": application.CreatedAt, "updated_at": application.UpdatedAt,
		})
	}
	if err := d.Set("applications", values); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(fmt.Sprintf("applications-%d", len(items)))
	return nil
}
