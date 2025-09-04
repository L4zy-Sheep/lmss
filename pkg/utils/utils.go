package utils

import (
	"bytes"
	"errors"
	"lmss/pkg/log"
	"net"
	"slices"
	"sync"
)

var (
	serviceFeature map[string][]byte = map[string][]byte{
		"http":  []byte("HTTP/"),
		"mysql": []byte("MySQL"),
		"redis": []byte("Redis"),
		"ssh":   []byte("SSH"),
	}
)

func IdentifyService(buf []byte) (string, error) {
	for k, v := range serviceFeature {
		if bytes.Contains(buf, v) {
			return k, nil
		}
	}
	return "", errors.New("unknown service, response:\n" + string(buf))
}

func ParseCIDR(cidrs []string, ips chan string) {
	for _, cidr := range cidrs {
		_, ipnet, err := net.ParseCIDR(cidr)
		if err != nil {
			ip := net.ParseIP(cidr)
			if ip == nil {
				log.Error("parse error: " + cidr)
			} else {
				ips <- ip.String()
			}
			continue
		}
		ip := slices.Clone(ipnet.IP)
		for {
			ips <- ip.String()
			for i := 3; i >= 0; i-- {
				ip[i]++
				if ip[i] > 0 {
					break
				}
			}
			if !ipnet.Contains(ip) {
				break
			}
		}
	}
	close(ips)
}

type ThreadPool struct {
	f  chan func()
	wg *sync.WaitGroup
}

func (t *ThreadPool) Start(f func()) {
	t.f <- f
}
func (t *ThreadPool) Stop() {
	close(t.f)
	t.wg.Wait()
}
func NewPool(num int) *ThreadPool {
	tp := &ThreadPool{
		f:  make(chan func(), num*2),
		wg: &sync.WaitGroup{},
	}
	for i := 0; i < num; i++ {
		tp.wg.Add(1)
		go func() {
			for f := range tp.f {
				f()
			}
			tp.wg.Done()
		}()
	}
	return tp
}
