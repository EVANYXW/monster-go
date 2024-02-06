package retry

import (
	"fmt"
	"time"
)

type Retry struct {
	retryCount int           // 最大重试次数
	retryDelay time.Duration // 重试之间的延迟
}

type RetriableError struct {
	error error
}

func (e RetriableError) Error() string {
	return e.error.Error()
}

// 重试方法
func (r *Retry) Retry(fn func() error) error {
	var err error
	retryCount := 0

	for {
		err = fn()

		//调用成功
		if err == nil {
			return nil
		}

		//错误属于重试
		if _, ok := err.(RetriableError); !ok {
			return err
		}

		//已经超过重试限制
		if retryCount > r.retryCount-1 {
			break
		}

		fmt.Println(1)
		retryCount++
		time.Sleep(r.retryDelay)
	}
	return err
}
