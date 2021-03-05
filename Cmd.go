package run

type cmd struct {
    what int
    fn   func()
}

func (j *cmd) Run() {
    if j.fn != nil {
        j.fn()
    }

}
