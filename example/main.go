package example

import (
    "github.com/DGHeroin/run.go"
    "log"
    "time"
)

func main() {
    d := run.NewDispatcher(3, 1000)
    fn := func(i int) func() {
        return func() {
            log.Println(i)
            time.Sleep(time.Second)
        }
    }
    for i := 0; i < 5; i++ {
        d.Run(fn(i))
    }

    time.AfterFunc(time.Second*2, func() {
        d.StopWorker(1)
    })

    time.AfterFunc(time.Second*8, func() {
        d.Stop()
    })

    time.Sleep(time.Second * 15)
}
