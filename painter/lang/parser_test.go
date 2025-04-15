package lang

import (
	"fmt"
	"strings"
	"testing"

	"github.com/roman-mazur/architecture-lab-3/painter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_parse_struct(t *testing.T) {
	tests := []struct {
		name        string
		command     string
		expectOp    painter.Operation
		expectError bool
	}{
		{
			name:        "rectangle",
			command:     "bgrect 0.25 0.25 0.75 0.75",
			expectOp:    &painter.BgRectOp{X1: 200, Y1: 200, X2: 600, Y2: 600},
		},
		{
			name:        "figure",
			command:     "figure 0.5 0.5",
			expectOp:    &painter.FigureOp{X: 400, Y: 400},
		},
		{
			name:        "move",
			command:     "move 0.1 0.1",
			expectOp:    &painter.MoveOp{X: 80, Y: 80},
		},
		{
			name:        "update",
			command:     "update",
			expectOp:    painter.UpdateOp,
		},
		{
			name:        "invalid command",
			command:     "invalidcommand",
			expectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			parser := &Parser{}
			ops, err := parser.Parse(strings.NewReader(tc.command))
			
			if tc.expectError {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			
			found := false
			for _, op := range ops {
				if assertOperationsEqual(t, tc.expectOp, op) {
					found = true
					break
				}
			}
			assert.True(t, found, "Operation %T not found in result", tc.expectOp)
		})
	}
}

func assertOperationsEqual(t *testing.T, expected, actual painter.Operation) bool {
	t.Helper()

	switch exp := expected.(type) {
	case *painter.BgRectOp:
		act, ok := actual.(*painter.BgRectOp)
		if !ok {
			return false
		}
		return exp.X1 == act.X1 && exp.Y1 == act.Y1 && exp.X2 == act.X2 && exp.Y2 == act.Y2
	
	case *painter.FigureOp:
		act, ok := actual.(*painter.FigureOp)
		if !ok {
			return false
		}
		return exp.X == act.X && exp.Y == act.Y
	
	case *painter.MoveOp:
		act, ok := actual.(*painter.MoveOp)
		if !ok {
			return false
		}
		return exp.X == act.X && exp.Y == act.Y
	
	case painter.Operation:

		return fmt.Sprintf("%v", exp) == fmt.Sprintf("%v", actual)
	
	default:
		return false
	}
}