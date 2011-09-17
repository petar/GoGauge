// Copyright 2011 Petar Maymounkov. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package context

import "sync"

type Context struct {
	sync.Mutex
	name     string
	selected bool
	parent   *Context
	children map[string]*Context
	attr     map[string]interface{}
}

func NewContext(name string) *Context {
	return &Context{ 
		name:     name, 
		selected: true, 
		parent:   nil,
		children: make(map[string]*Context),
		attr:     make(map[string]interface{}),
	}
}

func (t *Context) Make(name string) *Context {
	t.Lock()
	defer t.Unlock()

	_, ok := t.children[name]
	if ok {
		panic("child context exists")
	}
	c := NewContext(name)
	c.parent = t
	t.children[name] = c
	return c
}

func (t *Context) Get(name string) *Context {
	t.Lock()
	defer t.Unlock()

	return t.children[name]
}

func (t *Context) Select(v bool) {
	t.Lock()
	defer t.Unlock()

	t.selected = v
}

func (t *Context) Name() string {
	t.Lock()
	defer t.Unlock()

	return t.name
}

func (t *Context) SetAttr(attr string, value interface{}) {
	t.Lock()
	defer t.Unlock()

	t.attr[attr] = value
}

func (t *Context) GetAttr(attr string) interface{} {
	t.Lock()
	defer t.Unlock()

	return t.attr[attr]
}

func (t *Context) Path() []*Context {
	r := make([]*Context, 0, 3)
	c := t
	for c != nil {
		r = append(r, c)
		d := c
		d.Lock()
		c = d.parent
		d.Unlock()
	}
	return r
}

func (t *Context) NamePath() []string {
	p := t.Path()
	q := make([]string, len(p))
	for i, c := range p {
		c.Lock()
		q[i] = c.name
		c.Unlock()
	}
	return q
}

func (t *Context) Selected() bool {
	p := t.Path()
	for _, c := range p {
		c.Lock()
		s := c.selected
		c.Unlock()
		if !s {
			return false
		}
	}
	return true
}
