package model

type Value struct {
	ID             string          `json:"id"`
	Enabled        bool            `json:"enabled"`
	Description    string          `json:"description"`
	DefaultVariant string          `json:"defaultVariant"`
	Variants       ValueVariants   `json:"variants"`
	Targeting      *ValueTargeting `json:"targeting"`
}

type (
	ValueVariants map[string]ValueEvaluation

	ValueEvaluation struct {
		BooleanValue bool           `json:"booleanValue"`
		StringValue  string         `json:"stringValue"`
		JSONValue    map[string]any `json:"jsonValue"`
	}

	ValueBooleanValue struct {
		Value bool
	}
)

type ValueTargeting struct {
	Rules []ValueTargetingRule `json:"rules"`
}

type ValueTargetingRule struct {
	Variant string                 `json:"variant"`
	Spec    ValueTargetingRuleSpec `json:"spec"`
	Exp     string                 `json:"exp"`
}

type ValueTargetingRuleSpec int32

const (
	ValueTargetingRuleSpecInvalid ValueTargetingRuleSpec = iota
	ValueTargetingRuleSpecCEL
)

func ValueTargetingRuleSpecFrom(v string) ValueTargetingRuleSpec {
	switch v {
	case "cel":
		return ValueTargetingRuleSpecCEL
	default:
		return ValueTargetingRuleSpecInvalid
	}
}
