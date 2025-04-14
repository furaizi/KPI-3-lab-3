package lang

import (
	"bufio"
	"errors"
	"github.com/roman-mazur/architecture-lab-3/painter"
	"github.com/roman-mazur/architecture-lab-3/ui"
	"io"
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
		op, err := mapToOp(comm, args)
		if err != nil {
			return nil, err
		}
		res = append(res, op)
	}

	return res, nil
}

func parseLine(line string) (comm string, args []int, err error) {
	fields := strings.Fields(line)
	if len(fields) == 0 {
		return "", nil, errors.New("line is empty")
	}

	nums, err := Map(fields[1:], floatStrToInt)
	if err != nil {
		return "", nil, errors.New("args are not integers")
	}
	return fields[0], nums, nil
}

func mapToOp(comm string, args []int) (painter.Operation, error) {
	switch comm {
	case "white":
		return painter.OperationFunc(painter.WhiteFill), nil
	case "green":
		return painter.OperationFunc(painter.GreenFill), nil
	case "update":
		return painter.UpdateOp, nil
	case "bgrect":
		return &painter.BgRectOp{
			X1: args[0],
			Y1: args[1],
			X2: args[2],
			Y2: args[3],
		}, nil
	case "figure":
		return &painter.FigureOp{
			X: args[0],
			Y: args[1],
		}, nil
	case "move":
		return &painter.MoveOp{
			X: args[0],
			Y: args[1],
		}, nil
	case "reset":
		return painter.OperationFunc(painter.Reset), nil
	default:
		return nil, errors.New("unknown command: " + comm)
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

func floatStrToInt(arg string) (int, error) {
	fl, err := strconv.ParseFloat(arg, 64)
	if err != nil {
		return 0, err
	}
	return int(fl * ui.WINDOW_SIZE), nil
}
