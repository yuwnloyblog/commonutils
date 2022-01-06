package commonutils

import (
	"fmt"
	"hash/crc32"
	"sort"
	"strconv"
	"sync"
)

const DEFAULT_REPLICAS = 160

type HashRing []uint32

func (c HashRing) Len() int {
	return len(c)
}

func (c HashRing) Less(i, j int) bool {
	return c[i] < c[j]
}

func (c HashRing) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

type Node struct {
	Name   string
	Entry  interface{}
	Weight int
}

type ConsistentHash struct {
	Nodes      map[uint32]*Node
	numReps    int
	Resources  map[string]bool
	isAutoSort bool
	ring       HashRing
	sync.RWMutex
}

func NewConsistentHash(isAutoSort bool) *ConsistentHash {
	return &ConsistentHash{
		Nodes:      make(map[uint32]*Node),
		numReps:    DEFAULT_REPLICAS,
		Resources:  make(map[string]bool),
		ring:       HashRing{},
		isAutoSort: isAutoSort,
	}
}

func (c *ConsistentHash) Add(name string, entry interface{}, weight int) bool {
	node := &Node{
		Name:   name,
		Entry:  entry,
		Weight: weight,
	}
	return c.addNode(node)
}

func (c *ConsistentHash) addNode(node *Node) bool {
	c.Lock()
	defer c.Unlock()

	if _, ok := c.Resources[node.Name]; ok {
		return false
	}

	count := c.numReps * node.Weight

	c.Nodes[c.hashStr(node.Name)] = node
	for i := 0; i < count; i++ {
		str := c.joinStr(i, node)
		c.Nodes[c.hashStr(str)] = node
	}
	c.Resources[node.Name] = true
	if c.isAutoSort {
		c.sortHashRing()
	}
	return true
}

func (c *ConsistentHash) sortHashRing() {
	c.ring = HashRing{}
	for k := range c.Nodes {
		c.ring = append(c.ring, k)
	}
	sort.Sort(c.ring)
}

func (c *ConsistentHash) Prepare() {
	c.sortHashRing()
}

func (c *ConsistentHash) joinStr(i int, node *Node) string {
	return node.Name + "*" + strconv.Itoa(node.Weight) +
		"-" + strconv.Itoa(i)
}

// MurMurHash算法 :https://github.com/spaolacci/murmur3
func (c *ConsistentHash) hashStr(key string) uint32 {
	return crc32.ChecksumIEEE([]byte(key))
}

func (c *ConsistentHash) Get(key string) *Node {
	c.RLock()
	defer c.RUnlock()
	if len(c.Resources) <= 0 {
		return nil
	}
	hash := c.hashStr(key)
	i := c.search(hash)

	return c.Nodes[c.ring[i]]
}

func (c *ConsistentHash) search(hash uint32) int {

	i := sort.Search(len(c.ring), func(i int) bool { return c.ring[i] >= hash })
	if i < len(c.ring) {
		if i == len(c.ring)-1 {
			return 0
		} else {
			return i
		}
	} else {
		return len(c.ring) - 1
	}
}

func (c *ConsistentHash) Remove(node *Node) {
	c.Lock()
	defer c.Unlock()

	if _, ok := c.Resources[node.Name]; !ok {
		return
	}

	delete(c.Resources, node.Name)

	count := c.numReps * node.Weight
	for i := 0; i < count; i++ {
		str := c.joinStr(i, node)
		delete(c.Nodes, c.hashStr(str))
	}
	if c.isAutoSort {
		c.sortHashRing()
	}
}

func main() {

	cHashRing := NewConsistentHash(true)

	for i := 0; i < 10; i++ {
		si := fmt.Sprintf("%d", i)
		cHashRing.Add("name"+si, 8080, 1)
	}

	for k, v := range cHashRing.Nodes {
		fmt.Println("Hash:", k, " Name:", v.Name)
	}

	ipMap := make(map[string]int, 0)
	for i := 0; i < 1000; i++ {
		si := fmt.Sprintf("key%d", i)
		k := cHashRing.Get(si)
		if _, ok := ipMap[k.Name]; ok {
			ipMap[k.Name] += 1
		} else {
			ipMap[k.Name] = 1
		}
	}

	for k, v := range ipMap {
		fmt.Println("Node Name:", k, " count:", v)
	}

}
