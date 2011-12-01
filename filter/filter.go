// Copyright 2011 Petar Maymounkov. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package filter

import (
	"hash"
	"hash/fnv"
	"sync"
)

// TODO: 
//	* add negated literals
//	* global enable/disable switch

func uint64Bytes(s []byte) uint64 {
	var x uint64
	for i := 0; i < 8; i++ {
		x |= uint64(s[i]) << uint(8*i)
	}
	return x
}

// Filter keeps track of all literals, the current selection (a boolean predicate), and term attributes
type Filter struct {
	sync.Mutex
	selected map[string]int
	attr     map[int64]*attrSet
	hlock    sync.Mutex
	hash     hash.Hash
}

func (x *Filter) hashStringSlice(ss []string) int64 {
	x.hlock.Lock()
	defer x.hlock.Unlock()
	x.hash.Reset()
	var h uint64
	for _, s := range ss {
		x.hash.Write([]byte(s))
		h ^= uint64Bytes(x.hash.Sum())
	}
	return int64(h)
}

func NewFilter() *Filter {
	x := &Filter{}
	x.init()
	return x
}

func (x *Filter) init() {
	x.selected = make(map[string]int)
	x.attr = make(map[int64]*attrSet)
	x.hash = fnv.New64a()
}

// Set the select status of a literal

func (x *Filter) Select(literals ...string) {
	x.Lock()
	defer x.Unlock()
	for _, literal := range literals {
		x.selected[literal] = 1
	}
}

func (x *Filter) Unselect(literals ...string) {
	x.Lock()
	defer x.Unlock()
	for _, literal := range literals {
		delete(x.selected, literal)
	}
}

func (x *Filter) Selected(literals ...string) bool {
	x.Lock()
	defer x.Unlock()
	for _, literal := range literals {
		_, ok := x.selected[literal]
		if !ok {
			return false
		}
	}
	return true
}

// XXX: varargs for term
func (x *Filter) SetAttr(term []string, attr string, value interface{}) {
	if value == nil {
		x.UnsetAttr(term, attr)
	}
	x.Lock()
	defer x.Unlock()
	h := x.hashStringSlice(term)
	a, ok := x.attr[h]
	if !ok {
		a = newAttrSet()
		x.attr[h] = a
	}
	a.SetAttr(attr, value)
}

// XXX: varargs for term
func (x *Filter) UnsetAttr(term []string, attr string) {
	x.Lock()
	defer x.Unlock()
	h := x.hashStringSlice(term)
	a, ok := x.attr[h]
	if !ok {
		return
	}
	a.UnsetAttr(attr)
	if a.Len() == 0 {
		delete(x.attr, h)
	}
}

// XXX: varargs for term
func (x *Filter) GetAttr(term []string, attr string) interface{} {
	x.Lock()
	defer x.Unlock()
	h := x.hashStringSlice(term)
	a, ok := x.attr[h]
	if !ok {
		return nil
	}
	return a.GetAttr(attr)
}

// attrSet is a set of attributes of the form (attrName, value)
type attrSet struct {
	sync.Mutex
	attr map[string]interface{}
}

func newAttrSet() *attrSet {
	return &attrSet{
		attr: make(map[string]interface{}),
	}
}

func (t *attrSet) Len() int {
	t.Lock()
	defer t.Unlock()
	return len(t.attr)
}

func (t *attrSet) SetAttr(attr string, value interface{}) {
	t.Lock()
	defer t.Unlock()
	if value == nil {
		delete(t.attr, attr)
	} else {
		t.attr[attr] = value
	}
}

func (t *attrSet) UnsetAttr(attr string) {
	t.Lock()
	defer t.Unlock()
	delete(t.attr, attr)
}

func (t *attrSet) GetAttr(attr string) interface{} {
	t.Lock()
	defer t.Unlock()
	return t.attr[attr]
}
