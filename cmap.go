// cmap.go
package util

import "sync"

type CMap struct {
	rwMutex sync.RWMutex
	_map    map[interface{}]interface{}
}

func NewCMap(size int) *CMap {
	m := make(map[interface{}]interface{}, size)
	return &CMap{_map: m}
}

//获取长度
func (cm *CMap) Lenght() int {
	cm.rwMutex.RLock()
	l := len(cm._map)
	cm.rwMutex.RUnlock()
	return l
}

//查找key对应的value  没有返回nil
func (cm *CMap) Get(key interface{}) interface{} {
	cm.rwMutex.RLock()
	v, ok := cm._map[key]
	cm.rwMutex.RUnlock()
	if ok {
		return v
	}
	return nil
}

//存key-value
func (cm *CMap) Put(key interface{}, value interface{}) {
	cm.rwMutex.Lock()
	cm._map[key] = value
	cm.rwMutex.Unlock()
}

//如果不存在key,存进去，返回true  如果存在key，不存进去,返回false
func (cm *CMap) PutNotExist(key interface{}, value interface{}) bool {
	if cm.HasKey(key) {
		return false
	}
	cm.rwMutex.Lock()
	if _, ok := cm._map[key]; ok {
		cm.rwMutex.Unlock()
		return false
	}
	cm._map[key] = value
	cm.rwMutex.Unlock()
	return true
}

//替换key-value  返回旧值(没有为nil)
func (cm *CMap) Replace(key interface{}, value interface{}) interface{} {
	cm.rwMutex.RLock()
	v, ok := cm._map[key]
	cm._map[key] = value
	cm.rwMutex.RUnlock()
	if ok {
		return v
	}
	return nil
}

//检测是否有key值
func (cm *CMap) HasKey(key interface{}) bool {
	cm.rwMutex.RLock()
	_, ok := cm._map[key]
	cm.rwMutex.RUnlock()
	return ok
}

//根据key删除
func (cm *CMap) Delete(key interface{}) {
	if cm.HasKey(key) {
		cm.rwMutex.Lock()
		delete(cm._map, key)
		cm.rwMutex.Unlock()
	}
}

//range删除一个返回 没有返回nil
func (cm *CMap) DeleteAnyOne() (interface{}, interface{}) {
	cm.rwMutex.Lock()
	var key interface{}
	var value interface{}
	for k, v := range cm._map {
		key = k
		value = v
		break
	}
	if key != nil {
		delete(cm._map, key)
	}
	cm.rwMutex.Unlock()
	return key, value
}

//得到所有key的列表
func (cm *CMap) GetKeys() []interface{} {
	cm.rwMutex.RLock()
	ls := make([]interface{}, 0, len(cm._map))
	for k, _ := range cm._map {
		ls = append(ls, k)
	}
	cm.rwMutex.RUnlock()
	return ls
}

//得到所有value的列表
func (cm *CMap) GetValues() []interface{} {
	cm.rwMutex.RLock()
	ls := make([]interface{}, 0, len(cm._map))
	for _, v := range cm._map {
		ls = append(ls, v)
	}
	cm.rwMutex.RUnlock()
	return ls
}

//复制一份map副本
func (cm *CMap) GetAll() map[interface{}]interface{} {
	cm.rwMutex.RLock()
	m := make(map[interface{}]interface{}, len(cm._map))
	for k, v := range cm._map {
		m[k] = v
	}
	cm.rwMutex.RUnlock()
	return m
}

//在所有key-value上执行f方法
func (cm *CMap) Foreach(f func(interface{}, interface{})) {
	cm.rwMutex.RLock()
	for k, v := range cm._map {
		f(k, v)
	}
	cm.rwMutex.RUnlock()
}

//在所有key-value上执行f方法并返回f方法的值的列表
func (cm *CMap) Map(f func(interface{}, interface{}) interface{}) []interface{} {
	cm.rwMutex.RLock()
	ls := make([]interface{}, 0, len(cm._map))
	for k, v := range cm._map {
		ls = append(ls, f(k, v))
	}
	cm.rwMutex.RUnlock()
	return ls
}
