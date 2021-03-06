package event_queue

import (
	"github.com/sniperHW/kendynet/util"
	"sync/atomic"
	"fmt"
)

type EventQueue struct {
	eventQueue *util.BlockQueue
	started     int32
}

func New() *EventQueue {
	r := &EventQueue{}
	r.eventQueue = util.NewBlockQueue()
	return r
}

func (this *EventQueue) PostEvent(ev interface{}) error {
	return this.eventQueue.Add(ev)
}

func (this *EventQueue) Close() {
	this.eventQueue.Close()
}

func (this *EventQueue) Start(onEvent func(interface{})) error {
	
	if nil == onEvent {
		return fmt.Errorf("onEvent == nil")
	}

	if !atomic.CompareAndSwapInt32(&this.started,0,1) {
        return fmt.Errorf("already started")
    }

	for {
		closed,localList := this.eventQueue.Get()
		if closed {
			return nil
		}
		size := len(localList)
		for i := 0; i < size; i++ {
			onEvent(localList[i])
		}
	}
}
