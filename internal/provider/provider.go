package provider

import (
	"context"
	"terraform-provider-onecloud/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func New(version string) *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"api_token": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("1CLOUD_API_TOKEN", nil),
			},
			"api_url": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "https://api.1cloud.ru",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"onecloud_dns_domain": resourceDomain(),
			"onecloud_dns_record": resourceRecord(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"onecloud_dns_domain": dataSourceDomain(),
			"onecloud_dns_record": dataSourceRecord(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	token := d.Get("api_token").(string)
	baseURL := d.Get("api_url").(string)
	c, err := client.NewClient(baseURL, token)
	if err != nil {
		return nil, diag.FromErr(err)
	}
	return c, nil
}
