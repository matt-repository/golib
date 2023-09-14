package event

import (
	"hash/crc32"
	"log"
	"sync"
	"time"
)

type ExecType = int

const (
	HashExec ExecType = iota //hash轮询
	RRExec                   //循环轮询
)

type Handle func(data any) error

type Hash func(key string) int

// Data ...
type Data struct {
	Key         string
	QueueIndex  int
	InQueueTime time.Time
	Data        any
	Handle      Handle //增加一个简化的接口, 用于兼容
}

// Ctrl ...
type Ctrl struct {
	Name           string
	EventChan      []chan *Data
	ChanBufferSize int64
	QueueIndex     int
	mu             sync.Mutex
	Hash           Hash
}

// DefaultHash ...
func DefaultHash(k string) int {
	return int(crc32.ChecksumIEEE([]byte(k)))
}

// NewData ...
func NewData(key string, data any, handle Handle) *Data {
	eventData := &Data{
		Data:   data,
		Key:    key,
		Handle: handle,
	}
	return eventData
}

// GetQueueIndexByHash ...
func (c *Ctrl) GetQueueIndexByHash(key string) int {
	c.mu.Lock()
	defer c.mu.Unlock()
	ret := c.Hash(key) % len(c.EventChan)
	return ret
}

// GetQueueIndexByRR ...
func (c *Ctrl) GetQueueIndexByRR() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	ret := c.QueueIndex
	c.QueueIndex = (c.QueueIndex + 1) % len(c.EventChan)
	return ret
}

// EventPut ....
func (c *Ctrl) EventPut(data *Data, execType ExecType) bool {
	if data == nil {
		return false
	}
	switch execType {
	case HashExec:
		data.QueueIndex = c.GetQueueIndexByHash(data.Key)
	case RRExec:
		data.QueueIndex = c.GetQueueIndexByRR()
	}
	data.InQueueTime = time.Now()
	//log in queue time
	log.Printf("EventData in :%+v", data.Data)

	c.EventChan[data.QueueIndex] <- data
	return true
}

// run ...
func (c *Ctrl) run(data chan *Data) {
	for {
		e := <-data
		if e == nil {
			continue
		}
		//log in queue => out queue time
		execStart := time.Now()
		log.Printf("EventData ExecStart:%+v, wait cost:%fs", e.Data, execStart.Sub(e.InQueueTime).Seconds())
		if e.Handle != nil {
			e.Handle(e.Data)
		}
		//log out queue => exec end time
		execEnd := time.Now()
		log.Printf("EventData ExecEnd:%+v, exec cost:%fs", e.Data, execEnd.Sub(execStart).Seconds())
	}
}

// Run ...
func (c *Ctrl) Run() {
	for i := 0; i < len(c.EventChan); i++ {
		c.EventChan[i] = make(chan *Data, c.ChanBufferSize)
		log.Printf("evenChan[%d]:bufferSize:[%d]", i, c.ChanBufferSize)
		go c.run(c.EventChan[i])
	}
}

// NewCtrl ...
func NewCtrl(name string, queueCount int, chanBufferSize int64, hash Hash) *Ctrl {
	ctrl := &Ctrl{
		Name:           name,
		EventChan:      make([]chan *Data, queueCount),
		ChanBufferSize: chanBufferSize,
		QueueIndex:     0,
		Hash:           hash,
	}
	if ctrl.Hash == nil {
		ctrl.Hash = DefaultHash
	}
	return ctrl
}
