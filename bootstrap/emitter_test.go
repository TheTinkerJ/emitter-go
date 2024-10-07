package bootstrap_test

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/maogejing/emitter/bootstrap"
	"github.com/maogejing/emitter/domain"
)

type TestTaskProducerFactory struct {
	taskNoIndicator int32
	domain.TaskProducerFactory[int]
}

func (tp *TestTaskProducerFactory) Generate() func() int {
	return func() int {
		atomic.AddInt32(&tp.taskNoIndicator, 1)
		return int(tp.taskNoIndicator)
	}
}

type TestTaskConsumerFactory struct {
	domain.TaskConsumerFactory[int]
	t *testing.T
}

func (tc *TestTaskConsumerFactory) Generate() func(int) {
	return func(i int) {
		tc.t.Log(i)
	}
}

func TestEmitterBasicRun(t *testing.T) {
	rootwg := sync.WaitGroup{}
	ttp := &TestTaskProducerFactory{taskNoIndicator: 0}
	ttc := &TestTaskConsumerFactory{t: t}
	emitter := bootstrap.NewEmitter[int](ttp, ttc, &rootwg)

	t.Log("start")
	rootwg.Add(1)
	go emitter.StartWorkingLoop()
	time.Sleep(1 * time.Second)
	emitter.IncrConsumer()
	time.Sleep(2 * time.Second)
	emitter.DecrConsumer()
	// t.Log("DecrConsumer")
	time.Sleep(2 * time.Second)
	t.Log("StopWorkingLoop")
	emitter.StopWorkingLoop()
	time.Sleep(1 * time.Second)
	rootwg.Wait()
}
