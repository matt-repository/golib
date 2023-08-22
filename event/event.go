package event

import (
	"hash/crc32"
	"log"
	"reflect"
	"runtime"
	"sync"
	"time"
)

type ExecType = int

const (
	HashExec ExecType = iota //hash轮询
	RRExec                   //循环轮询
)

type Handle func(data interface{}) error

type Hash func(key string) int

func (h *Handle) GetName() string {
	return runtime.FuncForPC(reflect.ValueOf(h).Pointer()).Name()
}

//Data ...
type Data struct {
	Key        string
	QueueIndex int
	InTime     time.Time
	OutTime    time.Time
	Data       interface{}
	Handle     Handle //增加一个简化的接口, 用于兼容
	HandleName string
}

//Ctrl ...
type Ctrl struct {
	Name           string
	EventChan      []chan *Data
	ChanBufferSize int64
	QueueIndex     int
	mu             sync.Mutex
	Hash           Hash
	doneChan       chan struct{}
}

func (e *Data) DumpIn(queueIndex int) {
	e.QueueIndex = queueIndex
	e.InTime = time.Now()
	log.Printf("in EventData Key:%s QueueIndex:%d,Handle:%s Data:%v", e.Key, e.QueueIndex, e.HandleName, e.Data)
}

func (e *Data) DumpOut() {
	e.OutTime = time.Now()
	cost := (e.OutTime.UnixNano() - e.InTime.UnixNano()) / int64(time.Microsecond)
	log.Printf("out EventData Key:%s QueueIndex:%d Handle:%s wait:%dus", e.Key, e.QueueIndex, e.HandleName, cost)
}

func (e *Data) DumpStats() {
	cost := (time.Now().UnixNano() - e.OutTime.UnixNano()) / int64(time.Microsecond)
	log.Printf("EventData Key:%s QueueIndex:%d Handle:%s cost:%dus", e.Key, e.QueueIndex, e.HandleName, cost)
}

//DefaultHash ...
func DefaultHash(key string) int {
	return int(crc32.ChecksumIEEE([]byte(key)))
}

//NewData ...
func NewData(Key string, data interface{}, Handle Handle) *Data {
	eventData := &Data{
		Data:       data,
		Key:        Key,
		Handle:     Handle,
		HandleName: Handle.GetName(),
	}
	return eventData
}

//GetQueueIndexByHash ...
func (r *Ctrl) GetQueueIndexByHash(Key string) int {
	r.mu.Lock()
	defer r.mu.Unlock()
	ret := r.Hash(Key) % len(r.EventChan)
	return ret
}

//GetQueueIndexByRR ...
func (r *Ctrl) GetQueueIndexByRR() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	ret := r.QueueIndex
	r.QueueIndex = (r.QueueIndex + 1) % len(r.EventChan)
	return ret
}

//EventPut ....
func (r *Ctrl) EventPut(event *Data, execType ExecType) bool {
	if event == nil {
		return false
	}
	queueIndex := 0
	switch execType {
	case HashExec:
		queueIndex = r.GetQueueIndexByHash(event.Key)
	case RRExec:
		queueIndex = r.GetQueueIndexByRR()
	}
	event.DumpIn(queueIndex)
	r.EventChan[queueIndex] <- event
	return true
}

//run ...
func (r *Ctrl) run(evenChan chan *Data) {
	for {
		select {
		case event := <-evenChan:
			if event == nil {
				continue
			}
			event.DumpOut()
			if event.Handle != nil {
				event.Handle(event.Data)
			}
			event.DumpStats()
		case <-r.doneChan:
			close(evenChan)
			return
		}
	}
}

//Run ...
func (r *Ctrl) Run() {
	for i := 0; i < len(r.EventChan); i++ {
		r.EventChan[i] = make(chan *Data, r.ChanBufferSize)
		log.Printf("evenChan[%d]:[%v},bufferSize:[%d]", i, r.EventChan[i], r.ChanBufferSize)
		go r.run(r.EventChan[i])
	}
}

//Stop ...
func (r *Ctrl) Stop() {
	r.doneChan <- struct{}{}
}

//NewCtrl ...
func NewCtrl(name string, queueCount int, chanBufferSize int64, hash Hash) *Ctrl {
	ctrl := &Ctrl{
		Name:           name,
		EventChan:      make([]chan *Data, queueCount),
		ChanBufferSize: chanBufferSize,
		QueueIndex:     0,
		Hash:           hash,
		doneChan:       make(chan struct{}),
	}
	if ctrl.Hash == nil {
		ctrl.Hash = DefaultHash
	}
	return ctrl
}
