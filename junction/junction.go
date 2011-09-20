// Copyright 2011 Petar Maymounkov. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package junction

import "sync"

func Select(literal string, selected bool) {
	_omega.fetchLiteral(literal).Select(selected)
}

func SetAttr(literal string, attr string, value interface{}) {
	_omega.fetchLiteral(literal).SetAttr(attr, value)
}

func UnsetAttr(literal string, attr string) {
	_omega.fetchLiteral(literal).UnsetAttr(attr)
}

type omega struct {
	sync.Mutex
	literals map[string]*Literal
}

var _omega omega

func init() {
	_omega.Init()
}

func (t *omega) Init() {
	t.literals = make(map[string]*Literal)
}

func (t *omega) fetchLiteral(literal string) *Literal {
	t.Lock()
	defer t.Unlock()
	l, ok := t.literals[literal]
	if ok {
		return l
	}
	l = newLiteral(literal)
	t.literals[literal] = l
	return l
}

func (t *omega) selectLiteral(literal string, selected bool) {
	t.Lock()
	l := t.literals[literal]
	t.Unlock()
	l.Select(selected)
}

type Literal struct {
	sync.Mutex
	literal  string
	selected bool
	attr     map[string]interface{}
}

func newLiteral(name string) *Literal {
	return &Literal{
		literal:  name,
		selected: true,
		attr:     make(map[string]interface{}),
	}
}

func (t *Literal) SetAttr(attr string, value interface{}) {
	t.Lock()
	defer t.Unlock()
	t.attr[attr] = value
}

func (t *Literal) UnsetAttr(attr string) {
	t.Lock()
	defer t.Unlock()
	t.attr[attr] = nil, false
}

func (t *Literal) GetAttr(attr string) interface{} {
	t.Lock()
	defer t.Unlock()
	return t.attr[attr]
}

func (t *Literal) Select(value bool) {
	t.Lock()
	defer t.Unlock()
	t.selected = value
}

func (t *Literal) Selected() bool {
	t.Lock()
	defer t.Unlock()
	return t.selected
}

type Clause []string

func (t Clause) Selected() bool {
	for _, l := range t {
		if !_omega.fetchLiteral(l).Selected() {
			return false
		}
	}
	return true
}
