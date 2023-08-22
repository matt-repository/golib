package event

import (
	"fmt"
	"log"
	"time"
)

func ExampleEventPut() {
	ctrl := NewCtrl("test1", 10, 256, DefaultHash)
	ctrl.Run()
	for i := 0; i < 20; i++ {
		event := NewData(fmt.Sprintf("event:[%d]", i), fmt.Sprintf("data:[%d]", i), func(data interface{}) error {
			d := data.(string)
			log.Printf("start : %s", d)
			time.Sleep(2 * time.Second)
			log.Printf("end : %s", d)
			return nil
		})
		ctrl.EventPut(event, RRExec)
	}
	//测试停止
	go func() {
		time.Sleep(2 * time.Second)
		ctrl.Stop()
	}()

	time.Sleep(7 * time.Second)
	fmt.Println("a")
	//OutPut:a
}
