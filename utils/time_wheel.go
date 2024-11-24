package utils

import (
	"fmt"
	"sync"
	"time"
)

type TimeWheel interface {
	BaseTickDuration() time.Duration
	TickBaseUnits() int64
	Start()
	Reset()
	TickChan() <-chan time.Time
	WheelSize() int64
	CurrentTime() int64
	IsBaseTicker() bool
	LowestTimeWheel() TimeWheel
	SetUpperTimeWheel(tw TimeWheel)
	TryAddTask(task TimingTask)
	AddTask(task TimingTask)
}

type TimeTicker struct {
	ticker       *time.Ticker
	tickDuration time.Duration
}

func (t *TimeTicker) BaseTickDuration() time.Duration {
	return t.tickDuration
}

func (t *TimeTicker) TickBaseUnits() int64 {
	return 1
}

func (t *TimeTicker) Reset() {
	t.ticker.Reset(t.tickDuration)
}

func (t *TimeTicker) TickChan() <-chan time.Time {
	return t.ticker.C
}

func (t *TimeTicker) Start() {
	t.ticker = time.NewTicker(t.tickDuration)
}

func (t *TimeTicker) CurrentTime() int64 {
	return 0
}

func (t *TimeTicker) WheelSize() int64 {
	return 1
}

func (t *TimeTicker) AddTask(task TimingTask) {}

func (t *TimeTicker) TryAddTask(task TimingTask) {}

func (t *TimeTicker) IsBaseTicker() bool {
	return true
}

func (t *TimeTicker) LowestTimeWheel() TimeWheel {
	return nil
}

func (t *TimeTicker) SetUpperTimeWheel(tw TimeWheel) {}

func NewTimeTicker(duration time.Duration) TimeWheel {
	return &TimeTicker{
		tickDuration: duration,
	}
}

type TimingWheel struct {
	wheelSize       int64
	subTimeWheel    TimeWheel
	upperTimeWheel  TimeWheel
	lowestTimeWheel TimeWheel
	currentTime     int64
	timerTaskLists  [][]TimingTask
	taskRwMutex     *sync.RWMutex
	tickBaseUnits   int64
	intervalCh      chan time.Time
}

// 3层时间轮，每层走完1轮会触发上层时间轮走一格，第1层1ms一格，时间跨度为1s，第二层时间跨度为1分钟，第三层时间跨度为1小时
var stdTimingWheel = NewTimingWheel(60, NewTimingWheel(60, NewTimingWheel(1000, NewTimeTicker(time.Millisecond))))

func Instance() *TimingWheel {
	return stdTimingWheel
}

func NewTimingWheel(wheelSize int64, tw TimeWheel) *TimingWheel {
	wheel := &TimingWheel{
		wheelSize:      wheelSize,
		subTimeWheel:   tw,
		timerTaskLists: make([][]TimingTask, wheelSize),
		taskRwMutex:    new(sync.RWMutex),
		tickBaseUnits:  tw.TickBaseUnits() * tw.WheelSize(),
		intervalCh:     make(chan time.Time),
	}
	if tw.IsBaseTicker() {
		wheel.lowestTimeWheel = wheel
	} else {
		wheel.lowestTimeWheel = wheel.subTimeWheel.LowestTimeWheel()
	}
	tw.SetUpperTimeWheel(wheel)
	for i := range wheel.timerTaskLists {
		wheel.timerTaskLists[i] = make([]TimingTask, 0)
	}
	return wheel
}

// BaseTickDuration implements TimeWheel.
func (t *TimingWheel) BaseTickDuration() time.Duration {
	return t.subTimeWheel.BaseTickDuration()
}

func (t *TimingWheel) SetUpperTimeWheel(tw TimeWheel) {
	t.upperTimeWheel = tw
}

func (t *TimingWheel) TickBaseUnits() int64 {
	if t.tickBaseUnits == 0 {
		t.tickBaseUnits = t.subTimeWheel.TickBaseUnits() * t.subTimeWheel.WheelSize()
	}
	return t.tickBaseUnits
}

func (t *TimingWheel) MaxTaskDuration() time.Duration {
	return t.BaseTickDuration() * time.Duration(t.TickBaseUnits()) * time.Duration(t.WheelSize())
}

func (t *TimingWheel) Reset() {
	t.subTimeWheel.Reset()
}

func (t *TimingWheel) Start() {
	t.subTimeWheel.Start()
	tickCh := t.subTimeWheel.TickChan()
	for tt := range tickCh {
		t.currentTime++
		if t.currentTime >= t.wheelSize {
			t.currentTime %= t.wheelSize
			t.intervalCh <- tt
		}
		t.TriggerCurrentTimeTasks()
	}
}

func (t *TimingWheel) TriggerCurrentTimeTasks() {
	t.taskRwMutex.RLock()
	defer t.taskRwMutex.RUnlock()
	for _, task := range t.timerTaskLists[t.currentTime] {
		tickNum := task.TimeUnitsLeft() / t.TickBaseUnits()
		task.CurrentSpendBaseTicks += tickNum * t.TickBaseUnits()
		// TODO
		leftTime := task.TimeUnitsLeft()
		if leftTime == 0 {

		}
	}
}

func (t *TimingWheel) TickChan() <-chan time.Time {
	return t.intervalCh
}

func (t *TimingWheel) WheelSize() int64 {
	return t.wheelSize
}

func (t *TimingWheel) CurrentTime() int64 {
	return t.currentTime*t.subTimeWheel.WheelSize() + t.subTimeWheel.CurrentTime()
}

type TimerTaskDealer interface {
	ExecByTime(triggerTime time.Time)
}

type TimingTask struct {
	TriggerCh             chan time.Time
	OnlyOnce              bool
	IntervalBaseTicks     int64
	CurrentSpendBaseTicks int64
	Dealer                TimerTaskDealer
}

func (t TimingTask) TimeUnitsLeft() int64 {
	return t.IntervalBaseTicks - t.CurrentSpendBaseTicks
}

func (t *TimingTask) Trigger(triggerTime time.Time) {
	t.CurrentSpendBaseTicks = 0
	t.Dealer.ExecByTime(triggerTime)
}

func (t TimingWheel) WrapTaskDealer(intervalDuration time.Duration, onlyOnce bool, dealer TimerTaskDealer) (task TimingTask, err error) {
	if intervalDuration > t.MaxTaskDuration() {
		return task, fmt.Errorf("interval duration is out of limit: %d sec", t.MaxTaskDuration()/time.Second)
	}
	if intervalDuration < t.BaseTickDuration() {
		return task, fmt.Errorf("interval duration is smaller than the base tick duration: %d milsec", t.BaseTickDuration()/time.Millisecond)
	}
	return TimingTask{
		TriggerCh:         make(chan time.Time),
		OnlyOnce:          onlyOnce,
		IntervalBaseTicks: int64(intervalDuration / t.BaseTickDuration()),
		Dealer:            dealer,
	}, nil
}

func (t *TimingWheel) AddTask(task TimingTask) {
	t.lowestTimeWheel.TryAddTask(task)
}

func (t *TimingWheel) TryAddTask(task TimingTask) {
	leftTime := task.TimeUnitsLeft()
	if t.upperTimeWheel != nil && leftTime > t.WheelSize()*t.TickBaseUnits() {
		// 如果任务剩余等待时间超过当前时间轮走一轮所需的时间，则尝试放入上层时间轮中
		t.upperTimeWheel.TryAddTask(task)
	}
	// 任务剩余等待时间在当前时间轮的周期时间内，尝试放入当前时间轮
	tickNum := leftTime / t.TickBaseUnits()
	timeIdx := (t.currentTime + tickNum) % t.wheelSize
	t.timerTaskLists[timeIdx] = append(t.timerTaskLists[timeIdx], task)
}

func (t *TimingWheel) IsBaseTicker() bool {
	return false
}

func (t *TimingWheel) LowestTimeWheel() TimeWheel {
	return t.lowestTimeWheel
}
