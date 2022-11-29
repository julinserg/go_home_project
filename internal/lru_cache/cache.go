package lrucache

import "sync"

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity    int
	queue       List
	items       map[Key]*ListItem
	itemsHelper map[*ListItem]Key
	mu          sync.Mutex
}

func (c *lruCache) Set(key Key, value interface{}) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	if valListItem, isExist := c.items[key]; isExist {
		valListItem.Value = value
		c.queue.MoveToFront(valListItem)
		return true
	}
	if c.queue.Len() == c.capacity {
		backElement := c.queue.Back()
		c.queue.Remove(backElement)
		keyForBackElement := c.itemsHelper[backElement]
		delete(c.items, keyForBackElement)
		delete(c.itemsHelper, backElement)
	}
	element := c.queue.PushFront(value)
	c.items[key] = element
	c.itemsHelper[element] = key
	return false
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if valListItem, isExist := c.items[key]; isExist {
		c.queue.MoveToFront(valListItem)
		return valListItem.Value, true
	}
	return nil, false
}

func (c *lruCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items = make(map[Key]*ListItem, c.capacity)
	c.itemsHelper = make(map[*ListItem]Key, c.capacity)
	c.queue = NewList()
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity:    capacity,
		queue:       NewList(),
		items:       make(map[Key]*ListItem, capacity),
		itemsHelper: make(map[*ListItem]Key, capacity),
	}
}
