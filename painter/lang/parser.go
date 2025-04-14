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
	// Наразі структура порожня, але може містити внутрішній стан, якщо знадобиться.
}

// Parse читає вхідний потік (наприклад, тіло HTTP запиту), розбиває його на рядки,
// виконує парсинг кожного рядка та повертає список painter.Operation.
// Якщо виникає помилка при парсингу якоїсь команди, повертається відповідна помилка.
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

// parseLine розбиває рядок на команду та аргументи.
// Повертає:
//   - comm: перше слово як назва команди;
//   - args: перетворені в int значення для решти полів;
//     перетворення здійснюється через виклик generic-функції Map із використанням floatStrToInt.
//
// Якщо рядок порожній або аргументи не є числами, повертається відповідна помилка.
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

// mapToOp приймає назву команди (comm) та числові аргументи (args)
// і повертає відповідну операцію (яка реалізує painter.Operation) для виклику.
// Якщо команда не підтримується, повертається помилка.
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

// Map — узагальнена функція, яка приймає слайс in типу T, застосовує до кожного елемента функцію f,
// що повертає значення типу U або помилку, і повертає слайс значень типу U або помилку.
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

// floatStrToInt перетворює рядок, що містить число з плаваючою точкою,
// у ціле число, масштабоване за допомогою ui.WINDOW_SIZE.
// Наприклад, якщо рядок містить "0.1" і ui.WINDOW_SIZE дорівнює 800,
// повернеться int(0.1 * 800) = 80.
func floatStrToInt(arg string) (int, error) {
	fl, err := strconv.ParseFloat(arg, 64)
	if err != nil {
		return 0, err
	}
	return int(fl * ui.WINDOW_SIZE), nil
}
