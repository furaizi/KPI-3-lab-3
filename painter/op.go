package painter

import (
	"image"
	"image/color"
	"image/draw"

	"golang.org/x/exp/shiny/screen"
)

// Operation змінює вхідну текстуру.
type Operation interface {
	// Do виконує зміну операції, повертаючи true, якщо текстура вважається готовою для відображення.
	Do(t screen.Texture) (ready bool)
}

// OperationList групує список операції в одну.
type OperationList []Operation

func (ol OperationList) Do(t screen.Texture) (ready bool) {
	for _, o := range ol {
		ready = o.Do(t) || ready
	}
	return
}

// UpdateOp операція, яка не змінює текстуру, але сигналізує, що текстуру потрібно розглядати як готову.
var UpdateOp = updateOp{}

type updateOp struct{}

func (op updateOp) Do(t screen.Texture) bool { return true }

// OperationFunc використовується для перетворення функції оновлення текстури в Operation.
type OperationFunc func(t screen.Texture)

func (f OperationFunc) Do(t screen.Texture) bool {
	f(t)
	return false
}

// WhiteFill зафарбовує тестуру у білий колір. Може бути викоистана як Operation через OperationFunc(WhiteFill).
func WhiteFill(t screen.Texture) {
	t.Fill(t.Bounds(), color.White, screen.Src)
}

// GreenFill зафарбовує тестуру у зелений колір. Може бути викоистана як Operation через OperationFunc(GreenFill).
func GreenFill(t screen.Texture) {
	t.Fill(t.Bounds(), color.RGBA{G: 0xff, A: 0xff}, screen.Src)
}

type BgRectOp struct {
	X1 int
	Y1 int
	X2 int
	Y2 int
}

type FigureOp struct {
	X int
	Y int
}

type MoveOp struct {
	X int
	Y int
	Figures []*FigureOp
}

// Do виконує операцію на об'єкті BgRectOp, заповнюючи прямокутник чорним кольором на переданому текстурному об'єкті
func (op *BgRectOp) Do(t screen.Texture) bool {

	t.Fill(image.Rect(op.X1, op.Y1, op.X2, op.Y2), color.RGBA{0, 0, 0, 255}, screen.Src)

	return false

}

// Do виконує операцію на об'єкті FigureOp, малюючи T-образну фігуру на текстурі.
func (op *FigureOp) Do(t screen.Texture) bool {

	t.Fill(image.Rect(op.X-60, op.Y-40, op.X+60, op.Y), color.RGBA{0, 0, 255, 255}, draw.Src)
	t.Fill(image.Rect(op.X-20, op.Y, op.X+20, op.Y+40), color.RGBA{0, 0, 255, 255}, draw.Src)

	return false

}

// Do виконує операцію переміщення всіх фігур FigureOp на екран з вказаними зміщеннями по осях X та Y.
func (op *MoveOp) Do(t screen.Texture) bool {

	for i := range op.Figures {

		op.Figures[i].Y += op.Y
		op.Figures[i].X += op.X

	}

	return false
	
}

// Reset скидає текстуру до її початкового стану, заповнюючи всю область текстури чорним кольором.
func Reset(t screen.Texture) {

	t.Fill(t.Bounds(), color.RGBA{0, 0, 0, 255}, draw.Src)

}
