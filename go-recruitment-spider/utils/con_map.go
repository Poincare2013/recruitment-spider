package utils

import "sync"

type ConMap struct {
	Data map[string]int
	*sync.Mutex
}

func (c *ConMap) Get(k string) (int, bool) {
	c.Lock()
	defer c.Unlock()
	if v, ok := c.Data[k]; ok {
		return v, true
	}
	return 0, false
}

func (c *ConMap) Set(k string, v int) {
	c.Lock()
	defer c.Unlock()
	c.Data[k] = v
}

func (c *ConMap) Delete(k string) {
	c.Lock()
	defer c.Unlock()
	delete(c.Data, k)
}
