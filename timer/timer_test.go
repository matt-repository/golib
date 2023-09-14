package timer

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"log"
	"testing"
	"time"
)

func TestTimer(t *testing.T) {
	ctrl := NewCron()
	count := 0
	_timer, err := NewTimer("test", PeriodExec, 1*time.Second, func(data any) {
		log.Printf(fmt.Sprintf("time:%s,count:%d", time.Now().Format(DateFmtSecond), count))
		count++
	}, nil)
	if err != nil {
		log.Printf(err.Error())
		return
	}
	ctrl.addTimer(_timer, true)
	ctrl.Run()
	select {}
}

func TestCron(t *testing.T) {
	count := 0
	cr := cron.New(cron.WithSeconds())
	cr.AddFunc("@every 1s", func() {
		count++
		fmt.Println(fmt.Sprintf("time:%s,count:%d", time.Now().Format(DateFmtSecond), count))
	})
	cr.Start()
	select {}
}
