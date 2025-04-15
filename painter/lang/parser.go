package lang

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/roman-mazur/architecture-lab-3/painter"
	"github.com/roman-mazur/architecture-lab-3/ui"
)

// Parser уміє прочитати дані з вхідного io.Reader та повернути список операцій представлені вхідним скриптом.
type Parser struct {

	figures     []*painter.FigureOp
	moveOps     []painter.Operation
	lastBgColor painter.Operation
	lastBgRect  *painter.BgRectOp
	updateOp    painter.Operation

}

// initializeParserState ініціалізує початковий стан парсера
func (p *Parser) initializeParserState() {

	if p.lastBgColor == nil {

		p.lastBgColor = painter.OperationFunc(painter.Reset)

	}

	if p.updateOp != nil {

		p.updateOp = nil

	}
}

// Parse читає вхідний потік (наприклад, тіло HTTP запиту), розбиває його на рядки,
// виконує парсинг кожного рядка та повертає список painter.Operation.
// Якщо виникає помилка при парсингу якоїсь команди, повертається відповідна помилка.
func (p *Parser) Parse(in io.Reader) ([]painter.Operation, error) {

	p.initializeParserState()
	scanner := bufio.NewScanner(in)
	scanner.Split(bufio.ScanLines)
	
	for scanner.Scan() {

		command := scanner.Text()
		err := p.parse(command)

		if err != nil {

			return nil, err

		}

	}
	
	return p.finalParseResult(), nil

}

// finalParseResult формує кінцевий список операцій на основі поточного стану парсера
func (p *Parser) finalParseResult() []painter.Operation {

	var res []painter.Operation
	
	if p.lastBgColor != nil {

		res = append(res, p.lastBgColor)

	}
	
	if p.lastBgRect != nil {

		res = append(res, p.lastBgRect)

	}
	
	if len(p.moveOps) != 0 {

		res = append(res, p.moveOps...)

	}
	p.moveOps = nil
	
	if len(p.figures) != 0 {

		for _, figure := range p.figures {

			res = append(res, figure)
		}
	}
	
	if p.updateOp != nil {

		res = append(res, p.updateOp)

	}
	
	return res

}

// resetParserState скидає стан парсера
func (p *Parser) resetParserState() {

	p.figures = nil
	p.moveOps = nil
	p.lastBgColor = nil
	p.lastBgRect = nil
	p.updateOp = nil

}

// parse обробляє окремий рядок команди
func (p *Parser) parse(commandLine string) error {

	fields := strings.Fields(commandLine)

	if len(fields) == 0 {

		return errors.New("line is empty")

	}
	
	comm := fields[0]
	args, err := Map(fields[1:], floatStrToInt)

	if err != nil && len(fields) > 1 {

		return errors.New("args are not integers")
		
	}
	
	switch comm {
	case "white":
		p.lastBgColor = painter.OperationFunc(painter.WhiteFill)
	case "green":
		p.lastBgColor = painter.OperationFunc(painter.GreenFill)
	case "update":
		p.updateOp = painter.UpdateOp
	case "bgrect":
		if len(args) < 4 {
			return errors.New("not enough arguments for bgrect")
		}
		p.lastBgRect = &painter.BgRectOp{
			X1: args[0],
			Y1: args[1],
			X2: args[2],
			Y2: args[3],
		}
	case "figure":
		if len(args) < 2 {
			return errors.New("not enough arguments for figure")
		}
		figure := &painter.FigureOp{
			X: args[0],
			Y: args[1],
		}
		p.figures = append(p.figures, figure)
	case "move":
		if len(args) < 2 {
			return errors.New("not enough arguments for move")
		}
		moveOp := &painter.MoveOp{
			X:       args[0],
			Y:       args[1],
			Figures: p.figures,
		}
		p.moveOps = append(p.moveOps, moveOp)
	case "reset":
		p.resetParserState()
		p.lastBgColor = painter.OperationFunc(painter.Reset)
	default:
		return fmt.Errorf("unknown command: %s", comm)
	}
	
	return nil
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