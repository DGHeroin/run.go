package run

import (
    "context"
)

type Dispatcher struct {
    context    context.Context
    cancel     context.CancelFunc
    workerPool chan *worker
    cmdQueue   chan *cmd
}

func NewDispatcher(maxWorkers int, backlogJob int) *Dispatcher {
    ctx, cancel := context.WithCancel(context.Background())
    d := &Dispatcher{
        workerPool: make(chan *worker, maxWorkers),
        cmdQueue:   make(chan *cmd, backlogJob),
        context:    ctx,
        cancel:     cancel,
    }
    go d.start(maxWorkers)
    return d
}

func (d *Dispatcher) start(maxWorkers int) {
    for i := 0; i < maxWorkers; i++ {
        w := newWorker(d.context, d.workerPool)
        go w.start()
    }
    d.dispatch()
}
func (d *Dispatcher) Run(cb func()) {
    job := &cmd{
        what: 0,
        fn: cb,
    }
    d.cmdQueue <- job
}
func (d *Dispatcher) dispatch() {
    for {
        select {
        case w := <-d.workerPool: // pick a worker
            c := <-d.cmdQueue // pick cmd from queue
            w.jobChan <- c    // dispatch the worker
        }
    }
}

func (d *Dispatcher) Stop() {
    close(d.workerPool)
    d.cancel()
}

func (d *Dispatcher) StopWorker(n int) {
    for i :=0;i<n;i++ {
        c := &cmd{
            what: 1,
        }
        d.cmdQueue <- c
    }
}
