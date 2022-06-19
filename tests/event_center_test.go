package tests

import (
	"context"
	"event_center/core"
	"log"
	"testing"
)

type TestEvent struct {
	data interface{}
}

func (e TestEvent) Data() interface{} {
	return e.data
}

func (e TestEvent) Name() core.EventName {
	return core.EventName("TestEvent")
}

func TestEventCenter(t *testing.T) {
	event := TestEvent{"testing......"}
	ctx, cancel := context.WithCancel(context.Background())
	ec := core.CreateEventCenter(ctx)

	ec.Register(event.Name())
	ec.Subscribe(event.Name(), "test_subscribe1", func(param interface{}) {
		log.Panicf("test_subscribe1 testting %s \n", param.(string))
	})

	//ec.Subscribe(event.Name(), "test_subscribe2", func(param interface{}) {
	//	log.Panicf("test_subscribe2 testting %s \n", param.(string))
	//})
	//
	//ec.Subscribe(event.Name(), "test_subscribe3", func(param interface{}) {
	//	log.Panicf("test_subscribe3 testting %s \n", param.(string))
	//})

	ec.SendEvent(event)
	cancel()
	//time.Sleep(2 * time.Second)
}
