package main

import (
	"context"
	"event_center/core"
	"log"
	"sync"
)

type BaseEvent struct {
	data interface{}
}

func (be BaseEvent) Name() core.EventName {
	return core.EventName("base_event")
}

func (be BaseEvent) Data() interface{} {
	return be.data
}

func Printf(handleName string, params interface{}) {
	log.Printf("%s:发送的消息是%s\n", handleName, params)
}

func main() {
	event := BaseEvent{}
	ctx, _ := context.WithCancel(context.Background())
	ec := core.CreateEventCenter(ctx)

	var wg sync.WaitGroup

	wg.Add(1)
	ec.Register(event.Name())
	ec.SetSubscribeHook(func(eventName core.EventName) {
		log.Printf("[%s] has been subscribe! \n", eventName)
	})

	ec.SetSendEventHook(func(event core.Event) {
		log.Printf("[%s] send data is 【%v】! \n", event.Name(), event.Data())
	})

	ec.Subscribe(event.Name(), "base_handle", func(param interface{}) {
		Printf("base_handle", param)
		wg.Done()
	})
	ec.SendEvent(BaseEvent{data: "hello"})
	wg.Wait()
}
