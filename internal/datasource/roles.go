package datasource

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/kitokinha/terraform-provider/internal/client"
	"github.com/kitokinha/terraform-provider/internal/roles"
)

// Roles returns all default and custom roles visible to the configured token.
func Roles() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRolesRead,
		Schema: map[string]*schema.Schema{
			"roles": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{Schema: map[string]*schema.Schema{
					"id":          {Type: schema.TypeString, Computed: true},
					"name":        {Type: schema.TypeString, Computed: true},
					"description": {Type: schema.TypeString, Computed: true},
					"color":       {Type: schema.TypeString, Computed: true},
					"is_default":  {Type: schema.TypeBool, Computed: true},
					"created_at":  {Type: schema.TypeString, Computed: true},
				}},
			},
		},
	}
}

func dataSourceRolesRead(_ context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	items, err := roles.List(meta.(*client.PhaseClient))
	if err != nil {
		return diag.FromErr(err)
	}

	values := make([]any, 0, len(items))
	for _, role := range items {
		description := ""
		if role.Description != nil {
			description = *role.Description
		}
		values = append(values, map[string]any{
			"id": role.ID, "name": role.Name, "description": description,
			"color": role.Color, "is_default": role.IsDefault, "created_at": role.CreatedAt,
		})
	}
	if err := d.Set("roles", values); err != nil {
		return diag.FromErr(err)
	}
	d.SetId("roles")
	return nil
}
