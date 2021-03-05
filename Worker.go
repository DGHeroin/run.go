package run

import (
    "context"
    "sync/atomic"
)

type worker struct {
    context context.Context
    pool    chan *worker
    jobChan chan *cmd
    quit    chan bool
    working int32
}

func newWorker(ctx context.Context, pool chan *worker) *worker {
    return &worker{
        context: ctx,
        pool:    pool,
        jobChan: make(chan *cmd),
        quit:    make(chan bool),
    }
}

func (w *worker) start() {
    defer func() {
        recover()
    }()
    for {
        // registering to workers pool
        w.pool <- w
        select {
        case <-w.context.Done():
            return
        case c := <-w.jobChan:
            if c == nil {
                continue
            }
            switch c.what {
            case 0:
                atomic.StoreInt32(&w.working, 1)
                c.Run()
                atomic.StoreInt32(&w.working, 0)
            case 1:
                w.stop()
            }
        case <-w.quit:
            return
        }
    }
}

func (w *worker) stop() {
    go func() {
        close(w.jobChan)
        w.quit <- true
    }()
}
