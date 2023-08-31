package timer

import (
	"fmt"
	"log"
	"sync"
	"time"
)

type ExecType int

const (
	OnceExec   ExecType = iota //一次性定时器节点
	PeriodExec                 //周期性定时器
)

const (
	SleepTime     = 100
	DateFmtSecond = "2006-01-02 15:04:05"
	DateFmtDay    = "2006-01-02"
	DateFmtDay2   = "2006.01.02"
	DateFmtHour   = "2006.01.02.15"
)

type Handler func(data interface{})

type Timer struct {
	name      string        //定时器名字
	execType  ExecType      //类型
	handler   Handler       //定时器句柄
	period    time.Duration //定时器的周期
	execTime  time.Time     //触发时间
	data      interface{}   //数据
	execCount int64         //执行次数
}

func NewTimer(name string, execType ExecType, period time.Duration, handler Handler, data interface{}) (*Timer, error) {
	if period < 1*time.Second {
		return nil, fmt.Errorf("NewTimer period must >= 1s")
	}
	ptr := &Timer{
		name:     name,
		execType: execType,
		period:   period,
		handler:  handler,
		data:     data,
	}
	return ptr, nil
}

func (t *Timer) GetExecCount() int64 {
	return t.execCount
}

type Cron struct {
	Timers      map[string]*Timer
	sleepPeriod time.Duration
	lastTime    time.Time
	lock        sync.RWMutex
	isRun       bool
}

func NewCron() *Cron {
	gCtrl := new(Cron)
	gCtrl.sleepPeriod = time.Duration(SleepTime) * time.Millisecond
	gCtrl.Timers = make(map[string]*Timer)
	return gCtrl
}

func (c *Cron) String() string {
	return "timer ctrl "
}

func (c *Cron) addTimer(timer *Timer, force bool) bool {
	if timer == nil {
		return false
	}
	c.lock.Lock()
	defer c.lock.Unlock()

	_, ok := c.Timers[timer.name]
	if !force && ok {
		return false
	}
	c.Timers[timer.name] = timer
	return true
}

func (c *Cron) run() {
	for {
		time.Sleep(c.sleepPeriod)
		c.execute()
	}
}

func (c *Cron) execute() {
	needExecTimers := c.getNeedExecTimers()
	if needExecTimers == nil {
		return
	}
	for _, t := range needExecTimers {
		if t.execType == OnceExec {
			c.DeleteTimer(t.name)
		}
		go t.handler(t.data)
	}
}

func (c *Cron) getNeedExecTimers() []*Timer {
	if c.Timers == nil {
		log.Printf("GetNeedExeTask : Timers is nil")
		return nil
	}
	needExeTaskTimers := make([]*Timer, 0)
	curTime := time.Now()
	c.lock.RLock()
	defer c.lock.RUnlock()

	for _, _timer := range c.Timers {
		// 如果周期时间跟循环时间差不多。那可能会出问题。
		if curTime.After(c.lastTime.Add(2 * c.sleepPeriod)) {
			log.Printf("%s became big from %v to %v ", c, c.lastTime, curTime)
		}

		if curTime.After(_timer.execTime) {
			switch _timer.execType {
			case OnceExec:
				if _timer.execCount > 0 {
					continue
				}
			case PeriodExec:
				if _timer.execType == PeriodExec {
					_timer.execTime = curTime.Add(_timer.period)
				}
			}
			_timer.execCount += 1
			needExeTaskTimers = append(needExeTaskTimers, _timer)
		}
	}
	c.lastTime = curTime
	return needExeTaskTimers
}

func (c *Cron) DeleteTimer(name string) bool {
	c.lock.Lock()
	defer c.lock.Unlock()
	_, ok := c.Timers[name]
	if !ok {
		return false
	}
	delete(c.Timers, name)
	return true
}

func (c *Cron) GetTimer(name string) *Timer {
	c.lock.RLock()
	defer c.lock.RUnlock()
	_timer, ok := c.Timers[name]
	if !ok {
		return nil
	}
	return _timer
}

// AddTimer force:whether to replace names when they are the same
func (c *Cron) AddTimer(timer *Timer, force bool) bool {
	if 2*c.sleepPeriod > timer.period {
		return false
	}
	timer.execTime = time.Now().Add(timer.period)
	return c.addTimer(timer, force)
}

func (c *Cron) Run() {
	c.lock.Lock()
	defer c.lock.Unlock()
	if c.isRun {
		return
	}
	go c.run()
	c.isRun = true
}
