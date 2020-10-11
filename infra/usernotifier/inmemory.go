package usernotifier

import (
	"go-users-example/domain/users"
)

// Based on the way the chan are architectured, if the listener don't listen to the chan, it will stop the app from
// working at some point. should be fixed for production use.
const defaultChanSize = 500

// InMemory is a basic implementation of an in memory notifier infra as kafka, nats, etc.
type InMemory struct {
	reader []chan *users.ChangeEvent
}

// NewInMemory will instantiate properly an InMemory
func NewInMemory() *InMemory {
	return &InMemory{}
}

// Notify will send a notification of user change in the systems. implements users.ChangeNotifier
func (i *InMemory) Notify(event *users.ChangeEvent) error {
	for _, c := range i.reader {
		c <- event
	}
	return nil
}

// Listen will generate a new subcription to the ChangeEvent notification
func (i *InMemory) Listen() chan *users.ChangeEvent {
	c := make(chan *users.ChangeEvent, defaultChanSize)
	i.reader = append(i.reader, c)
	return c
}
