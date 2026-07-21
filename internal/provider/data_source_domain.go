package provider

import (
	"context"
	"strconv"

	"terraform-provider-onecloud/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceDomain() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDomainRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_delegate": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceDomainRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)
	domID := d.Get("id").(string)
	domain, err := c.GetDomain(domID)
	if err != nil {
		return diag.FromErr(err)
	}
	if domain == nil {
		return diag.Errorf("domain not found")
	}
	d.SetId(strconv.Itoa(domain.ID))
	d.Set("name", domain.Name)
	d.Set("state", domain.State)
	if domain.IsDelegate != nil {
		d.Set("is_delegate", strconv.FormatBool(*domain.IsDelegate))
	}
	return nil
}
