package main

type Dispatcher struct {
	WorkerPool chan chan Job
	maxWorkers int
}

func NewDispatcher() *Dispatcher {
	workersCount := 10
	pool := make(chan chan Job, workersCount)
	return &Dispatcher{
		WorkerPool: pool,
		maxWorkers: workersCount,
	}
}

func (d *Dispatcher) Run() {
	for i := 0; i < d.maxWorkers; i++ {
		worker := NewWorker(d.WorkerPool)
		worker.Start()
	}
	go d.dispatch()
}

func (d *Dispatcher) dispatch() {
	for {
		select {
		case job := <-JobQueue:
			go func(job Job) {
				jobChannel := <-d.WorkerPool
				jobChannel <- job
			}(job)
		}
	}
}
