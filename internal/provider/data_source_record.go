package provider

import (
	"context"
	"strconv"

	"terraform-provider-onecloud/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceRecord() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRecordRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID записи.",
			},
			"domain_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID домена, к которому принадлежит запись (может быть пустым).",
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"content": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ttl": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"priority": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"mnemonic_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ext_host_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"service": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"proto": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"weight": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"port": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"target": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"text": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceRecordRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)
	recID := d.Get("id").(string)

	rec, err := c.GetRecord(recID)
	if err != nil {
		return diag.FromErr(err)
	}
	if rec == nil {
		return diag.Errorf("запись с ID %s не найдена. Возможно, используйте ресурс с импортом, указав domain_id:record_id.", recID)
	}

	d.SetId(strconv.Itoa(rec.ID))
	d.Set("type", rec.TypeRecord)
	d.Set("name", rec.HostName)
	d.Set("ttl", rec.TTL)

	switch rec.TypeRecord {
	case "A", "AAAA":
		d.Set("content", rec.IP)
	case "CNAME":
		d.Set("content", rec.MnemonicName)
	case "MX":
		d.Set("content", rec.ExtHostName)
	case "NS":
		d.Set("content", rec.ExtHostName)
	case "SRV":
		d.Set("content", rec.Target)
	case "TXT":
		d.Set("content", rec.Text)
	}

	if rec.Priority != "" {
		p, _ := strconv.Atoi(rec.Priority)
		d.Set("priority", p)
	}
	d.Set("mnemonic_name", rec.MnemonicName)
	d.Set("ext_host_name", rec.ExtHostName)
	d.Set("service", rec.Service)
	d.Set("proto", rec.Proto)
	if rec.Weight != "" {
		w, _ := strconv.Atoi(rec.Weight)
		d.Set("weight", w)
	}
	if rec.Port != "" {
		p, _ := strconv.Atoi(rec.Port)
		d.Set("port", p)
	}
	d.Set("target", rec.Target)
	d.Set("text", rec.Text)

	return nil
}
