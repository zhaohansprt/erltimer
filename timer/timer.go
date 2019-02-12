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

func Stats() int{
	fmt.Printf("	timer wheel ETS length:%v \n", cks[0].ets.Len())
	return cks[0].ets.Len()
}

func emit(timer *ErlTimer) {
	cks[rand.Intn(len(cks))].channel <- timer
}

//预留多开  当前不支持
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

func NewTimertest(du time.Duration, trace string) (ert ErlTimer) {
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

func sleep(tch *chanlocker, _ time.Duration) {
	var ert, ert0 *ErlTimer
	t := time.NewTimer(maxDuration)

	var delta, swap int64
	for {
		select {
		case ert0 = <-tch.channel:
			switch {
			case ert == nil:
				delta = time.Now().UnixNano() - ert0.ts //处理 timer wheel 增量
				if delta >= 0 {                         //tag1:
					fmt.Println("delta guard failed ......!!!!!!!!!!!!!!!!!!!!!!!!!!! TrackMark: ", ert0.TrackMark)
					close(ert0.C)
				} else {
					ert = ert0     //重置ert
					tch.push(ert0) //入有序数组
					t.Reset(time.Duration(-delta))
				}
			case ert.ts < 0:
				panic("ts can't < 0")
			default:
				swap = ert.ts - ert0.ts
				if swap > 0 {
					delta = time.Now().UnixNano() - ert0.ts //处理 timer wheel 增量
					if delta >= 0 {
						fmt.Println("delta guard failed ......!!!!!!!!!!!!!!!!!!!!!!!!!!! TrackMark: ", ert0.TrackMark)
						close(ert0.C)
					} else {
						//goto tag1
						ert = ert0     //重置ert
						tch.push(ert0) //入有序数组
						t.Reset(time.Duration(-delta))
					}
				} else {
					tch.push(ert0) //入有序数组

				}

			}
		case <-t.C:
			close(ert.C)
		tag0:
			if ert = tch.popmod(); ert != nil {
				delta = time.Now().UnixNano() - ert.ts //处理 timer wheel 增量
				if delta >= 0 {
					close(ert.C)
					//ert=nil//goto tagx1
					goto tag0 //递回处理相同timer
					//t.Reset(0)
				} else {
					t.Reset(time.Duration(-delta))
				}
			} else {
				t.Reset(maxDuration)
			}

		}
	}

}
