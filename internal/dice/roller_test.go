package dice

import (
	"reflect"
	"testing"
)

func TestParseExpression(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		want  Expression
	}{
		{
			name:  "simple",
			input: "2d6",
			want: Expression{
				Dice: []DiceTerm{
					{Count: 2, Sides: 6, Sign: 1},
				},
			},
		},
		{
			name:  "implicit count",
			input: "d6",
			want: Expression{
				Dice: []DiceTerm{
					{Count: 1, Sides: 6, Sign: 1},
				},
			},
		},
		{
			name:  "multiple dice and modifier",
			input: "2d6 + 1d4 - 3",
			want: Expression{
				Dice: []DiceTerm{
					{Count: 2, Sides: 6, Sign: 1},
					{Count: 1, Sides: 4, Sign: 1},
				},
				Modifier: -3,
			},
		},
		{
			name:  "negative dice term",
			input: "d8 - 2d4 + 5",
			want: Expression{
				Dice: []DiceTerm{
					{Count: 1, Sides: 8, Sign: 1},
					{Count: 2, Sides: 4, Sign: -1},
				},
				Modifier: 5,
			},
		},
		{
			name:  "uppercase with spaces",
			input: "  3D8   -   2 ",
			want: Expression{
				Dice: []DiceTerm{
					{Count: 3, Sides: 8, Sign: 1},
				},
				Modifier: -2,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := ParseExpression(tt.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("unexpected result: %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestParseExpressionErrors(t *testing.T) {
	t.Parallel()

	tests := []string{
		"",
		"d0",
		"0d6",
		"2dx",
		"2d6 ++ 1",
		"5",
	}

	for _, input := range tests {
		input := input
		t.Run(input, func(t *testing.T) {
			t.Parallel()
			if _, err := ParseExpression(input); err == nil {
				t.Fatalf("expected error for %q", input)
			}
		})
	}
}

func TestRoll(t *testing.T) {
	t.Parallel()

	expr := Expression{
		Dice: []DiceTerm{
			{Count: 3, Sides: 6, Sign: 1},
			{Count: 1, Sides: 4, Sign: -1},
		},
		Modifier: 2,
	}

	result, err := Roll(expr)
	if err != nil {
		t.Fatalf("Roll returned error: %v", err)
	}

	expectedRolls := 4
	if len(result.Rolls) != expectedRolls {
		t.Fatalf("expected %d rolls, got %d", expectedRolls, len(result.Rolls))
	}

	total := expr.Modifier
	index := 0
	for _, term := range expr.Dice {
		for i := 0; i < term.Count; i++ {
			roll := result.Rolls[index]
			index++

			if term.Sign == 1 {
				if roll < 1 || roll > term.Sides {
					t.Fatalf("positive roll out of range: %d", roll)
				}
			} else {
				if roll > -1 || roll < -term.Sides {
					t.Fatalf("negative roll out of range: %d", roll)
				}
			}
			total += roll
		}
	}

	if result.Total != total {
		t.Fatalf("unexpected total %d, want %d", result.Total, total)
	}
}

func TestRollValidation(t *testing.T) {
	t.Parallel()

	_, err := Roll(Expression{})
	if err == nil {
		t.Fatalf("expected error for empty expression")
	}

	_, err = Roll(Expression{
		Dice: []DiceTerm{{Count: 0, Sides: 6, Sign: 1}},
	})
	if err == nil {
		t.Fatalf("expected error for zero dice count")
	}

	_, err = Roll(Expression{
		Dice: []DiceTerm{{Count: 1, Sides: 0, Sign: 1}},
	})
	if err == nil {
		t.Fatalf("expected error for invalid dice sides")
	}
}

func TestRollSingleSideDie(t *testing.T) {
	t.Parallel()

	expr := Expression{
		Dice: []DiceTerm{{Count: 2, Sides: 1, Sign: 1}},
	}
	result, err := Roll(expr)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, roll := range result.Rolls {
		if roll != 1 {
			t.Fatalf("expected roll 1, got %d", roll)
		}
	}
	expectedTotal := expr.Modifier + expr.Dice[0].Count
	if result.Total != expectedTotal {
		t.Fatalf("expected total %d, got %d", expectedTotal, result.Total)
	}
}
