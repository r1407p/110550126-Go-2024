package intset

import (
	"bytes"
	"fmt"
)

type IntSet struct {
	words []uint64
}

func (self *IntSet) Add(x int) {
	word, bit := x/64, uint(x%64)
	for word >= len(self.words) {
		self.words = append(self.words, 0)
	}
	self.words[word] |= 1 << bit
}

func (self *IntSet) Has(x int) bool {
	word, bit := x/64, uint(x%64)
	return word < len(self.words) && self.words[word]&(1<<bit) != 0
}


func (self *IntSet) UnionWith(t *IntSet) {
	for i, tword := range t.words {
		if i < len(self.words) {
			self.words[i] |= tword
		} else {
			self.words = append(self.words, tword)
		}
	}
}

func (self *IntSet) Len() int {
	count := 0
	for _, word := range self.words {
		count += popCount(word)
	}
	return count
}

func (self *IntSet) Remove(x int) {
	word, bit := x/64, uint(x%64)
	if word < len(self.words) {
		self.words[word] &^= 1 << bit
	}
}

func (self *IntSet) Clear() {
	self.words = nil
}

func (self *IntSet) Copy() *IntSet {
	newSet := &IntSet{}
	newSet.words = append([]uint64{}, self.words...)
	return newSet
}

func (self *IntSet) AddAll(values ...int) {
	for _, value := range values {
		self.Add(value)
	}
}

func (self *IntSet) IntersectWith(t *IntSet) {
	for i := range self.words {
		if i < len(t.words) {
			self.words[i] &= t.words[i]
		} else {
			self.words[i] = 0
		}
	}
}

func (self *IntSet) DifferenceWith(t *IntSet) {
	for i := range self.words {
		if i < len(t.words) {
			self.words[i] &^= t.words[i]
		}
	}
}

func (self *IntSet) SymmetricDifference(t *IntSet) {
	for i, tword := range t.words {
		if i < len(self.words) {
			self.words[i] ^= tword
		} else {
			self.words = append(self.words, tword)
		}
	}
}

func (self *IntSet) String() string {
	var buf bytes.Buffer
	buf.WriteByte('{')
	for i, word := range self.words {
		if word == 0 {
			continue
		}
		for j := 0; j < 64; j++ {
			if word&(1<<uint(j)) != 0 {
				if buf.Len() > len("{") {
					buf.WriteByte(' ')
				}
				fmt.Fprintf(&buf, "%d", 64*i+j)
			}
		}
	}
	buf.WriteByte('}')
	return buf.String()
}

func popCount(x uint64) int {
	count := 0
	for x != 0 {
		count++
		x &= x - 1
	}
	return count
}

