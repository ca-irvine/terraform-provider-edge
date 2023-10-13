package model

import (
	"testing"
)

func TestValueTargetingRuleSpecFrom(t *testing.T) {
	t.Parallel()
	tests := []struct {
		v    string
		want ValueTargetingRuleSpec
	}{
		{
			v:    "cel",
			want: ValueTargetingRuleSpecCEL,
		},
		{
			v:    "json",
			want: ValueTargetingRuleSpecJsonLogic,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run("", func(t *testing.T) {
			t.Parallel()
			got := ValueTargetingRuleSpecFrom(tt.v)
			if got != tt.want {
				t.Fatalf("expected %d, but got %d", tt.want, got)
			}
		})
	}
}

func TestValueTransformSpecFrom(t *testing.T) {
	t.Parallel()
	tests := []struct {
		v    string
		want ValueTransformSpec
	}{
		{
			v:    "cel",
			want: ValueTransformSpecCEL,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run("", func(t *testing.T) {
			t.Parallel()
			got := ValueTransformSpecFrom(tt.v)
			if got != tt.want {
				t.Fatalf("expected %d, but got %d", tt.want, got)
			}
		})
	}
}
