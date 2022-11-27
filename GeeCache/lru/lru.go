package lru

import "container/list"

type Cache struct {
	maxBytes  int64
	nBytes    int64
	ll        *list.List                    //这里是一个双向量表存数据
	cache     map[string]*list.Element      // map存储数据索引
	OnEvicted func(key string, value Value) // 回调函数
}

// entry 实际存储在双向链表中的数据
type entry struct {
	key   string
	value Value
}

// Value 接口给存入的数据留下了足够的空间，只要是实现了当前接口的数据都可以存储
type Value interface {
	Len() int
}

func New(maxBytes int64, onEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}

// Get 查找功能 bare return
func (c *Cache) Get(key string) (value Value, Bool bool) {
	if ele, ok := c.cache[key]; ok {
		// 访问该element之后将其提到front
		c.ll.MoveToFront(ele)
		// 这里的Value并非我们自定义的Value
		// 类型断言x.(T)
		// T是结构体x.(*T) T是接口x.(T)
		kv := ele.Value.(*entry)
		return kv.value, true
	}
	return
}

// RemoveOldest 淘汰缓存
func (c *Cache) RemoveOldest() {
	// 取出最旧的数据
	ele := c.ll.Back()
	if ele != nil {
		// 移除指定element
		c.ll.Remove(ele)
		kv := ele.Value.(*entry)
		// 删除map中的key value
		delete(c.cache, kv.key)
		// 更新占用bytes
		c.nBytes -= int64(len(kv.key)) + int64(kv.value.Len())
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}
	}
}

// AddUpdate 修改删除数据
func (c *Cache) AddUpdate(key string, value Value) {
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		c.nBytes += int64(value.Len()) - int64(kv.value.Len())
		kv.value = value
	} else {
		ele := c.ll.PushFront(&entry{key: key, value: value})
		c.cache[key] = ele
		c.nBytes += int64(len(key)) + int64(value.Len())
	}
	// 之后检测缓存是否溢出
	for c.maxBytes != 0 && c.maxBytes < c.nBytes {
		c.RemoveOldest()
	}
}

func (c *Cache) Len() int {
	return c.ll.Len()
}
