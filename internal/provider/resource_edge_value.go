package provider

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ca-irvine/terraform-provider-edge/internal/model"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceEdgeValue() *schema.Resource {
	return &schema.Resource{
		Description:   "Edge value resource.",
		CreateContext: resourceValueCreate,
		ReadContext:   resourceValueRead,
		UpdateContext: resourceValueUpdate,
		DeleteContext: resourceValueDelete,

		Schema: map[string]*schema.Schema{
			"value_id": {
				Description: "Primary key of value to be resolved.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"description": {
				Description: "Value description. This is not exposed to resolver client.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"enabled": {
				Description: "Flag of data lifecycle. If set false, resolver does not resolve value.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
			},
			"default_variant": {
				Description: "Default variant name used by resolver.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"boolean_value": {
				Description: "Boolean value variant.",
				Type:        schema.TypeSet,
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"variant": {
							Type:     schema.TypeString,
							Required: true,
						},
						"value": {
							Type:     schema.TypeBool,
							Required: true,
						},
					},
				},
			},
			"string_value": {
				Description: "String value variant.",
				Type:        schema.TypeSet,
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"variant": {
							Type:     schema.TypeString,
							Required: true,
						},
						"value": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"json_value": {
				Description: "JSON value variant.",
				Type:        schema.TypeSet,
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"variant": {
							Type:     schema.TypeString,
							Required: true,
						},
						"value": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"targeting": {
				Description: "Value targeting expression. Google CEL is supported.",
				Type:        schema.TypeList,
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"variant": {
							Type:     schema.TypeString,
							Required: true,
						},
						"spec": {
							Description: "Google CEL: cel",
							Type:        schema.TypeString,
							Required:    true,
						},
						"expr": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"test": {
				Description: "Test value targeting expression. Google CEL is supported.",
				Type:        schema.TypeList,
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"variables": {
							Description: "Variables for test cases.",
							Type:        schema.TypeMap,
							Required:    true,
						},
						"expected": {
							Description: "Expected variant.",
							Type:        schema.TypeString,
							Required:    true,
						},
					},
				},
			},
		},
	}
}

func mapValueBooleanEvaluations(variants model.ValueVariants, v any) {
	set := v.(*schema.Set).List()
	for i := range set {
		m := set[i].(map[string]any)
		key := m["variant"].(string)
		value := m["value"].(bool)
		variants[key] = model.ValueEvaluation{BooleanValue: &model.ValueBooleanValue{Value: value}}
	}
}

func mapValueStringEvaluations(variants model.ValueVariants, v any) {
	set := v.(*schema.Set).List()
	for i := range set {
		m := set[i].(map[string]any)
		key := m["variant"].(string)
		value := m["value"].(string)
		variants[key] = model.ValueEvaluation{StringValue: &model.ValueStringValue{Value: value}}
	}
}

func mapValueJSONEvaluations(variants model.ValueVariants, v any) error {
	set := v.(*schema.Set).List()
	for i := range set {
		m := set[i].(map[string]any)
		key := m["variant"].(string)
		value := m["value"].(string)
		eval := new(model.ValueJSONValue)
		err := json.Unmarshal([]byte(value), &eval.Value)
		if err != nil {
			return err
		}
		variants[key] = model.ValueEvaluation{JSONValue: eval}
	}
	return nil
}

func buildValueVariants(d *schema.ResourceData) (model.ValueVariants, error) {
	variants := model.ValueVariants{}
	types := make([]bool, 0, 3)

	boolSet, hasBool := d.GetOk("boolean_value")
	if hasBool {
		mapValueBooleanEvaluations(variants, boolSet)
		types = append(types, hasBool)
	}

	stringSet, hasString := d.GetOk("string_value")
	if hasString {
		mapValueStringEvaluations(variants, stringSet)
		types = append(types, hasString)
	}

	jsonSet, hasJSON := d.GetOk("json_value")
	if hasJSON {
		err := mapValueJSONEvaluations(variants, jsonSet)
		if err != nil {
			return nil, err
		}
		types = append(types, hasJSON)
	}

	if len(types) != 1 {
		return nil, fmt.Errorf("one of `boolean_value` or `string_value` or `json_value` must be set")
	}

	return variants, nil
}

func buildValueTargeting(d *schema.ResourceData) (*model.ValueTargeting, error) {
	targetingList, ok := d.GetOk("targeting")
	if !ok {
		return &model.ValueTargeting{}, nil
	}

	rules := make([]model.ValueTargetingRule, 0, 5)
	list := targetingList.([]any)
	for i := range list {
		m := list[i].(map[string]any)
		variant := m["variant"].(string)
		spec := m["spec"].(string)
		expr := m["expr"].(string)
		specInt := model.ValueTargetingRuleSpecFrom(spec)
		rules = append(rules, model.ValueTargetingRule{
			Variant: variant,
			Spec:    specInt,
			Expr:    expr,
		})
	}

	return &model.ValueTargeting{Rules: rules}, nil
}

func buildEvaluationTests(d *schema.ResourceData) ([]*model.EvaluationTest, error) {
	test, ok := d.GetOk("test")
	if !ok {
		return []*model.EvaluationTest{}, nil
	}

	list := test.([]any)
	tests := make([]*model.EvaluationTest, 0, len(list))
	for i := range list {
		m := list[i].(map[string]any)
		tests = append(tests, &model.EvaluationTest{
			Variables: m["variables"].(map[string]any),
			Expected:  m["expected"].(string),
		})
	}

	return tests, nil
}

func resourceValueCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	id := d.Get("value_id").(string)
	description := d.Get("description").(string)
	enabled := d.Get("enabled").(bool)
	defaultVariant := d.Get("default_variant").(string)

	variants, err := buildValueVariants(d)
	if err != nil {
		return diag.Errorf("variant block %s", err)
	}

	targeting, err := buildValueTargeting(d)
	if err != nil {
		return diag.Errorf("targeting block %s", err)
	}

	tests, err := buildEvaluationTests(d)
	if err != nil {
		return diag.Errorf("test block %s", err)
	}

	value := &model.Value{
		ID:             id,
		Enabled:        enabled,
		Description:    description,
		DefaultVariant: defaultVariant,
		Variants:       variants,
		Targeting:      targeting,
		Tests:          tests,
	}

	client := meta.(*config)
	err = client.CreateValue(ctx, value)
	if err != nil {
		return diag.Errorf("create error %s", err)
	}

	d.SetId(id)

	tflog.Trace(ctx, "created a value resource")

	return nil
}

func resourceValueRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	id := d.Id()

	client := meta.(*config)
	_, err := client.GetValue(ctx, id)
	if err != nil {
		return diag.Errorf("get error %s", err)
	}

	d.SetId(id)

	tflog.Trace(ctx, "get a value resource")

	return nil
}

func resourceValueUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	id := d.Id()
	description := d.Get("description").(string)
	enabled := d.Get("enabled").(bool)
	defaultVariant := d.Get("default_variant").(string)

	variants, err := buildValueVariants(d)
	if err != nil {
		return diag.Errorf("variant block %s", err)
	}

	targeting, err := buildValueTargeting(d)
	if err != nil {
		return diag.Errorf("targeting block %s", err)
	}

	tests, err := buildEvaluationTests(d)
	if err != nil {
		return diag.Errorf("test block %s", err)
	}

	value := &model.Value{
		ID:             id,
		Enabled:        enabled,
		Description:    description,
		DefaultVariant: defaultVariant,
		Variants:       variants,
		Targeting:      targeting,
		Tests:          tests,
	}

	client := meta.(*config)
	err = client.UpdateValue(ctx, value)
	if err != nil {
		return diag.Errorf("update error %s", err)
	}

	d.SetId(id)

	tflog.Trace(ctx, "update a value resource")

	return nil
}

func resourceValueDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	id := d.Id()

	client := meta.(*config)

	if err := client.DeleteValue(ctx, id); err != nil {
		return diag.Errorf("delete error %s", err)
	}

	d.SetId(id)

	tflog.Trace(ctx, "deleted a value resource")

	return nil
}
