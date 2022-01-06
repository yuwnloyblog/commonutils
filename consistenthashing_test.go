package commonutils

import (
	"fmt"
	"strconv"
	"sync"
	"testing"
)

func TestNewConsistentHash(t *testing.T) {
	hashRing := NewConsistentHash(true)

	hashRing.Add("aaa", "aaaaa", 1)
	hashRing.Add("bbb", "bbb", 1)
	hashRing.Add("ccc", "ccc", 1)

	m := sync.Map{}

	for i := 0; i < 100000; i++ {
		nod := hashRing.Get("abc" + strconv.Itoa(i))

		act, isExt := m.LoadOrStore(nod.Name, 1)
		if isExt {
			m.Store(nod.Name, act.(int)+1)
		}
	}
	m.Range(func(key, value interface{}) bool {
		fmt.Println(key, "", value)
		return true
	})
}
