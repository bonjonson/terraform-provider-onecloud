package provider

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"terraform-provider-onecloud/client"
)

func resourceRecord() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRecordCreate,
		ReadContext:   resourceRecordRead,
		UpdateContext: resourceRecordUpdate,
		DeleteContext: resourceRecordDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceRecordImport,
		},
		Schema: map[string]*schema.Schema{
			"domain_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"A", "AAAA", "CNAME", "MX", "NS", "SRV", "TXT"}, false),
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"content": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ttl": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  3600,
			},
			"priority": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"mnemonic_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ext_host_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"service": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"proto": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"weight": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"port": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"target": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"text": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func recordFromResourceData(d *schema.ResourceData) *client.Record {
	rec := &client.Record{
		TypeRecord: d.Get("type").(string),
		HostName:   d.Get("name").(string),
		TTL:        d.Get("ttl").(int),
	}
	switch strings.ToUpper(rec.TypeRecord) {
	case "A", "AAAA":
		rec.IP = d.Get("content").(string)
	case "CNAME":
		rec.MnemonicName = d.Get("content").(string)
	case "MX":
		rec.ExtHostName = d.Get("content").(string)
	case "NS":
		rec.ExtHostName = d.Get("content").(string)
	case "SRV":
		rec.Target = d.Get("content").(string)
	case "TXT":
		rec.Text = d.Get("content").(string)
	}
	if v, ok := d.GetOk("priority"); ok {
		rec.Priority = strconv.Itoa(v.(int))
	}
	if v, ok := d.GetOk("mnemonic_name"); ok {
		rec.MnemonicName = v.(string)
	}
	if v, ok := d.GetOk("ext_host_name"); ok {
		rec.ExtHostName = v.(string)
	}
	if v, ok := d.GetOk("service"); ok {
		rec.Service = v.(string)
	}
	if v, ok := d.GetOk("proto"); ok {
		rec.Proto = v.(string)
	}
	if v, ok := d.GetOk("weight"); ok {
		rec.Weight = strconv.Itoa(v.(int))
	}
	if v, ok := d.GetOk("port"); ok {
		rec.Port = strconv.Itoa(v.(int))
	}
	if v, ok := d.GetOk("target"); ok {
		rec.Target = v.(string)
	}
	if v, ok := d.GetOk("text"); ok {
		rec.Text = v.(string)
	}
	return rec
}

func resourceDataFromRecord(d *schema.ResourceData, rec *client.Record) {
	d.Set("type", rec.TypeRecord)
	d.Set("name", rec.HostName)
	d.Set("ttl", rec.TTL)
	switch strings.ToUpper(rec.TypeRecord) {
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
}

func resourceRecordCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)
	domainID := d.Get("domain_id").(string)
	rec := recordFromResourceData(d)

	created, err := c.CreateRecord(domainID, rec)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(strconv.Itoa(created.ID))
	resourceDataFromRecord(d, created)

	return resourceRecordRead(ctx, d, m)
}

func resourceRecordRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)
	domainID := d.Get("domain_id").(string)
	rec, err := c.GetRecordByDomainID(domainID, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	if rec == nil {
		d.SetId("")
		return nil
	}
	resourceDataFromRecord(d, rec)
	return nil
}

func resourceRecordUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)
	domainID := d.Get("domain_id").(string)
	rec := recordFromResourceData(d)

	updated, err := c.UpdateRecord(domainID, d.Id(), rec)
	if err != nil {
		return diag.FromErr(err)
	}
	resourceDataFromRecord(d, updated)

	return resourceRecordRead(ctx, d, m)
}

func resourceRecordDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)
	domainID := d.Get("domain_id").(string)
	if err := c.DeleteRecord(domainID, d.Id()); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceRecordImport(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	parts := strings.SplitN(d.Id(), ":", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid import ID, expected domain_id:record_id")
	}
	d.Set("domain_id", parts[0])
	d.SetId(parts[1])
	return []*schema.ResourceData{d}, nil
}
