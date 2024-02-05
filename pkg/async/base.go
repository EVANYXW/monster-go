package async

import (
	"fmt"
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

func Go(fn func()) {
	waitAll.Add(1)
	/*var debugStr string
	id := atomic.AddUint32(&goid, 1)
	c := atomic.AddInt32(&gocount, 1)
	if DefLog.Level() <= LogLevelDebug {
		debugStr = LogSimpleStack()
		LogDebug("goroutine start id:%d count:%d from:%s", id, c, debugStr)
	}*/
	atomic.AddInt32(&statis.GoCount, 1)
	go func() {
		defer func() {
			atomic.AddInt32(&statis.GoCount, -1)
			waitAll.Done()
		}()

		Try(fn, nil)

		/*c = atomic.AddInt32(&gocount, -1)

		if DefLog.Level() <= LogLevelDebug {
			LogDebug("goroutine end id:%d count:%d from:%s", id, c, debugStr)
		}*/
	}()
}

func Try(fun func(), handler func(interface{})) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println()
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
