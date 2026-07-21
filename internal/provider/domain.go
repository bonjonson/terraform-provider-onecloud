package provider

import (
	"context"
	"strconv"

	"terraform-provider-onecloud/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceDomain() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDomainCreate,
		ReadContext:   resourceDomainRead,
		DeleteContext: resourceDomainDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"migrate": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				ForceNew: true,
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

func resourceDomainCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)
	name := d.Get("name").(string)
	migrate := d.Get("migrate").(bool)

	domain, err := c.CreateDomain(name, migrate)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(strconv.Itoa(domain.ID))
	d.Set("name", domain.Name)
	d.Set("state", domain.State)
	if domain.IsDelegate != nil {
		d.Set("is_delegate", strconv.FormatBool(*domain.IsDelegate))
	}
	return resourceDomainRead(ctx, d, m)
}

func resourceDomainRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)
	domain, err := c.GetDomain(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	if domain == nil {
		d.SetId("")
		return nil
	}
	d.Set("name", domain.Name)
	d.Set("state", domain.State)
	if domain.IsDelegate != nil {
		d.Set("is_delegate", strconv.FormatBool(*domain.IsDelegate))
	}
	return nil
}

func resourceDomainDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)
	if err := c.DeleteDomain(d.Id()); err != nil {
		return diag.FromErr(err)
	}
	return nil
}
