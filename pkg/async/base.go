package async

import (
	"github.com/evanyxw/monster-go/pkg/output"
	"sync"
	"sync/atomic"
	"time"
)

var waitAll = &sync.WaitGroup{} //等待所有goroutine
type Statis struct {
	GoCount     int32
	MsgqueCount int
	StartTime   time.Time
	LastPanic   int
	PanicCount  int32
}

var statis = &Statis{}

func GetGoCount() int32 {
	return statis.GoCount
}

func Wait() {
	waitAll.Wait()
}

func OputUpdate() {
	if output.Oput == nil {
		return
	}
	output.Oput.SetGoNum(GetGoCount())
}

func Go(fn func()) {
	waitAll.Add(1)
	//fmt.Println(string(debug.Stack()))

	go func() {
		defer func() {
			waitAll.Done()
			atomic.AddInt32(&statis.GoCount, -1)
			OputUpdate()
		}()

		atomic.AddInt32(&statis.GoCount, 1)
		OputUpdate()
		Try(fn, nil)
	}()
}

func Try(fun func(), handler func(interface{})) {
	defer func() {
		if err := recover(); err != nil {
			if handler == nil {
				//LogStack()
				//LogError("error catch:%v", err)

			} else {
				handler(err)
			}
			atomic.AddInt32(&statis.PanicCount, 1)
			//statis.LastPanic = int(Timestamp)
		}
	}()
	fun()
}
