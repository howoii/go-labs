package lru

import (
	"container/list"
	"sync"
)

type Cache struct {
	MaxEntries int
	OnEvicted  func(key Key, value interface{})

	ll        *list.List
	cache     map[Key]*list.Element
	entryPool sync.Pool
	usePool   bool
}

type Key interface{}

type entry struct {
	key   Key
	value interface{}
	data  [1024]byte
}

type Option struct {
	UsePool bool
}

func New(maxEntries int, opt *Option) *Cache {
	c := &Cache{
		MaxEntries: maxEntries,
		ll:         list.New(),
		cache:      make(map[Key]*list.Element, 0),
	}
	if opt != nil && opt.UsePool {
		c.usePool = true
		c.entryPool = sync.Pool{New: newPoolEntry}
	}
	return c
}

func newPoolEntry() interface{} {
	return &entry{}
}

func (c *Cache) Add(key Key, value interface{}) {
	if c.cache == nil {
		c.ll = list.New()
		c.cache = make(map[Key]*list.Element, 0)
	}
	if e, ok := c.cache[key]; ok {
		c.ll.MoveToFront(e)
		e.Value.(*entry).value = value
		return
	}

	if c.ll.Len() == c.MaxEntries {
		c.RemoveOldest()
	}

	e := c.ll.PushFront(c.newEntry(key, value))
	c.cache[key] = e
}

func (c *Cache) Get(key Key) (value interface{}, ok bool) {
	if c.cache == nil {
		return
	}
	if e, ok := c.cache[key]; ok {
		c.ll.MoveToFront(e)
		return e.Value.(*entry).value, true
	}

	return
}

func (c *Cache) Remove(key Key) {
	if c.cache == nil {
		return
	}
	if e, ok := c.cache[key]; ok {
		c.removeElement(e)
	}
}

func (c *Cache) RemoveOldest() {
	if c.cache == nil {
		return
	}
	e := c.ll.Back()
	if e != nil {
		c.removeElement(e)
	}
}

func (c *Cache) Len() int {
	if c.cache == nil {
		return 0
	}
	return c.ll.Len()
}

func (c *Cache) Clear() {
	if c.OnEvicted != nil {
		for _, e := range c.cache {
			ent := e.Value.(*entry)
			c.OnEvicted(ent.key, ent.value)
		}
	}
	c.ll = nil
	c.cache = nil
}

func (c *Cache) newEntry(key Key, value interface{}) *entry {
	if c.usePool {
		v := c.entryPool.Get()
		ent := v.(*entry)
		ent.key = key
		ent.value = value
		return ent
	}
	return &entry{key: key, value: value}
}

func (c *Cache) removeElement(e *list.Element) {
	ent := e.Value.(*entry)
	c.ll.Remove(e)
	delete(c.cache, ent.key)

	if c.OnEvicted != nil {
		c.OnEvicted(ent.key, ent.value)
	}

	if c.usePool {
		c.entryPool.Put(ent)
	}
}
