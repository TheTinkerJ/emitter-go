package domain

type TaskProducerFactory[TaskType interface{}] interface {
	Generate() func() TaskType
}

type TaskConsumerFactory[TaskType interface{}] interface {
	Generate() func(TaskType)
}
