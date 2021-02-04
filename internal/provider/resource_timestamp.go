package provider

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceTimestamp() *schema.Resource {
	return &schema.Resource{
		Description: "The resource `pinned_timestamp` generates a timestamp when one of the trigger values changes.\n" +
			"\n" +
			"This resource can be used in conjunction with resources that have the `create_before_destroy` " +
			"lifecycle flag set, to avoid conflicts with unique names during the brief period where both the " +
			"old and new resources exist concurrently.",
		Create: CreateTimestamp,
		Read:   schema.Noop,
		Delete: schema.RemoveFromState,
		Importer: &schema.ResourceImporter{
			State: ImportTimestamp,
		},

		Schema: map[string]*schema.Schema{
			"triggers": {
				Description: "Arbitrary map of values that, when changed, will trigger recreation of " +
					"resource. See [the main provider documentation](../index.html) for more information.",
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
			},

			"format": {
				Description: "A custom date/time format - must be written in Go syntax.  " +
					"e.g. see: https://yourbasic.org/golang/format-parse-string-time-date-example/",
				Type:     schema.TypeString,
				Default:  "2006-01-02 15:04:05",
				Optional: true,
				ForceNew: true,
			},

			"timestamp": {
				Description: "The generated timestamp.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func CreateTimestamp(d *schema.ResourceData, meta interface{}) error {
	format := d.Get("format").(string)
	now := time.Now()

	value := now.Format(format)

	d.Set("timestamp", value)
	d.SetId(value)

	return nil
}

func ImportTimestamp(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), ",")
	if len(parts) != 1 && len(parts) != 2 {
		return nil, fmt.Errorf("Invalid import usage: expecting {timestamp} or {timestamp},{format}")
	}

	d.Set("timestamp", parts[0])

	if len(parts) == 2 {
		d.Set("format", parts[1])
	}

	d.SetId(parts[0])
	return []*schema.ResourceData{d}, nil
}
