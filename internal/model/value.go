package model

type Value struct {
	ID             string            `json:"id"`
	Enabled        bool              `json:"enabled"`
	Description    string            `json:"description"`
	DefaultVariant string            `json:"default_variant"`
	Variants       ValueVariants     `json:"variants"`
	Targeting      *ValueTargeting   `json:"targeting"`
	CreateTime     int64             `json:"create_time,omitempty"`
	UpdateTime     int64             `json:"update_time,omitempty"`
	Tests          []*EvaluationTest `json:"tests,omitempty"`
}

type (
	ValueVariants map[string]ValueEvaluation

	ValueEvaluation struct {
		BooleanValue *ValueBooleanValue `json:"boolean_value"`
		StringValue  *ValueStringValue  `json:"string_value"`
		JSONValue    *ValueJSONValue    `json:"json_value"`
		IntegerValue *ValueIntegerValue `json:"integer_value"`
	}

	ValueBooleanValue struct {
		Value bool `json:"value"`
	}

	ValueStringValue struct {
		Value string `json:"value"`
	}

	ValueJSONValue struct {
		Value map[string]any `json:"value"`
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

type ValueTargetingRuleSpec string

const (
	ValueTargetingRuleSpecInvalid   ValueTargetingRuleSpec = "SPEC_INVALID"
	ValueTargetingRuleSpecCEL       ValueTargetingRuleSpec = "SPEC_GOOGLE_CEL"
	ValueTargetingRuleSpecJsonLogic ValueTargetingRuleSpec = "SPEC_JSON_LOGIC"
)

func ValueTargetingRuleSpecFrom(v string) ValueTargetingRuleSpec {
	switch v {
	case "cel":
		return ValueTargetingRuleSpecCEL
	case "json":
		return ValueTargetingRuleSpecJsonLogic
	default:
		return ValueTargetingRuleSpecInvalid
	}
}

type GetValueRequest struct {
	ID string `json:"id"`
}

type DeleteValueRequest struct {
	ID string `json:"id"`
}
