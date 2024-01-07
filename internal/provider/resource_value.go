package provider

import (
	"context"
	"encoding/json"
	"regexp"

	"github.com/ca-irvine/terraform-provider-edge/internal/model"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource = &ValueResource{}
)

func NewValueResource() resource.Resource {
	return &ValueResource{}
}

type ValueResource struct {
	c *config
}

type (
	valueResourceModel struct {
		ID             types.String                     `tfsdk:"id"`
		ValueID        types.String                     `tfsdk:"value_id"`
		Description    types.String                     `tfsdk:"description"`
		Enabled        types.Bool                       `tfsdk:"enabled"`
		DefaultVariant types.String                     `tfsdk:"default_variant"`
		BooleanValue   []valueResourceBooleanValueModel `tfsdk:"boolean_value"`
		StringValue    []valueResourceStringValueModel  `tfsdk:"string_value"`
		JSONValue      []valueResourceJSONValueModel    `tfsdk:"json_value"`
		IntegerValue   []valueResourceIntegerValueModel `tfsdk:"integer_value"`
		Targeting      []valueResourceTargetingModel    `tfsdk:"targeting"`
		Test           []valueResourceTestModel         `tfsdk:"test"`
	}

	valueResourceBooleanValueModel struct {
		Variant types.String `tfsdk:"variant"`
		Value   types.Bool   `tfsdk:"value"`
	}

	valueResourceStringValueModel struct {
		Variant types.String `tfsdk:"variant"`
		Value   types.String `tfsdk:"value"`
	}

	valueResourceJSONValueModel struct {
		Variant   types.String                  `tfsdk:"variant"`
		Value     types.String                  `tfsdk:"value"`
		Transform []valueResourceTransformModel `tfsdk:"transform"`
	}

	valueResourceIntegerValueModel struct {
		Variant types.String `tfsdk:"variant"`
		Value   types.Int64  `tfsdk:"value"`
	}

	valueResourceTargetingModel struct {
		Variant types.String `tfsdk:"variant"`
		Spec    types.String `tfsdk:"spec"`
		Expr    types.String `tfsdk:"expr"`
	}

	valueResourceTestModel struct {
		Variables types.String `tfsdk:"variables"`
		Expected  types.String `tfsdk:"expected"`
	}

	valueResourceTransformModel struct {
		Spec types.String `tfsdk:"spec"`
		Expr types.String `tfsdk:"expr"`
	}
)

func (v *ValueResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	value, err := v.c.GetValue(ctx, req.ID)
	if err != nil {
		resp.Diagnostics.AddError("Error get value", err.Error())
		return
	}

	state := valueState(value)
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (v *ValueResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_value"
}

func (v *ValueResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Edge value resource.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Computed ID.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"value_id": schema.StringAttribute{
				Description: "The ID of this Value.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"description": schema.StringAttribute{
				Optional: true,
			},
			"enabled": schema.BoolAttribute{
				Required: true,
			},
			"default_variant": schema.StringAttribute{
				Required: true,
			},
		},
		Blocks: map[string]schema.Block{
			"boolean_value": schema.ListNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"variant": schema.StringAttribute{
							Required: true,
						},
						"value": schema.BoolAttribute{
							Required: true,
						},
					},
				},
			},
			"string_value": schema.ListNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"variant": schema.StringAttribute{
							Required: true,
						},
						"value": schema.StringAttribute{
							Required: true,
						},
					},
				},
			},
			"json_value": schema.ListNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"variant": schema.StringAttribute{
							Required: true,
						},
						"value": schema.StringAttribute{
							Required: true,
							Validators: []validator.String{
								stringvalidator.RegexMatches(
									regexp.MustCompile(`^{.+}$`),
									"Must be map object, not array",
								),
							},
						},
					},
					Blocks: map[string]schema.Block{
						"transform": schema.ListNestedBlock{
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"spec": schema.StringAttribute{
										Optional: true,
									},
									"expr": schema.StringAttribute{
										Required: true,
									},
								},
							},
						},
					},
				},
			},
			"integer_value": schema.ListNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"variant": schema.StringAttribute{
							Required: true,
						},
						"value": schema.Int64Attribute{
							Required: true,
						},
					},
				},
			},
			"targeting": schema.ListNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"variant": schema.StringAttribute{
							Required: true,
						},
						"spec": schema.StringAttribute{
							Optional: true,
						},
						"expr": schema.StringAttribute{
							Required: true,
						},
					},
				},
			},
			"test": schema.ListNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"variables": schema.StringAttribute{
							Required: true,
						},
						"expected": schema.StringAttribute{
							Required: true,
						},
					},
				},
			},
		},
	}
}

func (v *valueResourceModel) value() (*model.Value, error) {
	variants := model.ValueVariants{}
	for _, val := range v.BooleanValue {
		variants[val.Variant.ValueString()] = model.ValueEvaluation{
			BooleanValue: &model.ValueBooleanValue{
				Value: val.Value.ValueBool(),
			},
		}
	}
	for _, val := range v.StringValue {
		variants[val.Variant.ValueString()] = model.ValueEvaluation{
			StringValue: &model.ValueStringValue{
				Value: val.Value.ValueString(),
			},
		}
	}
	for _, val := range v.JSONValue {
		m := make(map[string]any)
		err := json.Unmarshal([]byte(val.Value.ValueString()), &m)
		if err != nil {
			return nil, err
		}
		transforms := make([]*model.ValueTransform, 0, len(val.Transform))
		for _, t := range val.Transform {
			transforms = append(transforms, &model.ValueTransform{
				Spec: model.ValueTransformSpecFrom(t.Spec.ValueString()),
				Expr: t.Expr.ValueString(),
			})
		}
		variants[val.Variant.ValueString()] = model.ValueEvaluation{
			JSONValue: &model.ValueJSONValue{
				Value:      m,
				Transforms: transforms,
			},
		}
	}
	for _, val := range v.IntegerValue {
		variants[val.Variant.ValueString()] = model.ValueEvaluation{
			IntegerValue: &model.ValueIntegerValue{
				Value: val.Value.ValueInt64(),
			},
		}
	}

	rules := make([]model.ValueTargetingRule, 0, len(v.Targeting))
	for _, t := range v.Targeting {
		rules = append(rules, model.ValueTargetingRule{
			Variant: t.Variant.ValueString(),
			Spec:    model.ValueTargetingRuleSpecFrom(t.Spec.ValueString()),
			Expr:    t.Expr.ValueString(),
		})
	}

	tests := make([]*model.EvaluationTest, 0, len(v.Test))
	for _, t := range v.Test {
		m := make(map[string]any)
		err := json.Unmarshal([]byte(t.Variables.ValueString()), &m)
		if err != nil {
			return nil, err
		}
		tests = append(tests, &model.EvaluationTest{
			Variables: m,
			Expected:  t.Expected.ValueString(),
		})
	}
	value := &model.Value{
		ID:             v.ValueID.ValueString(),
		Enabled:        v.Enabled.ValueBool(),
		Description:    v.Description.ValueString(),
		DefaultVariant: v.DefaultVariant.ValueString(),
		Variants:       variants,
		Targeting: model.ValueTargeting{
			Rules: rules,
		},
		Tests: tests,
	}
	return value, nil
}

func valueState(v *model.Value) *valueResourceModel {
	var (
		bools []valueResourceBooleanValueModel
		strs  []valueResourceStringValueModel
		jsons []valueResourceJSONValueModel
		ints  []valueResourceIntegerValueModel
	)

	for k, val := range v.Variants {
		if val.BooleanValue != nil {
			if bools == nil {
				bools = make([]valueResourceBooleanValueModel, 0, len(v.Variants))
			}
			bools = append(bools, valueResourceBooleanValueModel{
				Variant: types.StringValue(k),
				Value:   types.BoolValue(val.BooleanValue.Value),
			})
		}
		if val.StringValue != nil {
			if strs == nil {
				strs = make([]valueResourceStringValueModel, 0, len(v.Variants))
			}
			strs = append(strs, valueResourceStringValueModel{
				Variant: types.StringValue(k),
				Value:   types.StringValue(val.StringValue.Value),
			})
		}
		if val.JSONValue != nil {
			if jsons == nil {
				jsons = make([]valueResourceJSONValueModel, 0, len(v.Variants))
			}
			b, _ := json.Marshal(val.JSONValue.Value)
			jsons = append(jsons, valueResourceJSONValueModel{
				Variant: types.StringValue(k),
				Value:   types.StringValue(string(b)),
			})
		}
		if val.IntegerValue != nil {
			if ints == nil {
				ints = make([]valueResourceIntegerValueModel, 0, len(v.Variants))
			}
			ints = append(ints, valueResourceIntegerValueModel{
				Variant: types.StringValue(k),
				Value:   types.Int64Value(val.IntegerValue.Value),
			})
		}
	}

	targeting := make([]valueResourceTargetingModel, 0, len(v.Targeting.Rules))
	for _, t := range v.Targeting.Rules {
		targeting = append(targeting, valueResourceTargetingModel{
			Variant: types.StringValue(t.Variant),
			Spec:    types.StringValue(model.TFValueTargetingRuleSpec(t.Spec)),
			Expr:    types.StringValue(t.Expr),
		})
	}

	tests := make([]valueResourceTestModel, 0, len(v.Tests))
	for _, t := range v.Tests {
		b, _ := json.Marshal(t.Variables)
		tests = append(tests, valueResourceTestModel{
			Variables: types.StringValue(string(b)),
			Expected:  types.StringValue(t.Expected),
		})
	}

	return &valueResourceModel{
		ID:             types.StringValue(v.ID),
		ValueID:        types.StringValue(v.ID),
		Description:    types.StringValue(v.Description),
		Enabled:        types.BoolValue(v.Enabled),
		DefaultVariant: types.StringValue(v.DefaultVariant),
		BooleanValue:   bools,
		StringValue:    strs,
		JSONValue:      jsons,
		IntegerValue:   ints,
		Targeting:      targeting,
		Test:           tests,
	}
}

func (v *ValueResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan valueResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	value, err := plan.value()
	if err != nil {
		resp.Diagnostics.AddError("Error creating value", "Invalid Attribute(s): "+err.Error())
		return
	}

	value, err = v.c.CreateValue(ctx, value)
	if err != nil {
		resp.Diagnostics.AddError("Error creating value", err.Error())
		return
	}

	plan.ID = types.StringValue(value.ID)
	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

func (v *ValueResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state valueResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (v *ValueResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan valueResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	value, err := plan.value()
	if err != nil {
		resp.Diagnostics.AddError("Error updating value", "Invalid Attribute(s): "+err.Error())
		return
	}

	_, err = v.c.UpdateValue(ctx, value)
	if err != nil {
		resp.Diagnostics.AddError("Error updating value", err.Error())
		return
	}

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

func (v *ValueResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state valueResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := v.c.DeleteValue(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting value", err.Error())
		return
	}

	resp.State.RemoveResource(ctx)
}

func (v *ValueResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	v.c = req.ProviderData.(*config)
}
