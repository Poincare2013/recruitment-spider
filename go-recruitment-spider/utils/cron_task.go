package utils

import (
	"time"
)

func CronTask(f func(),d time.Duration) {
	// 这是一个使用time包实现的定时器
	t1 := time.NewTimer(d)
	for {
		select {
		case <-t1.C:
			//t1.Reset(time.Second * 180)
			t1.Reset(d)
			f()
		}
	}
}
