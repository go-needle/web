package web

import (
	"net/http"
	"strings"
	"sync"
	"time"
)

// RateLimit is a middleware which limits the frequency of access to the same IP address
func RateLimit(rate time.Duration) Handler {
	var blackList sync.Map
	return HandlerFunc(func(c *Context) {
		remoteAddr := c.Request.RemoteAddr
		forwardedFor := c.GetHeader("X-Forwarded-For")
		if forwardedFor != "" {
			remoteAddr = strings.Split(forwardedFor, ", ")[0]
		}
		if _, ok := blackList.Load(remoteAddr); !ok {
			blackList.Store(remoteAddr, struct{}{})
			go func() {
				time.Sleep(rate)
				blackList.Delete(remoteAddr)
			}()
		} else {
			c.Fail(http.StatusForbidden, "rate out of limit")
			return
		}
		c.Next()
	})
}

// TrafficLimit is a middleware which uses token bucket algorithm for traffic restriction
func TrafficLimit(flowTotal int, rate time.Duration) Handler {
	chg := make(chan struct{}, flowTotal)
	go func(ch chan<- struct{}) {
		for i := 0; i < flowTotal; i++ {
			ch <- struct{}{}
		}
		for {
			time.Sleep(rate)
			ch <- struct{}{}
		}
	}(chg)
	return HandlerFunc(func(c *Context) {
		var ch <-chan struct{}
		ch = chg
		select {
		case <-ch:
			c.Next()
		default:
			c.Fail(http.StatusForbidden, "traffic congestion")
			return
		}
	})
}
