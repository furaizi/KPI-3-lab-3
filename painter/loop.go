package painter

import (
	"image"
	"sync"

	"golang.org/x/exp/shiny/screen"
)

// Receiver отримує текстуру, яка була підготовлена в результаті виконання команд у циклі подій.
type Receiver interface {
	Update(t screen.Texture)
}

// Loop реалізує цикл подій для формування текстури отриманої через виконання операцій отриманих з внутрішньої черги.
type Loop struct {
	Receiver Receiver
	next screen.Texture
	prev screen.Texture
	stopReq bool
	stopped chan struct{}
	mq messageQueue
}

var size = image.Pt(400, 400)

// Start запускає цикл подій. Цей метод потрібно запустити до того, як викликати на ньому будь-які інші методи.
func (l *Loop) Start(s screen.Screen) {

	l.next, _ = s.NewTexture(size)
	l.prev, _ = s.NewTexture(size)
	l.mq = messageQueue{}

	go l.mainEventLoop() // Запуск обробника подій у окремій горутині
	
}

// Основний цикл обробки подій
func (l *Loop) mainEventLoop() {

	for {

		if op := l.mq.Pull(); op != nil {

			if update := op.Do(l.next); update {

				l.Receiver.Update(l.next)
				l.next, l.prev = l.prev, l.next

			}

		}

	}

}

// Post додає нову операцію у внутрішню чергу.
func (l *Loop) Post(op Operation) {

	if op != nil {

		l.mq.Push(op)

	}

}

// StopAndWait сигналізує про необхідність завершити цикл та блокується до моменту його повної зупинки.
func (l *Loop) StopAndWait() {

	l.Post(OperationFunc(func(screen.Texture) {

		l.stopReq = true

	}))

	<-l.stopped

}

// Реалізована черга подій
// Структура черги повідомлень
type messageQueue struct {
	mu      sync.Mutex // Мьютекс для захисту конкурентного доступу
	Queue   []Operation // Список операцій
	blocked chan struct{} // Канал, який блокує Pull, поки черга порожня
}

// Додавання операції в чергу
func (mq *messageQueue) Push(op Operation) {

	mq.mu.Lock()

	defer mq.mu.Unlock()

	mq.Queue = append(mq.Queue, op)

	if mq.blocked != nil {

		close(mq.blocked)
		mq.blocked = nil

	}

}

// Витягування операції з черги (якщо черга порожня, то блокує)
func (mq *messageQueue) Pull() Operation {

	mq.mu.Lock()

	defer mq.mu.Unlock()

	for len(mq.Queue) == 0 {

		mq.blocked = make(chan struct{})
		mq.mu.Unlock()
		<-mq.blocked
		mq.mu.Lock()

	}

	op := mq.Queue[0]
	mq.Queue[0] = nil
	mq.Queue = mq.Queue[1:]

	return op

}