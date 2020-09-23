package menandmice

import (
	"strconv"
	"terraform-provider-menandmice/diag"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// schema for DNSZone resource
func DataSourceDNSZone() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDNSZoneRead,
		Schema: map[string]*schema.Schema{
			// TODO add more search criteria

			"ref": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				// TODO ref or name and authority
			},
			// TODO add  .
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"dynamic": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			// TODO add  .
			"adintegrated": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},

			"dnsviewref": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"dnsviewrefs": &schema.Schema{
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed: true,
			},
			"type": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"authority": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"dnssecsigned": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			"kskids": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"zskids": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			// TODO "customProperties": &schema.Schema{
			// 	Type:     ?
			// 	Computed: true,
			// }

			"created": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"lastmodified": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"displayname": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceDNSZoneRead(d *schema.ResourceData, m interface{}) error {

	var diags diag.Diagnostics
	c := m.(*Mmclient)

	filter := make(map[string]string)

	filter["name"] = d.Get("name").(string)
	filter["authority"] = d.Get("authority").(string)

	err, dnszones := c.FindDNSZone(filter)

	if err != nil {
		return diag.FromErr(err)
	}
	switch {
	case len(dnszones) <= 0:
		return diag.Errorf("no DNSZOnes found matching you criteria")
	case len(dnszones) > 1:
		return diag.Errorf("%v DNSZOnes found matching you criteria, but should be only 1", len(dnszones))
	}

	writeDNSZoneSchema(d, dnszones[0])
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags

}
