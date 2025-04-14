package ui

import (
	"image"
	"image/color"
	"log"

	"golang.org/x/exp/shiny/driver"
	"golang.org/x/exp/shiny/imageutil"
	"golang.org/x/exp/shiny/screen"
	"golang.org/x/image/draw"
	"golang.org/x/mobile/event/key"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/mouse"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/size"
)

// Змінні для координат фігури
var figureX, figureY int = 400, 400 // Початкові координати фігури (центральні)

type Visualizer struct {
	Title         string
	Debug         bool
	OnScreenReady func(s screen.Screen)

	w    screen.Window
	tx   chan screen.Texture
	done chan struct{}

	sz  size.Event
	pos image.Rectangle
}

func (pw *Visualizer) Main() {
	pw.tx = make(chan screen.Texture)
	pw.done = make(chan struct{})
	pw.pos.Max.X = 200
	pw.pos.Max.Y = 200
	driver.Main(pw.run)
}

func (pw *Visualizer) Update(t screen.Texture) {
	pw.tx <- t
}

func (pw *Visualizer) run(s screen.Screen) {
	w, err := s.NewWindow(&screen.NewWindowOptions{
		Title: pw.Title,
		Width: 800, // Вікно розміру 800x800 px
		Height: 800,
	})
	if err != nil {
		log.Fatal("Failed to initialize the app window:", err)
	}
	defer func() {
		w.Release()
		close(pw.done)
	}()

	if pw.OnScreenReady != nil {
		pw.OnScreenReady(s)
	}

	pw.w = w

	events := make(chan any)
	go func() {
		for {
			e := w.NextEvent()
			if pw.Debug {
				log.Printf("new event: %v", e)
			}
			if detectTerminate(e) {
				close(events)
				break
			}
			events <- e
		}
	}()

	var t screen.Texture

	for {
		select {
		case e, ok := <-events:
			if !ok {
				return
			}
			pw.handleEvent(e, t)

		case t = <-pw.tx:
			w.Send(paint.Event{})
		}
	}
}

func detectTerminate(e any) bool {
	switch e := e.(type) {
	case lifecycle.Event:
		if e.To == lifecycle.StageDead {
			return true // Window destroy initiated.
		}
	case key.Event:
		if e.Code == key.CodeEscape {
			return true // Esc pressed.
		}
	}
	return false
}

func (pw *Visualizer) handleEvent(e any, t screen.Texture) {
	switch e := e.(type) {
	case size.Event: // Оновлення розміру вікна.
		pw.sz = e

	case error:
		log.Printf("ERROR: %s", e)

	case mouse.Event:
		if e.Button == mouse.ButtonRight && e.Direction == mouse.DirPress {
			// Переміщення фігури до нової позиції миші
			figureX, figureY = int(e.X), int(e.Y)
			pw.w.Send(paint.Event{}) // Оновлення малюнку
		}

	case paint.Event:
		// Малювання контенту.
		if t == nil {
			pw.drawDefaultUI()
		} else {
			pw.w.Scale(pw.sz.Bounds(), t, t.Bounds(), draw.Src, nil)
		}
		pw.w.Publish()
	}
}

func (pw *Visualizer) drawDefaultUI() {
	// Зелений фон
	pw.w.Fill(pw.sz.Bounds(), color.RGBA{G: 255, A: 255}, draw.Src)

	// Білий контур
	for _, br := range imageutil.Border(pw.sz.Bounds(), 10) {
		pw.w.Fill(br, color.White, draw.Src)
	}

	// Фігура Т-0 синього кольору
	blue := color.RGBA{B: 255, A: 255}
	cellSize := 40

	// Початкові координати фігури повинні бути по центру вікна
	startX := figureX - cellSize*3/2 // Центруємо по горизонталі (для фігури "Т")
	startY := figureY - cellSize*2   // Центруємо по вертикалі (для фігури "Т")

	// Верхній ряд фігури (3 квадрати)
	for i := 0; i < 3; i++ {
		rect := image.Rect(startX+i*cellSize, startY, startX+(i+1)*cellSize, startY+cellSize)
		pw.w.Fill(rect, blue, draw.Src)
	}

	// Нижній блок (під центральним верхнім)
	centerX := startX + cellSize
	rect := image.Rect(centerX, startY+cellSize, centerX+cellSize, startY+2*cellSize)
	pw.w.Fill(rect, blue, draw.Src)
}
