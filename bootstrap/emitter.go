package bootstrap

import (
	"fmt"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/maogejing/emitter/domain"
)

type Emitter[TaskType interface{}] struct {
	ActiveProducerCount int32
	ActiveConsumerCount int32
	taskChan            chan TaskType
	stopChan            chan struct{}
	reduceConsumerChan  chan struct{}
	taskProducerFactory domain.TaskProducerFactory[TaskType]
	taskConsumerFactory domain.TaskConsumerFactory[TaskType]
	rootWg              *sync.WaitGroup // rootwg
	producerWg          *sync.WaitGroup // 生产者wg
	consumerWg          *sync.WaitGroup // 消费者wg
}

func NewEmitter[TaskType interface{}](
	taskProducerFactory domain.TaskProducerFactory[TaskType],
	taskConsumerFactory domain.TaskConsumerFactory[TaskType],
	rootWg *sync.WaitGroup,
) *Emitter[TaskType] {
	return &Emitter[TaskType]{
		ActiveProducerCount: 0,
		ActiveConsumerCount: 0,
		taskConsumerFactory: taskConsumerFactory,
		taskProducerFactory: taskProducerFactory,
		rootWg:              rootWg,
		taskChan:            make(chan TaskType, 1),
		reduceConsumerChan:  make(chan struct{}, 1),
		stopChan:            make(chan struct{}, 1),
		producerWg:          &sync.WaitGroup{},
		consumerWg:          &sync.WaitGroup{},
	}
}

func (e *Emitter[TaskType]) StartWorkingLoop() {
	defer e.rootWg.Done()
	e.producerWg.Add(1)
	go e.startProducer("DefaultProducer")

	<-e.stopChan
	e.producerWg.Wait()
	e.consumerWg.Wait()
	close(e.taskChan)
	close(e.reduceConsumerChan)
}

func (e *Emitter[TaskType]) StopWorkingLoop() {
	close(e.stopChan)
	for restMsg := range e.taskChan {
		fmt.Println(restMsg)
	}
}

func (e *Emitter[TaskType]) IncrConsumer() {
	e.consumerWg.Add(1)
	atomic.AddInt32(&e.ActiveConsumerCount, 1)
	go e.startConsumer("IncrConsumer" + strconv.Itoa(int(e.ActiveConsumerCount)))
}

func (e *Emitter[TaskType]) DecrConsumer() {
	if atomic.AddInt32(&e.ActiveConsumerCount, -1) < 0 {
		atomic.AddInt32(&e.ActiveConsumerCount, 1)
	} else {
		e.reduceConsumerChan <- struct{}{}
	}
}

func (e *Emitter[TaskType]) IncrProducer() {

}

func (e *Emitter[TaskType]) DecrProducer() {

}

func (e *Emitter[TaskType]) startProducer(name string) {
	defer e.producerWg.Done()
	taskGenFunc := e.taskProducerFactory.Generate()
ProducerLoop:
	for {
		select {
		case _, ok := <-e.stopChan:
			if !ok {
				break ProducerLoop
			}
		default:
			time.Sleep(300 * time.Millisecond)
			task := taskGenFunc()
			e.taskChan <- task
		}
	}
	fmt.Println("producer", name, "exit")
}

func (e *Emitter[TaskType]) startConsumer(name string) {
	defer e.consumerWg.Done()
	taskConsumerFunc := e.taskConsumerFactory.Generate()
ConsumerLoop:
	for {
		select {
		case _, ok := <-e.stopChan:
			if !ok {
				break ConsumerLoop
			}
		case <-e.reduceConsumerChan:
			break ConsumerLoop
		case task := <-e.taskChan:
			// 执行消费动作
			taskConsumerFunc(task)
		}
	}
	fmt.Println("consumer", name, "exit")
}
