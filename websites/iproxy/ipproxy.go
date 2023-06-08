package iproxy

import (
	"log"
	"sync"
	"time"
)

type Proxy struct {
	Scheme string `json:"scheme"`
	Ip     string `json:"ip"`
	Port   int    `json:"port"`
}

var IpUseNum int

type ipProxy interface {
	// 刷新ip
	refresh()
	// 获取代理ip
	get() Proxy
}

type iproxy struct {
	f    RefreshFunc
	ch   chan Proxy
	lock sync.Mutex
}

var inst ipProxy

type RefreshFunc func() []Proxy

func Setup(f RefreshFunc) {
	inst = &iproxy{
		f:  f,
		ch: make(chan Proxy, 1000000),
	}
}

func (i *iproxy) refresh() {
	i.lock.Lock()
	defer i.lock.Unlock()
start:
	if len(i.ch) > 0 {
		return
	}
	for _, s := range i.f() {
		for n := 0; n < IpUseNum; n++ {
			i.ch <- s
		}
	}
	log.Printf("获取到了%d个代理ip", len(i.ch))
	if len(i.ch) == 0 {
		time.Sleep(30 * time.Minute)
		goto start
	}
}

func (i *iproxy) get() Proxy {
	select {
	case s := <-i.ch:
		return s
	default:
		i.refresh()
	}
	return <-i.ch
}

var Enable = true

func Get() Proxy {
	if Enable && inst != nil {
		return inst.get()
	}
	return Proxy{}
}
