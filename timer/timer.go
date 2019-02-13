package timer

import (
	"fmt"
	"math/rand"
	"sort"
	"time"
)

var cks = []*chanlocker{}

type chanlocker struct {
	channel chan *ErlTimer
	ets     ETS
}

func (ckobj *chanlocker) remove(i int) {
	if i == 0 {
		ckobj.ets = ckobj.ets[1:] //扔掉第一个
		return
	}
	ckobj.ets = append(ckobj.ets[:i], ckobj.ets[i+1:]...)

	return
}

//pop 修改版 先删除头部 再读取新头
func (ckobj *chanlocker) popmod() *ErlTimer {
	ckobj.remove(0)
	//fmt.Println("测试数组长度打印popmod", len(ckobj.ets))
	if len(ckobj.ets) == 0 {
		return nil
	}
	return ckobj.ets[0]
}
func (ckobj *chanlocker) push(et *ErlTimer) {
	ckobj.ets = append(ckobj.ets, et)
	sort.Sort(ckobj.ets)
	//fmt.Println("测试数组长度打印push", len(ckobj.ets))
	return
}

func initchlocker() (ch chanlocker) {

	ch.channel = make(chan *ErlTimer, 128)
	ch.ets = make(ETS, 0, 256)
	return
}

func Start(i int) {
	loop(i, time.Nanosecond)
}

func Stats() int {
	fmt.Printf("	timer wheel ETS length:%v \n", cks[0].ets.Len())
	return cks[0].ets.Len()
}

func emit(timer *ErlTimer) {
	cks[rand.Intn(len(cks))].channel <- timer
}

//多开
func loop(num int, duration time.Duration) {
	for ; num > 0; num-- {
		obj := initchlocker()
		cks = append(cks, &obj)
		go sleep(&obj, duration)

	}

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

func NewTimerTest(du time.Duration, trace string) (ert ErlTimer) {
	if du <= 0 {
		ert.C = make(chan uint8, 1)
		close(ert.C)
		return
	}
	ert.TrackMark = trace
	ert.C = make(chan uint8, 1)
	ert.ts = time.Now().UnixNano() + int64(du)

	emit(&ert)

	return
}

func NewTimer(du time.Duration) (ert ErlTimer) {
	if du <= 0 {
		ert.C = make(chan uint8, 1)
		close(ert.C)
		return
	}
	ert.C = make(chan uint8, 1)
	ert.ts = time.Now().UnixNano() + int64(du)

	emit(&ert)

	return
}

const maxDuration time.Duration = 1<<63 - 1

type State struct {
	delta,
	swap int64
	ert,
	ert0 *ErlTimer
	t   *time.Timer
	tch *chanlocker
}

func (s *State) init(tch0 *chanlocker) {
	s.t = time.NewTimer(maxDuration)
	s.tch = tch0
}

func sleep(tch0 *chanlocker, _ time.Duration) {
	state := new(State)
	state.init(tch0)

	for {
		select {
		case state.ert0 = <-state.tch.channel:
			switch {
			case state.ert == nil:
				handle0(state)
			case state.ert.ts < 0:
				panic("ts can't < 0")
			default:
				state.swap = state.ert.ts - state.ert0.ts
				if state.swap > 0 {
					handle0(state)
				} else {
					state.tch.push(state.ert0) //入有序数组

				}

			}
		case <-state.t.C:
			close(state.ert.C)
		tag0:
			if state.ert = state.tch.popmod(); state.ert != nil {
				state.delta = time.Now().UnixNano() - state.ert.ts //处理 timer wheel 增量
				if state.delta >= 0 {
					close(state.ert.C)
					//ert=nil//goto tagx1
					goto tag0 //递回处理相同timer
					//t.Reset(0)
				} else {
					state.t.Reset(time.Duration(-state.delta))
				}
			} else {
				state.t.Reset(maxDuration)
			}

		}
	}

}

func handle0(state *State) {
	state.delta = time.Now().UnixNano() - state.ert0.ts //处理 timer wheel 增量
	if state.delta >= 0 {                               //tag1:
		fmt.Println("delta guard failed ......!!!!!!!!!!!!!!!!!!!!!!!!!!! TrackMark: ", state.ert0.TrackMark)
		close(state.ert0.C)
	} else {
		state.ert = state.ert0     //重置ert
		state.tch.push(state.ert0) //入有序数组
		state.t.Reset(time.Duration(-state.delta))
	}
}
