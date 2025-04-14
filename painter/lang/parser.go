package lang

import (
	"bufio"
	"github.com/roman-mazur/architecture-lab-3/painter"
	"io"
	"net"
	"strconv"
	"strings"
)

// Parser уміє прочитати дані з вхідного io.Reader та повернути список операцій представлені вхідним скриптом.
type Parser struct {
}

func (p *Parser) Parse(in io.Reader) ([]painter.Operation, error) {
	scanner := bufio.NewScanner(in)
	scanner.Split(bufio.ScanLines)
	var res []painter.Operation

	for scanner.Scan() {
		command := scanner.Text()
		comm, args, err := parseLine(command)
		if err != nil {
			return nil, err
		}
		op := mapToOp(comm, args)
		res = append(res, op)
	}

	return res, nil
}

func parseLine(line string) (comm string, args []int, err error) {
	fields := strings.Fields(line)
	if len(fields) == 0 {
		return "", nil, &net.ParseError{Text: "Line is empty"}
	}

	nums, err := Map(fields[1:], strconv.Atoi)
	if err != nil {
		return "", nil, &net.ParseError{Text: "Args are not integers"}
	}

	return fields[0], nums, nil
}

func mapToOp(comm string, args []int) painter.Operation {
	switch comm {
	case "white":
		return painter.OperationFunc(painter.WhiteFill)
	case "green":
		return painter.OperationFunc(painter.GreenFill)
	case "update":
		return painter.UpdateOp
	case "bgrect":
		return &painter.BgRectOp{
			X1: args[0],
			Y1: args[1],
			X2: args[2],
			Y2: args[3],
		}
	case "figure":
		return &painter.FigureOp{
			X: args[0],
			Y: args[1],
		}
	case "move":
		return &painter.MoveOp{
			X: args[0],
			Y: args[1],
		}
	case "reset":
		return painter.OperationFunc(painter.Reset)
	default:
		return nil
	}
}

func Map[T any, U any](in []T, f func(T) (U, error)) ([]U, error) {
	out := make([]U, len(in))
	for i, v := range in {
		u, err := f(v)
		if err != nil {
			return nil, err
		}
		out[i] = u
	}
	return out, nil
}
