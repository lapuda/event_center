package core

import (
	"context"
	"errors"
	"log"
)

type EventName string
type HandlerName string

type EventHandler func(param interface{})

type EventCenter struct {
	Ctx             context.Context
	eventList       map[EventName]struct{}
	message         chan Event
	subscribeList   map[EventName]map[HandlerName]EventHandler
	subscribeHook   func(eventName EventName)
	unSubscribeHook func(eventName EventName)
	sendEventHook   func(event Event)
}

func (ev *EventCenter) Register(eventName EventName) {
	ev.eventList[eventName] = struct{}{}
	if _, ok := ev.subscribeList[eventName]; !ok {
		ev.subscribeList[eventName] = make(map[HandlerName]EventHandler)
	}
}

func (ev *EventCenter) UnRegister(event Event) {
	delete(ev.eventList, event.Name())
	delete(ev.subscribeList, event.Name())
}

func (ev *EventCenter) Subscribe(eventName EventName, handlerName HandlerName, handler EventHandler) {
	ev.subscribeList[eventName][handlerName] = handler
	ev.subscribeHook(eventName)
}

func (ev *EventCenter) UnSubscribe(eventName EventName, handlerName HandlerName) {
	delete(ev.subscribeList[eventName], handlerName)
	ev.unSubscribeHook(eventName)
}

func (ev *EventCenter) SendEvent(event Event) error {
	if _, ok := ev.eventList[event.Name()]; !ok {
		return errors.New("event " + string(event.Name()) + " not found!")
	}
	ev.message <- event
	return nil
}

func (ev *EventCenter) run() {
	for {
		select {
		case <-ev.Ctx.Done():
			log.Printf("force exit\n")
			return
		case event := <-ev.message:
			go ev.broadcast(event)
		}
	}
}

func (ev *EventCenter) broadcast(event Event) {
	ev.sendEventHook(event)
	if readySend, ok := ev.subscribeList[event.Name()]; ok {
		for _, handler := range readySend {
			go handler(event.Data())
		}
	}
	return
}

func (ev *EventCenter) SetSubscribeHook(handle func(eventName EventName)) {
	ev.subscribeHook = handle
}

func (ev *EventCenter) SetUnSubscribeHook(handle func(eventName EventName)) {
	ev.unSubscribeHook = handle
}

func (ev *EventCenter) SetSendEventHook(handle func(event Event)) {
	ev.sendEventHook = handle
}

func CreateEventCenter(ctx context.Context) *EventCenter {
	var ec *EventCenter
	ec = &EventCenter{
		Ctx:             ctx,
		eventList:       make(map[EventName]struct{}),
		message:         make(chan Event),
		subscribeList:   make(map[EventName]map[HandlerName]EventHandler),
		subscribeHook:   func(eventName EventName) {},
		unSubscribeHook: func(eventName EventName) {},
		sendEventHook:   func(event Event) {},
	}
	go ec.run()
	return ec
}
