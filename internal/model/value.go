package model

type Value struct {
	ID             string            `json:"id"`
	Enabled        bool              `json:"enabled"`
	Description    string            `json:"description"`
	DefaultVariant string            `json:"defaultVariant"`
	Variants       ValueVariants     `json:"variants"`
	Targeting      ValueTargeting    `json:"targeting"`
	CreateTime     string            `json:"createTime,omitempty"`
	UpdateTime     string            `json:"updateTime,omitempty"`
	Tests          []*EvaluationTest `json:"tests,omitempty"`
}

type (
	ValueVariants map[string]ValueEvaluation

	ValueEvaluation struct {
		BooleanValue *ValueBooleanValue `json:"booleanValue"`
		StringValue  *ValueStringValue  `json:"stringValue"`
		JSONValue    *ValueJSONValue    `json:"jsonValue"`
		IntegerValue *ValueIntegerValue `json:"integerValue"`
	}

	ValueBooleanValue struct {
		Value bool `json:"value"`
	}

	ValueStringValue struct {
		Value string `json:"value"`
	}

	ValueJSONValue struct {
		Value      map[string]any    `json:"value"`
		Transforms []*ValueTransform `json:"transforms,omitempty"`
	}

	ValueIntegerValue struct {
		Value int64 `json:"value"`
	}
)

type EvaluationTest struct {
	Variables map[string]any `json:"variables"`
	Expected  string         `json:"expected"`
}

type ValueTargeting struct {
	Rules []ValueTargetingRule `json:"rules"`
}

type ValueTargetingRule struct {
	Variant string                 `json:"variant"`
	Spec    ValueTargetingRuleSpec `json:"spec"`
	Expr    string                 `json:"expr"`
}

type ValueTargetingRuleSpec int32

const (
	ValueTargetingRuleSpecCEL ValueTargetingRuleSpec = iota
	ValueTargetingRuleSpecJsonLogic
)

func ValueTargetingRuleSpecFrom(v string) ValueTargetingRuleSpec {
	switch v {
	case "cel":
		return ValueTargetingRuleSpecCEL
	case "json":
		return ValueTargetingRuleSpecJsonLogic
	default:
		return ValueTargetingRuleSpecCEL
	}
}

func TFValueTargetingRuleSpec(v ValueTargetingRuleSpec) string {
	switch v {
	case ValueTargetingRuleSpecCEL:
		return "cel"
	case ValueTargetingRuleSpecJsonLogic:
		return "json"
	default:
		return "cel"
	}
}

type ValueTransform struct {
	Spec ValueTransformSpec `json:"spec"`
	Expr string             `json:"expr"`
}

type ValueTransformSpec int32

const (
	ValueTransformSpecCEL ValueTransformSpec = iota
)

func ValueTransformSpecFrom(v string) ValueTransformSpec {
	switch v {
	case "cel":
		return ValueTransformSpecCEL
	default:
		return ValueTransformSpecCEL
	}
}

func TFValueTransformSpec(v int32) string {
	switch ValueTransformSpec(v) {
	case ValueTransformSpecCEL:
		return "cel"
	default:
		return "cel"
	}
}

type GetValueRequest struct {
	ID string `json:"id"`
}

type DeleteValueRequest struct {
	ID string `json:"id"`
}
