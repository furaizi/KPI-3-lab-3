package painter

import (
	"image"
	"image/color"
	"image/draw"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/exp/shiny/screen"
)

type Mock struct {
	mock.Mock
}

func (_ *Mock) NewBuffer(size image.Point) (screen.Buffer, error) {
	return nil, nil
}

func (_ *Mock) NewWindow(opts *screen.NewWindowOptions) (screen.Window, error) {
	return nil, nil
}

func (m *Mock) Update(texture screen.Texture) {
	m.Called(texture)
}

func (m *Mock) NewTexture(size image.Point) (screen.Texture, error) {
	args := m.Called(size)
	return args.Get(0).(screen.Texture), args.Error(1)
}

func (m *Mock) Release() {
	m.Called()
}

func (m *Mock) Upload(dp image.Point, src screen.Buffer, sr image.Rectangle) {
	m.Called(dp, src, sr)
}

func (m *Mock) Bounds() image.Rectangle {
	args := m.Called()
	return args.Get(0).(image.Rectangle)
}

func (m *Mock) Fill(dr image.Rectangle, src color.Color, op draw.Op) {
	m.Called(dr, src, op)
}

func (m *Mock) Size() image.Point {
	args := m.Called()
	return args.Get(0).(image.Point)
}

func (m *Mock) Do(t screen.Texture) bool {
	args := m.Called(t)
	return args.Bool(0)
}

// Допоміжні функції

// newLoopWithMocks створює підготовлений тестовий цикл Loop з усіма необхідними моками.
// Повертає: цикл рендерингу, мок текстури, мок приймача (Receiver) та мок екрану (screen)
func newLoopWithMocks(t *testing.T) (*Loop, *Mock, *Mock, *Mock) {
	textureSize := image.Pt(400, 400)

	textureMock := new(Mock)
	receiverMock := new(Mock)
	screenMock := new(Mock)

	screenMock.On("NewTexture", textureSize).Return(textureMock, nil)
	receiverMock.On("Update", textureMock).Return()

	loop := &Loop{Receiver: receiverMock}
	loop.Start(screenMock)

	textureMock.On("Bounds").Return(image.Rectangle{})

	return loop, textureMock, receiverMock, screenMock
}

// waitForProcessing робить коротку паузу, щоб дати час циклe Loop обробити чергу.
func waitForProcessing() {
	time.Sleep(100 * time.Millisecond)
}

// Тести

func TestPostSingleSuccess(t *testing.T) {
	loop, textureMock, receiverMock, screenMock := newLoopWithMocks(t)

	op := new(Mock)
	op.On("Do", textureMock).Return(true)

	loop.Post(op)
	waitForProcessing()

	op.AssertCalled(t, "Do", textureMock)
	receiverMock.AssertCalled(t, "Update", textureMock)
	screenMock.AssertCalled(t, "NewTexture", image.Pt(400, 400))
	assert.Empty(t, loop.mq.Queue)
}

func TestPostMultipleSuccess(t *testing.T) {
	loop, textureMock, receiverMock, screenMock := newLoopWithMocks(t)

	op1 := new(Mock)
	op2 := new(Mock)
	op1.On("Do", textureMock).Return(true)
	op2.On("Do", textureMock).Return(true)

	loop.Post(op1)
	loop.Post(op2)
	waitForProcessing()

	op1.AssertCalled(t, "Do", textureMock)
	op2.AssertCalled(t, "Do", textureMock)
	receiverMock.AssertCalled(t, "Update", textureMock)
	screenMock.AssertCalled(t, "NewTexture", image.Pt(400, 400))
	assert.Empty(t, loop.mq.Queue)
}

func TestPostSingleFailure(t *testing.T) {
	loop, textureMock, receiverMock, screenMock := newLoopWithMocks(t)

	op := new(Mock)
	op.On("Do", textureMock).Return(false)

	loop.Post(op)
	waitForProcessing()

	op.AssertCalled(t, "Do", textureMock)
	receiverMock.AssertNotCalled(t, "Update", textureMock)
	screenMock.AssertCalled(t, "NewTexture", image.Pt(400, 400))
	assert.Empty(t, loop.mq.Queue)
}
