// Copyright 2011 Petar Maymounkov. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package unclosed

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"runtime"
	"strconv"
	"sync"
	//"time"
)

type readCloserTracker struct {
	io.ReadCloser
	lock     sync.Mutex
	name     string
	isClosed bool
}

func (t *readCloserTracker) Close() os.Error {
	t.lock.Lock()
	if !t.isClosed {
		t.isClosed = true
		Closed(t.name)
	}
	t.lock.Unlock()
	return t.ReadCloser.Close()
}

// NewReadCloserTracker creates a new io.ReadCloser which wraps rc
// and keeps track of when the returned io.ReadCloser is closed.
func NewReadCloserTracker(rc io.ReadCloser) io.ReadCloser {
	var name string
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		name = "unknown"
	} else {
		name = file + ":" + strconv.Itoa(line)
	}
	Opened(name)
	return &readCloserTracker{
		ReadCloser: rc, 
		name:       name,
		isClosed:   false,
	}
}

var (
	lk       sync.Mutex
	unclosed map[string]int = make(map[string]int)	// Maps name to number of opens
)

func Opened(name string) {
	lk.Lock()
	defer lk.Unlock()

	n, _ := unclosed[name]
	unclosed[name] = n + 1
}

func Closed(name string) {
	lk.Lock()
	defer lk.Unlock()

	n, ok := unclosed[name]
	if !ok {
		panic("open/close mismatch")
	}
	if n == 1 {
		unclosed[name] = 0, false
	} else {
		unclosed[name] = n-1
	}
}

func Report() string {
	lk.Lock()
	defer lk.Unlock()

	var w bytes.Buffer
	for n, k := range unclosed {
		fmt.Fprintf(&w, "%s: # open = %d\n", n, k)
	}
	return string(w.Bytes())
}
