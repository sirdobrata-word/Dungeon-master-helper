package dice

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"strings"
)

// Expression describes a dice expression consisting of dice terms and a constant modifier.
type Expression struct {
	Dice     []DiceTerm
	Modifier int
}

// DiceTerm represents a single dice group (e.g. "+2d6" or "-d4").
type DiceTerm struct {
	Count int
	Sides int
	Sign  int
}

// Result contains the detailed output of a dice roll.
type Result struct {
	Expression Expression
	Rolls      []int
	Total      int
}

// ParseExpression converts a textual dice expression into a structured representation.
func ParseExpression(input string) (Expression, error) {
	trimmed := strings.TrimSpace(input)
	if trimmed == "" {
		return Expression{}, fmt.Errorf("invalid expression: %q", trimmed)
	}

	tokens, err := splitTerms(trimmed)
	if err != nil {
		return Expression{}, err
	}

	expr := Expression{}
	for _, raw := range tokens {
		token := strings.TrimSpace(raw)
		if token == "" {
			return Expression{}, fmt.Errorf("invalid expression: %q", trimmed)
		}

		sign := 1
		switch token[0] {
		case '+':
			token = strings.TrimSpace(token[1:])
		case '-':
			sign = -1
			token = strings.TrimSpace(token[1:])
		}
		if token == "" {
			return Expression{}, fmt.Errorf("invalid expression: %q", trimmed)
		}

		idx := strings.IndexAny(token, "dD")
		if idx >= 0 {
			countStr := strings.TrimSpace(token[:idx])
			sidesStr := strings.TrimSpace(token[idx+1:])
			if sidesStr == "" {
				return Expression{}, fmt.Errorf("invalid expression: %q", trimmed)
			}

			count := 1
			if countStr != "" {
				value, err := parsePositive(countStr, "dice count")
				if err != nil {
					return Expression{}, err
				}
				count = value
			}

			sides, err := parsePositive(sidesStr, "dice sides")
			if err != nil {
				return Expression{}, err
			}

			expr.Dice = append(expr.Dice, DiceTerm{
				Count: count,
				Sides: sides,
				Sign:  sign,
			})
			continue
		}

		value, err := parsePositive(token, "modifier")
		if err != nil {
			return Expression{}, err
		}
		expr.Modifier += sign * value
	}

	if len(expr.Dice) == 0 {
		return Expression{}, fmt.Errorf("invalid expression: %q", trimmed)
	}

	return expr, nil
}

func splitTerms(input string) ([]string, error) {
	var tokens []string
	var current strings.Builder

	for _, r := range input {
		if r == '+' || r == '-' {
			if current.Len() == 0 {
				current.WriteRune(r)
				continue
			}
			tokens = append(tokens, current.String())
			current.Reset()
			current.WriteRune(r)
			continue
		}
		current.WriteRune(r)
	}

	if current.Len() == 0 {
		return nil, fmt.Errorf("invalid expression: %q", input)
	}
	tokens = append(tokens, current.String())

	return tokens, nil
}

func parsePositive(raw, field string) (int, error) {
	value := 0
	for _, ch := range raw {
		if ch < '0' || ch > '9' {
			return 0, fmt.Errorf("%s must be numeric", field)
		}
		value = value*10 + int(ch-'0')
		if value > 10_000 {
			return 0, fmt.Errorf("%s is too large", field)
		}
	}
	if value <= 0 {
		return 0, fmt.Errorf("%s must be positive", field)
	}
	return value, nil
}

// Roll executes a dice expression using a cryptographically secure RNG.
func Roll(expr Expression) (Result, error) {
	if len(expr.Dice) == 0 {
		return Result{}, errors.New("expression must include at least one dice term")
	}

	var rolls []int
	total := expr.Modifier
	for _, term := range expr.Dice {
		if term.Count <= 0 {
			return Result{}, errors.New("dice count must be positive")
		}
		if term.Sides < 1 {
			return Result{}, errors.New("dice must have at least one side")
		}
		sign := term.Sign
		if sign != 1 && sign != -1 {
			return Result{}, errors.New("invalid dice term")
		}

		for i := 0; i < term.Count; i++ {
			value, err := rollDie(term.Sides)
			if err != nil {
				return Result{}, err
			}
			adjusted := value * sign
			rolls = append(rolls, adjusted)
			total += adjusted
		}
	}

	return Result{
		Expression: expr,
		Rolls:      rolls,
		Total:      total,
	}, nil
}

func rollDie(sides int) (int, error) {
	max := big.NewInt(int64(sides))
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return 0, fmt.Errorf("rng failure: %w", err)
	}
	return int(n.Int64()) + 1, nil
}
