package provider

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceInteger() *schema.Resource {
	return &schema.Resource{
		Description: "The resource `sequential_integer` generates sequential values from a given range, described " +
			"by the `min` and `max` attributes of a given resource.\n" +
			"\n" +
			"This resource can be used in conjunction with resources that have the `create_before_destroy` " +
			"lifecycle flag set, to avoid conflicts with unique names during the brief period where both the " +
			"old and new resources exist concurrently.",
		Create: CreateInteger,
		Read:   ReadInteger,
		Update: UpdateInteger,
		Delete: schema.RemoveFromState,
		Importer: &schema.ResourceImporter{
			State: ImportInteger,
		},

		Schema: map[string]*schema.Schema{
			"keepers": {
				Description: "Arbitrary map of values that, when changed, will trigger recreation of " +
					"resource. See [the main provider documentation](../index.html) for more information.",
				Type:     schema.TypeMap,
				Optional: true,
			},

			"min": {
				Description: "The minimum inclusive value of the range. Must be greater than Zero.",
				Type:        schema.TypeInt,
				Required:    true,
			},

			"max": {
				Description: "The maximum inclusive value of the range. If generation exceeds Max, it will reset to Min.",
				Type:        schema.TypeInt,
				Required:    true,
			},

			"value": {
				Description: "The sequential integer result.",
				Type:        schema.TypeInt,
				Computed:    true,
			},
		},
		CustomizeDiff: func(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error {
			old_min, new_min := d.GetChange("min")
			old_max, new_max := d.GetChange("max")
			_, new_res := d.GetChange("value")
			if new_min.(int) != old_min.(int) || new_max.(int) != old_max.(int) {
				if new_min.(int) > new_res.(int) || new_max.(int) < new_res.(int) {
					d.SetNewComputed("value")
				}
			}
			if d.HasChange("keepers") {
				d.SetNewComputed("value")
			}
			return nil
		},
		UseJSONNumber: true,
	}
}

func CreateInteger(d *schema.ResourceData, meta interface{}) error {
	min := d.Get("min").(int)
	max := d.Get("max").(int)

	if min <= 0 {
		return fmt.Errorf("Minimum value cannot be less than or equal to Zero")
	}
	if max <= min {
		return fmt.Errorf("Maximum value needs to be greater than minimum value")
	}

	value := min

	d.Set("value", value)
	d.SetId(strconv.Itoa(value))

	return nil
}

func ReadInteger(d *schema.ResourceData, m interface{}) error {
	min := d.Get("min").(int)
	max := d.Get("max").(int)

	if min <= 0 {
		return fmt.Errorf("Minimum value cannot be less than or equal to Zero")
	}
	if max <= min {
		return fmt.Errorf("Maximum value needs to be greater than minimum value")
	}

	value, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("Cannot read resource id [%s]: %s", d.Id(), err)
	}

	if value > max || value < min {
		value = min
	}

	d.Set("value", value)

	return nil
}

func UpdateInteger(d *schema.ResourceData, m interface{}) error {
	min := d.Get("min").(int)
	max := d.Get("max").(int)

	if min <= 0 {
		return fmt.Errorf("Minimum value cannot be less than or equal to Zero")
	}
	if max <= min {
		return fmt.Errorf("Maximum value needs to be greater than minimum value")
	}

	value, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("Cannot read resource id [%s]: %s", d.Id(), err)
	}

	value += 1
	if value > max || value < min {
		value = min
	}

	d.Set("value", value)
	d.SetId(strconv.Itoa(value))

	return nil
}

func ImportInteger(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), ",")
	if len(parts) != 3 {
		return nil, fmt.Errorf("Invalid import usage: expecting {value},{min},{max}")
	}

	min, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, errwrap.Wrapf("Error parsing \"min\": {{err}}", err)
	}
	if min <= 0 {
		return nil, fmt.Errorf("Minimum value cannot be less than or equal to Zero")
	}
	d.Set("min", min)

	max, err := strconv.Atoi(parts[2])
	if err != nil {
		return nil, errwrap.Wrapf("Error parsing \"max\": {{err}}", err)
	}
	if max <= min {
		return nil, fmt.Errorf("Maximum value needs to be greater than minimum value")
	}
	d.Set("max", max)

	result, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil, errwrap.Wrapf("Error parsing \"value\": {{err}}", err)
	}
	if result < min || result > max {
		return nil, errwrap.Wrapf("Value must be between Min and Max (inclusive)", err)
	}
	d.Set("value", result)

	return []*schema.ResourceData{d}, nil
}
