package datasource

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/phasehq/terraform-provider/internal/client"
	"github.com/phasehq/terraform-provider/internal/secrets"
)

func Secrets() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSecretsRead,
		Schema: map[string]*schema.Schema{
			"app_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the Phase App.",
			},
			"env": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The environment name.",
			},
			"path": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "/",
				Description: "The path to fetch secrets from.",
			},
			"key": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The key of a specific secret to fetch.",
			},
			"tags": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List of tags to filter secrets by. Multiple tags are combined with OR logic.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"secrets": {
				Type:      schema.TypeMap,
				Computed:  true,
				Sensitive: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceSecretsRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*client.PhaseClient)

	appID := d.Get("app_id").(string)
	env := d.Get("env").(string)
	path := d.Get("path").(string)
	key := d.Get("key").(string)

	// Handle tags if present
	var tagsFilter string
	if v, ok := d.GetOk("tags"); ok {
		tags := make([]string, 0)
		for _, tag := range v.([]any) {
			tags = append(tags, tag.(string))
		}
		if len(tags) > 0 {
			tagsFilter = strings.Join(tags, ",")
		}
	}

	// Determine if we're fetching all secrets
	fetchingAll := path == ""

	secrets, err := secrets.ReadSecret(client, appID, env, key, tagsFilter)
	if err != nil {
		return diag.FromErr(err)
	}

	secretMap := make(map[string]string)
	for _, secret := range secrets {
		if fetchingAll || secret.Path == path {
			if secret.Override != nil && secret.Override.IsActive {
				secretMap[secret.Key] = secret.Override.Value
			} else {
				secretMap[secret.Key] = secret.Value
			}
		}
	}

	if err := d.Set("secrets", secretMap); err != nil {
		return diag.FromErr(err)
	}

	// Set the path in the state
	if err := d.Set("path", path); err != nil {
		return diag.FromErr(err)
	}

	// Generate a unique ID for the data source
	d.SetId(fmt.Sprintf("%s-%s-%s-%s-%s", appID, env, path, key, tagsFilter))

	return nil
}
