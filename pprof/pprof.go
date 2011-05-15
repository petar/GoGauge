// Copyright 2011 Petar Maymounkov. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pprof

import (
	"log"
	"http"
	_ "http/pprof"
)

func InstallPprof() {
	go func() {
		err := http.ListenAndServe(":5500", nil)
		if err != nil {
			log.Fatal("ListenAndServe: ", err.String())
		}
	}()
}
