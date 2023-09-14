package event

import (
	"fmt"
	"time"
)

func ExampleEventPut() {
	ctrl := NewCtrl("test1", 10, 256, DefaultHash)
	ctrl.Run()
	for i := 0; i < 20; i++ {
		event := NewData(fmt.Sprintf("event:[%d]", i), fmt.Sprintf("data:[%d]", i), func(data interface{}) error {
			time.Sleep(1 * time.Second)
			return nil
		})
		ctrl.EventPut(event, RRExec)
	}

	time.Sleep(5 * time.Second)
	fmt.Println("a")
	//OutPut:a
}
