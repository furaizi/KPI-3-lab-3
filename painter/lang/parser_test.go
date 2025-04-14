package lang

import (
	"github.com/roman-mazur/architecture-lab-3/painter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func TestParser(t *testing.T) {

	tests := []struct {
		name  string
		input string
		op    painter.Operation
		err   bool
	}{
		{"Test White Command", "white", painter.OperationFunc(painter.WhiteFill), false},
		{"Test Green Command", "green", painter.OperationFunc(painter.GreenFill), false},
		{"Test Update Command", "update", painter.UpdateOp, false},
		{"Test BgRect", "bgrect 0.1 0.1 0.5 0.5", &painter.BgRectOp{80, 80, 400, 400}, false},
		{"Test Figure Command", "figure 0.25 0.5", &painter.FigureOp{200, 400}, false},
		{"Test Move Command", "move 0.125 0.125", &painter.MoveOp{100, 100}, false},
		{"Test Reset Command", "reset", painter.OperationFunc(painter.Reset), false},
		{"Test Invalid Command", "gaysex", nil, true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			parser := &Parser{}
			ops, err := parser.Parse(strings.NewReader(test.input))
			if test.err {
				require.Error(t, err, "Error expected for input string %q", test.input)
				require.Nil(t, ops, "Nil expected as a return value for input string %q", test.input)
			} else {
				require.NoError(t, err, "Unexpected error for input string %q", test.input)
				require.Len(t, ops, 1, "Expected 1 operation for input string %q", test.input)
				switch test.op.(type) {
				case painter.OperationFunc:
					assert.IsType(t, test.op, ops[0], "Invalid type for input string %q", test.input)
				default:
					assert.Equal(t, test.op, ops[0], "Invalid type for input string %q", test.input)
				}
			}
		})
	}
}
