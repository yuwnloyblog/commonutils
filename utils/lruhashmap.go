package utils

import(
	"time"
	"sync"
	"sync/atomic"
)
type TimestampEntry struct {
	entry interface{}
	timestamp int64
}

func (self *TimestampEntry)updateTime(){
	self.timestamp = time.Now().Unix()
}

type LruHashMap struct{
	MaxSize int
	Duration int64
	entryMap map[interface{}]*TimestampEntry
	lock sync.RWMutex
	size int
	keyChan chan interface{}
	isCleanerRuning int32
}
/**
 * create hashmap with maxSize and duration
 */
func NewLruHashMap(maxSize int,duration int64)*LruHashMap{
	return &LruHashMap{
		MaxSize : maxSize,
		Duration : duration,
		entryMap : make(map[interface{}]*TimestampEntry),
		keyChan : make(chan interface{}, maxSize),
	}
}

func NewLruHashMapNoDura(maxSize int)*LruHashMap{
	return NewLruHashMap(maxSize,0)
}

func NewDefaultLruHashMap()*LruHashMap{
	return NewLruHashMapNoDura(10000)
}

/**
 * remove expire item goroutine
 */
func cleanerRuning(lru *LruHashMap){
	if lru.Duration>0{
		time.Sleep(time.Duration(lru.Duration)*time.Second)
		var interval int64
		for {
			interval = 1
			lru.lock.Lock()
			currentTime := time.Now().Unix()
			for k,v := range lru.entryMap{
				lapsed:=currentTime - v.timestamp
				if lapsed > lru.Duration{
					//remove key
					lru.removeNoLock(k)
				}else{
					interval = lru.Duration - lapsed
					break
				}
			}
			lru.lock.Unlock()
			time.Sleep(time.Duration(interval)*time.Second)	
		}
	}
}

/**
 * put obj to map
 */
func (self *LruHashMap)Put(key interface{}, entry interface{})interface{}{
	self.lock.Lock()
	self.lock.Unlock()
	old := self.getNoLock(key)
	self.putNoLock(key,entry)
	return old
}

func (self *LruHashMap)PutIfAbsent(key interface{}, entry interface{})interface{}{
	self.lock.Lock()
	self.lock.Unlock()
	old := self.getNoLock(key)
	if old == nil{
		self.putNoLock(key,entry)
	}
	return old
}

func (self *LruHashMap)putNoLock(key interface{}, entry interface{}){
	if entry == nil{
		return
	}
	if self.size>=self.MaxSize{
		//remove oldest entry
		self.removeNoLock(<-self.keyChan)
	}
	//create timestamp entry
	self.entryMap[key] = &TimestampEntry{entry,time.Now().Unix()}
	self.keyChan <- key
	self.size = self.size + 1
	//start clean worker
	if self.Duration>0{
		if atomic.CompareAndSwapInt32(&self.isCleanerRuning,0,1){
			go cleanerRuning(self)
		}
	}
}

func (self *LruHashMap)getNoLock(key interface{})interface{}{
	if key == nil{
		return nil
	}
	timeEntry := self.entryMap[key]
	if timeEntry!=nil{
		timeEntry.timestamp = time.Now().Unix()
		return timeEntry.entry
	}
	return nil
}

func (self *LruHashMap)Get(key interface{})interface{}{
	self.lock.RLock()
	defer self.lock.RUnlock()
	return self.getNoLock(key)
}

func (self *LruHashMap)removeNoLock(key interface{})interface{}{
	if key == nil{
		return nil
	}
	old := self.getNoLock(key)
	delete(self.entryMap,key)
	if old != nil{
		self.size = self.size - 1
	}
	return old
}

func (self *LruHashMap)Remove(key interface{})interface{}{
	self.lock.Lock()
	defer self.lock.Unlock()
	return self.removeNoLock(key)
}

func (self *LruHashMap)Size()int{
	self.lock.RLock()
	defer self.lock.RUnlock()
	return self.size
}
/**
 * Is contains key
 */
func (self *LruHashMap)ContainsKey(key interface{})bool{
	if key == nil{
		return false
	}
	self.lock.RLock()
	defer self.lock.RUnlock()
	if _,ok:=self.entryMap[key]; ok{
		return true
	}
	return false
}

/**
 * Clear hash map
 */
func (self *LruHashMap)Clear(){
	self.lock.Lock()
	defer self.lock.Unlock()
	self.entryMap = make(map[interface{}]*TimestampEntry)
	self.size = 0
}


//func main(){
//	lhm := NewLruHashMap(2,10)
//	lhm.Put("xiao","aaaa")
//	lhm.Put("xiao2","bbbb")
//	time.Sleep(1*time.Second)
//	//lhm.Remove("xiao")
//	fmt.Println(lhm.Get("xiao2")," ",lhm.Size())
//	
//	fmt.Println(lhm.ContainsKey("xiao"))
//	
//	
//}