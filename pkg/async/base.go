package async

import (
	"fmt"
	"github.com/evanyxw/monster-go/pkg/output"
	"runtime"
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

func oputAsyncGoUpdate() {
	if output.Oput == nil {
		return
	}
	output.Oput.SetGoNum(GetGoCount())
}

func printStack() {
	// 获取调用栈的信息
	var stack [64]uintptr             // 定义一个较小的栈数组
	n := runtime.Callers(2, stack[:]) // 2表示跳过当前函数和printStack
	if n == 0 {
		return
	}

	// 使用runtime.CallersFrames解析调用栈
	frames := runtime.CallersFrames(stack[:n])
	for {
		frame, more := frames.Next()
		fmt.Println(fmt.Sprintf("File: %s, Line: %d, Function: %s\n", frame.File, frame.Line, frame.Function))
		if !more {
			break
		}
	}
}

func Go(fn func()) {
	waitAll.Add(1)
	//fmt.Println(string(debug.Stack()))
	printStack()
	go func() {
		defer func() {
			waitAll.Done()
			atomic.AddInt32(&statis.GoCount, -1)
			oputAsyncGoUpdate()
		}()

		atomic.AddInt32(&statis.GoCount, 1)
		oputAsyncGoUpdate()
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
