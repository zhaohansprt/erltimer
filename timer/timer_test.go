package timer

import (
	"container/heap"
	"fmt"
	"testing"
	"time"
)

var st *Timer

func init() {
	st = new(Timer)
	st.exstate = new(exstate)
	st.exstate.ets = make(ETS, 0, 256)
	heap.Init(&st.exstate.ets)
	st.Start(256)
}

//to see the wheel queue
func Stats() int {
	fmt.Printf("	timer wheel ETS length:%v \n", st.exstate.ets.Len())
	return st.exstate.ets.Len()
}
func Test_NewTimer(t *testing.T) {
	go testwraper(time.Minute+8*time.Second, fmt.Sprintf("minute test1"))
	//
	go testwraper(62*time.Second, fmt.Sprintf("minute test2"))
	//
	go testwraper(38*time.Second, fmt.Sprintf("mili test1"))

	go testwraper(28*time.Second, fmt.Sprintf("mili test2"))

	go testwraper(18*time.Second, fmt.Sprintf("mili test3"))

	go testwraper(8*time.Second, fmt.Sprintf("mili test4"))

	//
	go testwraper(8*time.Millisecond, fmt.Sprintf("micro test"))

	//c := time.Tick(10*time.Microsecond)
	//wall:=time.Now()
	//<-time.After(10*time.Nanosecond)
	//fmt.Println("init test After:",time.Since(wall))
	//
	//for range c {
	//	interunix:=time.Now().UnixNano()-wall.UnixNano()
	//	inter:=time.Since(wall)
	//	fmt.Println("init test inter:",inter)
	//	fmt.Println("init test interunix:",interunix)
	//	panic("")
	//}

	//go wraper(8*time.Microsecond,fmt.Sprintf("nano test")) //微秒一下的测试需要在Linux下进行 windows 下 go只能提供毫秒级别的精度

	for {
		time.Sleep(time.Second)
		if Stats() == 0 {
			return
		}
		time.Sleep(3 * time.Second)
	}

}

func testwraper(du time.Duration, fmod string) {
	wall := time.Now()
	useExample(du, fmod)
	inter := time.Now().Sub(wall)
	msg := fmt.Sprintf("espect inter:%v   actual inter:%v \n", du, inter)
	if inter < du+(4200*time.Microsecond) && inter > du-(2000*time.Microsecond) { //windows 下误差在2毫秒之间
		fmt.Println(fmod, " passed !!!\n", msg)
	} else {
		panic(fmod + " failed !!!\n" + msg)
	}

}

func useExample(du time.Duration, fmod string) {
	t := st.NewTimerTest(du, fmod) //just like time.NewTimer(du)
	<-t.C
}
