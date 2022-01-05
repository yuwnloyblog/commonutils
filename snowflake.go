package commonutils

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

const (
	one int64 = -1
	//sequence         int64 = 0
	twepoch         int64 = 1288834974657
	maxWorkerId     int64 = -1 ^ (-1 << 5)
	maxDatacenterId int64 = -1 ^ (-1 << 5)
	sequenceBits    int64 = 12
	sequenceMask    int64 = -1 ^ (-1 << 12)

//	lastTimestamp      int64 = -1
)

type SnowflakeIdWorker struct {
	_workerId     int64
	_datacenterId int64
	sequence      int64
	lastTimestamp int64
	lock          sync.Mutex
}

func CreateSnowflake(workerId, datacenterId int64) (*SnowflakeIdWorker, error) {
	if workerId > maxWorkerId || workerId < 0 {
		return nil, errors.New(fmt.Sprintf("worker Id can't be greater than %d or less than 0", maxWorkerId))
	}
	if datacenterId > maxDatacenterId || datacenterId < 0 {
		return nil, errors.New(fmt.Sprintf("datacenter Id can't be greater than %d or less than 0", maxDatacenterId))
	}
	idWorker := &SnowflakeIdWorker{
		_workerId:     workerId,
		_datacenterId: datacenterId,
		sequence:      0,
		lastTimestamp: -1,
	}
	return idWorker, nil
}

func (self *SnowflakeIdWorker) nextId() (int64, error) {
	self.lock.Lock()
	defer self.lock.Unlock()
	timestamp := timeGen()

	if timestamp < self.lastTimestamp {
		return 0, errors.New(fmt.Sprintf("Clock moved backwards.  Refusing to generate id for %d milliseconds", self.lastTimestamp-timestamp))
	}
	if self.lastTimestamp == timestamp {
		self.sequence = (self.sequence + 1) & sequenceMask
		if self.sequence == 0 {
			timestamp = tilNextMillis(self.lastTimestamp)
		}
	} else {
		self.sequence = 0
	}
	self.lastTimestamp = timestamp

	return ((timestamp - twepoch) << 22) | (self._datacenterId << 17) | (self._workerId << 12) | self.sequence, nil
}

func tilNextMillis(lastTimestamp int64) int64 {
	ts := timeGen()
	for ts <= lastTimestamp {
		ts = timeGen()
	}
	return ts
}

func timeGen() int64 {
	return time.Now().UnixNano() / 1000000
}
