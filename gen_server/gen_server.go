package gen_server

import "time"

const Infinity time.Duration = 1<<63 - 1

type State struct {
	Ainf
	delta time.Duration
	t     *time.Timer
}

func (s *State) init(a Ainf) {
	s.t = time.NewTimer(Infinity)
	s.Ainf=a
}

type Ainf interface {
	Handle_msg(self, msg interface{}) time.Duration
	Handle_timeout(self interface{}) time.Duration
}

func Serve(channel chan interface{}, exstate Ainf) {
	state := new(State)
	state.init(exstate)
	go state.sleep(channel)
}

//exstate,channel : user define
func (state*State)sleep(channel chan interface{}) {

	var msg interface{}
	for {
		select {
		case msg = <-channel:
			state.delta = state.Ainf.Handle_msg(state.Ainf, msg)

		case <-state.t.C:
			state.delta = state.Ainf.Handle_timeout(state.Ainf)

		}
		state.t.Reset(time.Duration(abs(int64(state.delta))))

	}

}
func abs(n int64) int64 {
	y := n >> 63
	return (n ^ y) - y
}
