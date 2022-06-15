package core

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRules(t *testing.T) {

	type testCase struct {
		rule             map[string]any
		context          EvalContext
		shouldBeAccepted bool
	}

	var tests = []testCase{
		{
			rule: map[string]any{
				"$or": []any{
					map[string]any{
						"$gt": []any{
							"$balance", float64(0),
						},
					},
					map[string]any{
						"$eq": []any{
							map[string]any{
								"$meta": "approved",
							},
							"yes",
						},
					},
				},
			},
			context: EvalContext{
				Variables: map[string]any{
					"balance": float64(-10),
				},
				Metadata: map[string]json.RawMessage{
					"approved": json.RawMessage("yes"),
				},
			},
			shouldBeAccepted: true,
		},
		{
			rule: map[string]any{
				"$or": []any{
					map[string]any{
						"$gte": []any{
							"$balance", float64(0),
						},
					},
					map[string]any{
						"$lte": []any{
							"$balance", float64(0),
						},
					},
				},
			},
			context: EvalContext{
				Variables: map[string]any{
					"balance": float64(-100),
				},
				Metadata: map[string]json.RawMessage{},
			},
			shouldBeAccepted: true,
		},
		{
			rule: map[string]interface{}{
				"$lt": []interface{}{
					"$balance", float64(0),
				},
			},
			context: EvalContext{
				Variables: map[string]interface{}{
					"balance": float64(100),
				},
				Metadata: map[string]json.RawMessage{},
			},
			shouldBeAccepted: false,
		},
		{
			rule: map[string]interface{}{
				"$lte": []interface{}{
					"$balance", float64(0),
				},
			},
			context: EvalContext{
				Variables: map[string]interface{}{
					"balance": float64(0),
				},
				Metadata: map[string]json.RawMessage{},
			},
			shouldBeAccepted: true,
		},
		{
			rule: map[string]interface{}{
				"$and": []interface{}{
					map[string]interface{}{
						"$gt": []interface{}{
							"$balance", float64(0),
						},
					},
					map[string]interface{}{
						"$eq": []interface{}{
							map[string]interface{}{
								"$meta": "approved",
							},
							"yes",
						},
					},
				},
			},
			context: EvalContext{
				Variables: map[string]interface{}{
					"balance": float64(10),
				},
				Metadata: map[string]json.RawMessage{
					"approved": json.RawMessage("no"),
				},
			},
			shouldBeAccepted: false,
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("test%d", i), func(t *testing.T) {
			r, err := ParseRuleExpr(test.rule)
			assert.NoError(t, err)
			assert.Equal(t, test.shouldBeAccepted, r.Eval(test.context))
		})
	}

}
