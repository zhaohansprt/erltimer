package timer

import (
	"fmt"
	"sort"
	"test/erltimer/gen_server"
	"time"
)

type Timer struct {
	ertchan chan interface{}
	*exstate
}

//make sure to Start before you use NewTimer
func (et *Timer) Start(bufsize int) {
	et.ertchan = make(chan interface{}, bufsize)
	loop(et.ertchan, et.exstate)

}

//just like the time.NewTimer(du) does
func (et *Timer) NewTimer(du time.Duration) (ert ErlTimer) {
	if du <= 0 {
		ert.C = make(chan uint8, 1)
		close(ert.C)
		return
	}
	ert.C = make(chan uint8, 1)
	ert.ts = time.Now().UnixNano() + int64(du)
	et.ertchan <- &ert
	return
}

func (et *Timer) NewTimerTest(du time.Duration, trace string) (ert ErlTimer) {
	if du <= 0 {
		ert.C = make(chan uint8, 1)
		close(ert.C)
		return
	}
	ert.TrackMark = trace
	ert.C = make(chan uint8, 1)
	ert.ts = time.Now().UnixNano() + int64(du)

	et.ertchan <- &ert

	return
}

func (exstate) Handle_timeout(sinf interface{}) time.Duration {
	state := sinf.(*exstate)
	close(state.ert.C)
tag0:
	if state.ert = state.popmod(); state.ert != nil {
		state.delta = time.Duration(time.Now().UnixNano() - state.ert.ts) //处理 timer wheel 增量
		if state.delta >= 0 {
			close(state.ert.C)
			//ert=nil//goto tagx1
			goto tag0 //递回处理相同timer
			//t.Reset(0)
		} else {
			return -state.delta
		}
	} else {
		return gen_server.Infinity
	}
}

func (exstate) Handle_msg(sinf, msg interface{}) time.Duration {
	state := sinf.(*exstate)
	ert0 := msg.(*ErlTimer)
	switch {
	case state.ert == nil:
		state.handle0(ert0)

	case state.ert.ts < 0:
		panic("ts can't < 0")
	default:
		state.swap = state.ert.ts - ert0.ts
		if state.swap > 0 {
			state.handle0(ert0)
		} else {
			state.push(ert0) //入有序数组

		}

	}
	return state.delta
}
func (state *exstate) handle0(ert0 *ErlTimer) {
	state.delta = time.Duration(time.Now().UnixNano() - ert0.ts) //处理 timer wheel 增量
	if state.delta >= 0 {                                        //tag1:
		fmt.Println("delta guard failed ......!!!!!!!!!!!!!!!!!!!!!!!!!!! TrackMark: ", ert0.TrackMark)
		close(ert0.C)
	} else {
		state.ert = ert0 //重置ert
		state.push(ert0) //入有序数组
	}
}

////////////////////////////////internal uses//////////////////

type exstate struct {
	delta time.Duration
	swap  int64
	ert   *ErlTimer
	//m   sync.Mutex
	ets ETS
}

func (ckobj *exstate) remove(i int) {
	//ckobj.m.Lock()
	//defer ckobj.m.Unlock()
	if i == 0 {
		ckobj.ets = ckobj.ets[1:] //扔掉第一个
		return
	}
	ckobj.ets = append(ckobj.ets[:i], ckobj.ets[i+1:]...)

	return
}

//pop 修改版 先删除头部 再读取新头
func (ckobj *exstate) popmod() *ErlTimer {
	//ckobj.m.Lock()
	//defer ckobj.m.Unlock()
	ckobj.remove(0)
	//fmt.Println("测试数组长度打印popmod", len(ckobj.ets))
	if len(ckobj.ets) == 0 {
		return nil
	}
	return ckobj.ets[0]
}
func (ckobj *exstate) push(et *ErlTimer) {
	//ckobj.m.Lock()
	//defer ckobj.m.Unlock()
	ckobj.ets = append(ckobj.ets, et)
	sort.Sort(ckobj.ets)
	//fmt.Println("测试数组长度打印push", len(ckobj.ets))
	return
}

//不支持多开 原因是 多个timer 不能共享同一个ets
func loop(ck chan interface{}, s *exstate) {
	//for ; num > 0; num-- {

	//	cks = append(cks, channel)
	gen_server.Serve(ck, s)

	//}

}

type ErlTimer struct {
	C         chan uint8
	ts        int64
	TrackMark string
}

type ETS []*ErlTimer

func (s ETS) Len() int           { return len(s) }
func (s ETS) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s ETS) Less(i, j int) bool { return s[i].ts < s[j].ts }
