package resources

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/kitokinha/terraform-provider/internal/client"
	"github.com/kitokinha/terraform-provider/internal/serviceaccounts"
)

func ServiceAccount() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceServiceAccountCreate,
		ReadContext:   resourceServiceAccountRead,
		UpdateContext: resourceServiceAccountUpdate,
		DeleteContext: resourceServiceAccountDelete,
		Importer:      &schema.ResourceImporter{StateContext: serviceAccountImportState},
		Schema: map[string]*schema.Schema{
			"name":                 {Type: schema.TypeString, Required: true},
			"role_id":              {Type: schema.TypeString, Required: true},
			"token_name":           {Type: schema.TypeString, Optional: true, ForceNew: true},
			"team_id":              {Type: schema.TypeString, Optional: true, ForceNew: true},
			"initial_token":        {Type: schema.TypeString, Computed: true, Sensitive: true},
			"initial_bearer_token": {Type: schema.TypeString, Computed: true, Sensitive: true},
			"initial_token_id":     {Type: schema.TypeString, Computed: true},
			"created_at":           {Type: schema.TypeString, Computed: true},
			"updated_at":           {Type: schema.TypeString, Computed: true},
		},
	}
}

func resourceServiceAccountCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	c := meta.(*client.PhaseClient)
	account, err := serviceaccounts.Create(c, serviceaccounts.CreateRequest{
		Name: d.Get("name").(string), RoleID: d.Get("role_id").(string),
		TokenName: d.Get("token_name").(string), TeamID: d.Get("team_id").(string),
	})
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(account.ID)
	if account.InitialToken != nil {
		if err := d.Set("initial_token", account.InitialToken.Token); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("initial_bearer_token", account.InitialToken.BearerToken); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("initial_token_id", account.InitialToken.ID); err != nil {
			return diag.FromErr(err)
		}
	}
	return resourceServiceAccountRead(ctx, d, meta)
}

func resourceServiceAccountRead(_ context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	account, err := serviceaccounts.Read(meta.(*client.PhaseClient), d.Id())
	if err != nil {
		if client.IsNotFound(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}
	if account == nil || account.ID == "" {
		return diag.Errorf("no service account found")
	}
	d.SetId(account.ID)
	if err := d.Set("name", account.Name); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("role_id", account.Role.ID); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("created_at", account.CreatedAt); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("updated_at", account.UpdatedAt); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceServiceAccountUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	c := meta.(*client.PhaseClient)
	if d.HasChange("name") || d.HasChange("role_id") {
		request := serviceaccounts.UpdateRequest{}
		if d.HasChange("name") {
			name := d.Get("name").(string)
			request.Name = &name
		}
		if d.HasChange("role_id") {
			roleID := d.Get("role_id").(string)
			request.RoleID = &roleID
		}
		if _, err := serviceaccounts.Update(c, d.Id(), request); err != nil {
			return diag.FromErr(err)
		}
	}
	return resourceServiceAccountRead(ctx, d, meta)
}

func resourceServiceAccountDelete(_ context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	err := serviceaccounts.Delete(meta.(*client.PhaseClient), d.Id())
	if err != nil && !client.IsNotFound(err) {
		return diag.FromErr(err)
	}
	d.SetId("")
	return nil
}

func serviceAccountImportState(_ context.Context, d *schema.ResourceData, _ any) ([]*schema.ResourceData, error) {
	if d.Id() == "" {
		return nil, fmt.Errorf("service account ID cannot be empty")
	}
	return []*schema.ResourceData{d}, nil
}
