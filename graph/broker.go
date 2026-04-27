package graph

import (
	"sync"

	"github.com/jmontenegro/my-reminders-backend/graph/model"
)

type NotificationBroker struct {
	subscribers map[string]chan *model.Notification
	mutex       sync.Mutex
}

func NewBroker() *NotificationBroker {
	return &NotificationBroker{
		subscribers: make(map[string]chan *model.Notification),
	}
}

// Suscribirse a las notificaciones
func (b *NotificationBroker) Subscribe(userID string) chan *model.Notification {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	ch := make(chan *model.Notification, 1)
	b.subscribers[userID] = ch
	return ch
}

// Notificar a un usuario específico
func (b *NotificationBroker) Notify(userID string, n *model.Notification) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	if ch, ok := b.subscribers[userID]; ok {
		ch <- n
	}
}
