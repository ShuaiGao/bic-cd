package limiter

import (
	"time"
)

// LocalLimiter 计数器算法，进程内限流
// 支持按优先级排序触发，支持以最小限流频次倍增限流
// 使用方法:
//  1. 通过 NewLocalLimiter 方法新建限流器
//  2. 通过 Run 添加要限流的函数 f
//     在f执行时会增加限流计数，限流计数标识当前正在执行的函数总数
//     若正在执行的函数总数，达到频次上线，则被限流
type LocalLimiter struct {
	l *localLimiter
}

type localLimiter struct {
	itemTail       *localLimiterItem
	itemHead       *localLimiterItem
	addC           chan *localLimiterItem
	doingC         chan bool
	quitChan       chan bool
	tickerC        *time.Ticker
	tickerDuration time.Duration
	frequency      int
	duration       time.Duration
	tickCounter    int
	emptyLoop      int
}
type localLimiterItem struct {
	priority int
	multiple float32
	f        func()
	next     *localLimiterItem
}

// funcDone 函数执行完成，返还限流计数
func (l localLimiter) funcDone() {
	<-l.doingC
}

// RunMulti 以go routine方式执行一个函数f
// 参数 priority 排序优先级，限流器将优先执行优先级高的函数
// 参数 multiple 当前函数执行要占用的限流计数倍数
func (l LocalLimiter) RunMulti(priority int, multiple float32, f func()) {
	if multiple < 1 {
		multiple = 1
	}
	item := &localLimiterItem{
		priority: priority,
		multiple: multiple,
		f:        f,
	}
	l.l.addC <- item
}

func (l LocalLimiter) Quit() {
	l.l.quitChan <- true
}

// Run 以最低优先级，1倍限流占用率执行 f
func (l LocalLimiter) Run(f func()) {
	l.RunMulti(0, 1, f)
}

// NewLocalLimiter 新建本地限流器
// 参数 frequency 每秒请求的频次
// 例：
//
//	frequency = 2，即限制每1/2秒(0.5秒)触发一次
//	frequency = 0.5，即限制每1/0.5秒(2秒)触发一次
func NewLocalLimiter(frequency float64) Limiter {
	duration := time.Millisecond*time.Duration(1000/frequency)/10 + 1
	doingChanLen := int(frequency)
	if frequency < 1 {
		doingChanLen = 1
	}
	l := &LocalLimiter{}
	l.l = &localLimiter{
		duration:       duration,
		itemHead:       &localLimiterItem{},
		tickerC:        time.NewTicker(duration),
		tickerDuration: duration,
		doingC:         make(chan bool, doingChanLen),
		addC:           make(chan *localLimiterItem),
	}
	l.l.itemTail = l.l.itemHead
	l.l.itemHead.next = l.l.itemTail
	go l.l.run()
	return l
}

func (l *localLimiter) run() {
	for {
		select {
		case <-l.quitChan:
			return
		case item := <-l.addC:
			l.add(item)
		case <-l.tickerC.C:
			l.do()
		}
	}
}

func (l *localLimiter) do() {
	if l.itemHead == l.itemTail {
		if l.emptyLoop == 10 {
			l.tickerC.Reset(l.tickerDuration * 10)
		}
		l.emptyLoop++
		return
	}
	l.emptyLoop = 0
	l.tickCounter++
	item := l.itemHead.next
	if item.multiple*10 > float32(l.tickCounter) {
		return
	}
	if cap(l.doingC) < int(float32(len(l.doingC))*item.multiple) {
		// 最大限流到上限
		return
	}

	l.tickCounter = 0
	l.itemHead.next = item.next
	if item == l.itemTail {
		l.itemTail = l.itemHead
	}
	l.doingC <- true
	go func() {
		defer l.funcDone()
		item.f()
	}()
}

func (l *localLimiter) add(item *localLimiterItem) {
	if l.itemTail != l.itemHead && item.priority > 0 {
		p := l.itemHead
		for p.next != nil {
			if p.next.priority < item.priority {
				item.next = p.next
				p.next = item
				// 插入成功直接返回
				return
			}
			p = p.next
		}
	}
	// 在尾节点前直接插入
	l.itemTail.next = item
	l.itemTail = item
	if l.emptyLoop >= 10 {
		l.tickerC.Reset(l.tickerDuration)
		l.emptyLoop = 0
	}
}
